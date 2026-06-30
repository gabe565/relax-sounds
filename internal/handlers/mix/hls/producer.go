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
	framesPerSegment   = 230
	mp3SamplesPerFrame = 1152
	segmentSampleRate  = 44100
	pcmChunkSamples    = 1152
)

// SegmentDuration is the wall-clock length of one segment.
func SegmentDuration() time.Duration {
	return time.Duration(framesPerSegment) * mp3SamplesPerFrame * time.Second / segmentSampleRate
}

// bufferAhead caps how far the encoder gets ahead of wall-clock. The initial
// burst fills the ring buffer fast (so the first manifest request resolves
// quickly), then this gate makes the producer self-pace.
const bufferAhead = 30 * time.Second

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

	var (
		segFrames int
		start     = time.Now()
		encoded   time.Duration
		// scanPos tracks bytes of the current (not-yet-emitted) segment that
		// have already been parsed, measured from raw's read offset.
		scanPos int
	)

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

		// Drain complete MP3 frames from raw.
		for {
			frameLen, ok := parseFrameLen(raw.Bytes(), scanPos)
			if !ok {
				break
			}
			scanPos += frameLen
			segFrames++

			if segFrames < framesPerSegment {
				continue
			}

			dur := e.emitSegment(raw.Next(scanPos), segFrames)
			encoded += dur
			scanPos = 0
			segFrames = 0

			if e.throttleAhead(encoded - time.Since(start)) {
				return
			}
		}
	}
}

// emitSegment clones the given byte range as a finalized segment and pushes
// it into the ring buffer.
func (e *Entry) emitSegment(b []byte, segFrames int) time.Duration {
	dur := time.Duration(segFrames) * mp3SamplesPerFrame * time.Second / segmentSampleRate
	e.PushSegment(&Segment{
		Duration: dur,
		Bytes:    bytes.Clone(b),
		Created:  time.Now(),
	})
	return dur
}

// throttleAhead sleeps when the encoder has run too far ahead of wall-clock,
// returning true if the context was canceled during the sleep.
func (e *Entry) throttleAhead(ahead time.Duration) bool {
	if ahead <= bufferAhead {
		return false
	}
	select {
	case <-e.ctx.Done():
		return true
	case <-time.After(ahead - bufferAhead):
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
