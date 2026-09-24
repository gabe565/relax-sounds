package hls

import (
	"bytes"
	"errors"
	"time"

	"gabe565.com/relax-sounds/internal/config"
	"github.com/gopxl/beep/v2"
	"github.com/viert/go-lame"
)

const (
	framesPerSegment   = 460
	mp3SamplesPerFrame = 1152
	segmentSampleRate  = 44100
	pcmChunkSamples    = 1152
)

// SegmentDuration is the wall-clock length of one segment.
func SegmentDuration() time.Duration {
	return time.Duration(framesPerSegment) * mp3SamplesPerFrame * time.Second / segmentSampleRate
}

// bufferAhead is how far ahead of wall clock segments are published. A new
// entry bursts until it reaches this lead, then each segment waits until
// bufferAhead before its seqStart. Because the schedule is absolute, an entry
// recreated for the same UUID resumes at the sequence its predecessor reached.
const bufferAhead = 30 * time.Second

// produceState tracks byte/frame progress within a single Produce call.
//
// scanPos tracks bytes of the current (not-yet-emitted) segment that have
// already been parsed, measured from raw's read offset.
type produceState struct {
	scanPos   int
	segFrames int
	// totalFrames counts frames emitted in prior segments, giving each new
	// segment the presentation time of its first sample.
	totalFrames uint64
}

// Produce runs a single continuous LAME encoder and slices its byte stream at
// MP3 frame boundaries. This avoids the per-segment Xing-tag encoder-delay
// gap (~350ms) that a fresh-encoder-per-segment approach would introduce at
// every seam.
func (e *Entry) Produce(conf *config.Config) {
	format := beep.Format{SampleRate: segmentSampleRate, NumChannels: 2, Precision: 2}
	mix := e.Streams.Mix()

	var raw bytes.Buffer
	enc := lame.NewEncoder(&raw)
	defer enc.Close()
	if err := enc.SetVBR(lame.VBRDefault); err != nil {
		e.Log.Error("LAME SetVBR failed", "error", err)
		return
	}
	if err := enc.SetVBRQuality(conf.LAMEQuality); err != nil {
		e.Log.Error("LAME SetVBRQuality failed", "error", err)
		return
	}

	samples := make([][2]float64, pcmChunkSamples)
	pcm := make([]byte, len(samples)*format.Width())
	var state produceState

	for {
		if err := e.ctx.Err(); err != nil {
			return
		}

		// Feed one PCM chunk to LAME. LAME emits MP3 frames into raw as it
		// accumulates enough samples.
		n, _ := mix.Stream(samples)
		if n == 0 {
			e.Log.Warn("HLS mix returned 0 samples")
			return
		}
		var off int
		for _, s := range samples[:n] {
			off += format.EncodeSigned(pcm[off:], s)
		}
		if _, err := enc.Write(pcm[:n*format.Width()]); err != nil {
			if !errors.Is(err, ErrClosed) {
				e.Log.Error("LAME write failed", "error", err)
			}
			return
		}

		if e.drainFrames(&raw, &state) {
			return
		}
	}
}

// drainFrames extracts complete MP3 frames from raw, emitting a segment each
// time framesPerSegment accumulate. Returns true if the caller should stop
// (context canceled during throttle wait).
func (e *Entry) drainFrames(raw *bytes.Buffer, s *produceState) bool {
	for {
		frameLen, ok := parseFrameLen(raw.Bytes(), s.scanPos)
		if !ok {
			return false
		}
		s.scanPos += frameLen
		s.segFrames++

		if s.segFrames < framesPerSegment {
			continue
		}

		next := e.emitSegment(raw.Next(s.scanPos), s.segFrames, s.totalFrames)
		s.totalFrames += uint64(s.segFrames)
		s.scanPos = 0
		s.segFrames = 0

		if e.throttle(next) {
			return true
		}
	}
}

// emitSegment prefixes the given byte range with an ID3 timestamp tag and
// pushes it into the ring buffer as a finalized segment. It returns the
// sequence the following segment will receive.
func (e *Entry) emitSegment(b []byte, segFrames int, startFrame uint64) uint64 {
	dur := time.Duration(segFrames) * mp3SamplesPerFrame * time.Second / segmentSampleRate

	tag := id3TimestampTag(timestampFor(startFrame * mp3SamplesPerFrame))
	segment := make([]byte, 0, len(tag)+len(b))
	segment = append(segment, tag...)
	segment = append(segment, b...)

	return e.PushSegment(&Segment{
		Duration: dur,
		Bytes:    segment,
		Created:  time.Now(),
	})
}

// throttle holds the producer until bufferAhead before next's seqStart. The
// first time it has to wait, the startup burst is complete and the entry is
// marked ready. Returns true if the context was canceled during the sleep.
func (e *Entry) throttle(next uint64) bool {
	wait := time.Until(seqStart(next).Add(-bufferAhead))
	if wait <= 0 {
		return false
	}
	e.markReady()
	select {
	case <-e.ctx.Done():
		return true
	case <-time.After(wait):
		return false
	}
}

//nolint:gochecknoglobals // MPEG spec constant lookup tables
var (
	// MPEG-1 Layer III bitrate table (kbps), indexed by the 4-bit bitrate field.
	mp3Bitrates = [16]int{0, 32, 40, 48, 56, 64, 80, 96, 112, 128, 160, 192, 224, 256, 320, 0}
	// MPEG-1 sample rate table (Hz), indexed by the 2-bit sample-rate field.
	mp3SampleRates = [4]int{44100, 48000, 32000, 0}
)

// parseFrameLen looks at data starting from pos and, if a complete MPEG-1
// Layer III frame is available, returns its total length in bytes. Returns
// ok=false if not enough data or the header isn't a valid frame sync.
func parseFrameLen(data []byte, pos int) (int, bool) {
	if pos+4 > len(data) {
		return 0, false
	}
	if data[pos] != 0xFF || data[pos+1]&0xE0 != 0xE0 {
		return 0, false
	}
	// version=11 (MPEG-1) and layer=01 (Layer III) are what LAME emits.
	if (data[pos+1]>>3)&0x3 != 3 || (data[pos+1]>>1)&0x3 != 1 {
		return 0, false
	}
	bitrateIdx := (data[pos+2] >> 4) & 0xF
	sampleRateIdx := (data[pos+2] >> 2) & 0x3
	padding := int((data[pos+2] >> 1) & 0x1)
	bitrate := mp3Bitrates[bitrateIdx]
	sampleRate := mp3SampleRates[sampleRateIdx]
	if bitrate == 0 || sampleRate == 0 {
		return 0, false
	}
	frameLen := (144*bitrate*1000)/sampleRate + padding
	if pos+frameLen > len(data) {
		return 0, false
	}
	return frameLen, true
}
