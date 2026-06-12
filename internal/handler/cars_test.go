package handler

import (
	"bytes"
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/sindrema/teslamap/internal/model"
)

// syncBuf is a concurrency-safe io.Writer + http.Flusher so the test can read
// the streamed bytes while writePositionStream writes them from a goroutine.
type syncBuf struct {
	mu  sync.Mutex
	buf bytes.Buffer
}

func (s *syncBuf) Write(p []byte) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.Write(p)
}
func (s *syncBuf) Flush() {}
func (s *syncBuf) String() string {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.buf.String()
}

func dataFrames(s string) int { return strings.Count(s, "data: ") }

// waitFor polls until cond() or the deadline, failing the test on timeout.
func waitFor(t *testing.T, cond func() bool, msg string) {
	t.Helper()
	deadline := time.Now().Add(2 * time.Second)
	for time.Now().Before(deadline) {
		if cond() {
			return
		}
		time.Sleep(2 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for %s", msg)
}

func TestWritePositionStream_EmitsOnChangeOnly(t *testing.T) {
	car := model.Car{ID: 1, Model: "Model Y"}
	d1 := time.Date(2026, 1, 1, 12, 0, 0, 0, time.UTC)
	posA := &model.Position{Latitude: 1, Longitude: 2, Date: d1}
	posB := &model.Position{Latitude: 3, Longitude: 4, Date: d1.Add(time.Minute)}

	var mu sync.Mutex
	cur := posA
	fetch := func(context.Context) (*model.Position, error) {
		mu.Lock()
		defer mu.Unlock()
		return cur, nil
	}

	sb := &syncBuf{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	// heartbeat far in the future so this test sees only change-driven frames
	go func() {
		writePositionStream(ctx, sb, sb, car, fetch, 5*time.Millisecond, time.Hour)
		close(done)
	}()

	// initial snapshot frame for posA
	waitFor(t, func() bool { return dataFrames(sb.String()) == 1 }, "initial frame")

	// unchanged position should NOT produce more frames
	time.Sleep(40 * time.Millisecond)
	if n := dataFrames(sb.String()); n != 1 {
		t.Fatalf("expected no extra frames while unchanged, got %d:\n%s", n, sb.String())
	}

	// change → exactly one more frame
	mu.Lock()
	cur = posB
	mu.Unlock()
	waitFor(t, func() bool { return dataFrames(sb.String()) == 2 }, "frame after change")

	cancel()
	<-done

	out := sb.String()
	if dataFrames(out) != 2 {
		t.Fatalf("expected 2 data frames total, got %d:\n%s", dataFrames(out), out)
	}
	if !strings.Contains(out, `"latitude":1`) || !strings.Contains(out, `"latitude":3`) {
		t.Fatalf("frames missing expected positions:\n%s", out)
	}
	if strings.Contains(out, ": keepalive") {
		t.Fatalf("did not expect a heartbeat in this test:\n%s", out)
	}
	// each frame must be terminated by the SSE blank line
	if !strings.Contains(out, "\n\n") {
		t.Fatalf("frames not terminated with blank line:\n%s", out)
	}
}

func TestWritePositionStream_NoPositionThenHeartbeat(t *testing.T) {
	car := model.Car{ID: 2, Model: "Model 3"}
	fetch := func(context.Context) (*model.Position, error) {
		return nil, nil // car has never reported a position
	}

	sb := &syncBuf{}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan struct{})
	go func() {
		writePositionStream(ctx, sb, sb, car, fetch, 5*time.Millisecond, 10*time.Millisecond)
		close(done)
	}()

	// initial frame carries position:null, then idle heartbeats follow
	waitFor(t, func() bool {
		return strings.Contains(sb.String(), `"position":null`)
	}, "null-position frame")
	waitFor(t, func() bool {
		return strings.Contains(sb.String(), ": keepalive")
	}, "heartbeat")

	cancel()
	<-done

	if dataFrames(sb.String()) != 1 {
		t.Fatalf("expected exactly one data frame for an unchanging null position:\n%s", sb.String())
	}
}
