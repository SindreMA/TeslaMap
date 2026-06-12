package handler

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"time"

	"github.com/sindrema/teslamap/internal/db"
	"github.com/sindrema/teslamap/internal/model"
)

type Handler struct {
	DB *sql.DB
}

// streamPoll is how often the stream re-reads the latest position from the DB.
// Events are only pushed to the client when the position actually changes, so a
// short interval keeps things near-realtime without spamming the wire while parked.
const streamPoll = 1 * time.Second

// streamHeartbeat keeps the connection alive through proxies while idle.
const streamHeartbeat = 15 * time.Second

func (h *Handler) ListCars(w http.ResponseWriter, r *http.Request) {
	cars, err := db.ListCars(r.Context(), h.DB)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if cars == nil {
		cars = []model.Car{}
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(cars)
}

func (h *Handler) GetCarPosition(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	car, err := db.GetCar(r.Context(), h.DB, id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if car == nil {
		http.Error(w, "car not found", http.StatusNotFound)
		return
	}

	pos, err := db.GetLatestPosition(r.Context(), h.DB, id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}

	resp := model.CarPosition{
		Car:      *car,
		Position: pos,
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// StreamCarPosition pushes live position updates over Server-Sent Events.
// Each `data:` frame is the same CarPosition JSON the REST endpoint returns, so
// the client can reuse one code path. A new frame is sent only when the latest
// position changes; otherwise a comment heartbeat keeps the connection open.
func (h *Handler) StreamCarPosition(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.PathValue("id"))
	if err != nil {
		http.Error(w, "invalid car id", http.StatusBadRequest)
		return
	}

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	car, err := db.GetCar(r.Context(), h.DB, id)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	if car == nil {
		http.Error(w, "car not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no") // tell nginx/ingress not to buffer
	w.WriteHeader(http.StatusOK)
	flusher.Flush()

	fetch := func(ctx context.Context) (*model.Position, error) {
		return db.GetLatestPosition(ctx, h.DB, id)
	}
	writePositionStream(r.Context(), w, flusher, *car, fetch, streamPoll, streamHeartbeat)
}

// positionFetcher returns the latest position for a car (nil = no position yet).
type positionFetcher func(context.Context) (*model.Position, error)

// writePositionStream is the SSE loop, separated from HTTP/DB wiring so it can be
// unit-tested. It emits a `data:` frame on connect and whenever the position
// changes, and a `: keepalive` comment when idle past heartbeat. It returns when
// the context is cancelled (client disconnect).
func writePositionStream(
	ctx context.Context,
	w io.Writer,
	flusher http.Flusher,
	car model.Car,
	fetch positionFetcher,
	poll, heartbeat time.Duration,
) {
	ticker := time.NewTicker(poll)
	defer ticker.Stop()

	var lastMarker string
	var lastWrite time.Time

	tick := func() {
		pos, err := fetch(ctx)
		if err != nil {
			return // transient (or ctx cancelled) — the loop exits via ctx.Done
		}

		marker := "none"
		if pos != nil {
			marker = pos.Date.UTC().Format(time.RFC3339Nano)
		}

		if marker == lastMarker {
			if time.Since(lastWrite) >= heartbeat {
				fmt.Fprint(w, ": keepalive\n\n")
				flusher.Flush()
				lastWrite = time.Now()
			}
			return
		}

		payload, err := json.Marshal(model.CarPosition{Car: car, Position: pos})
		if err != nil {
			return
		}
		fmt.Fprintf(w, "data: %s\n\n", payload)
		flusher.Flush()
		lastMarker = marker
		lastWrite = time.Now()
	}

	tick() // initial snapshot immediately on connect

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			tick()
		}
	}
}
