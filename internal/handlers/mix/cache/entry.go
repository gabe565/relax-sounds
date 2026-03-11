package cache

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"sync"
	"time"

	"gabe565.com/relax-sounds/internal/handlers/mix/encoder"
	"gabe565.com/relax-sounds/internal/handlers/mix/stream"
	"github.com/gopxl/beep/v2"
	"github.com/pocketbase/pocketbase/core"
	"github.com/valkey-io/valkey-go"
)

type Entry struct {
	Log *slog.Logger

	UUID    string
	Preset  string
	Streams stream.Streams
	Format  beep.Format

	Writer      *ProxyWriter
	Encoder     encoder.Encoder
	positionErr error

	Mu      sync.Mutex
	Created time.Time
}

func NewEntry(e *core.RequestEvent, uuid, preset string) *Entry {
	return &Entry{
		Log: slog.With(
			"userIp", e.RealIP(),
			"userAgent", e.Request.UserAgent(),
			"url", e.Request.URL.String(),
			"id", uuid,
		),
		UUID:    uuid,
		Preset:  preset,
		Writer:  NewProxyWriter(),
		Created: time.Now(),
	}
}

func (e *Entry) Close() error {
	e.Mu.Lock()
	defer e.Mu.Unlock()
	return e.close()
}

func (e *Entry) close() error {
	if e.Writer != nil {
		e.Writer.Close()
		defer func() {
			e.Writer = nil
		}()
	}

	errs := make([]error, 0, 2)
	errs = append(errs, e.Streams.Close())
	if e.Encoder != nil {
		errs = append(errs, e.Encoder.Close())
	}
	return errors.Join(errs...)
}

func (e *Entry) valkeyKey(topic string) string {
	hash := sha256.Sum256([]byte(e.Preset))
	hashHex := hex.EncodeToString(hash[:])
	return e.UUID + ":" + hashHex + ":" + topic
}

func (e *Entry) StorePositions(ctx context.Context, v valkey.Client) error {
	positions := make([]int, 0, len(e.Streams))
	for _, s := range e.Streams {
		positions = append(positions, s.Closer.Position())
	}

	b, err := json.Marshal(positions)
	if err != nil {
		e.positionErr = err
		return err
	}

	res := v.Do(ctx,
		v.B().Set().Key(e.valkeyKey("position")).Value(valkey.BinaryString(b)).Ex(10*time.Minute).Build(),
	)
	e.positionErr = res.Error()
	return e.positionErr
}

var ErrPositionLenMismatch = errors.New("position length mismatch")

func (e *Entry) LoadPositions(ctx context.Context, v valkey.Client) (bool, error) {
	if e.positionErr != nil {
		return false, nil //nolint:nilerr
	}

	res := v.Do(ctx,
		v.B().Get().Key(e.valkeyKey("position")).Build(),
	)

	b, err := res.AsBytes()
	if err != nil {
		if valkey.IsValkeyNil(err) {
			return false, nil
		}
		return false, err
	}

	positions := make([]int, 0, len(e.Streams))
	if err := json.Unmarshal(b, &positions); err != nil {
		return false, err
	}

	if len(positions) != len(e.Streams) {
		return false, fmt.Errorf("%w: expected %d, got %d", ErrPositionLenMismatch, len(e.Streams), len(positions))
	}

	var errs []error
	for i, s := range e.Streams {
		if s.Closer.Position() != positions[i] {
			if err := s.Closer.Seek(positions[i]); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return len(errs) < len(e.Streams), errors.Join(errs...)
}
