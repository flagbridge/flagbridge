package sse

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
)

type Event struct {
	Type string         `json:"type"`
	Data map[string]string `json:"data"`
}

type client struct {
	ch  chan Event
	env string
}

type Hub struct {
	mu      sync.RWMutex
	clients map[*client]struct{}
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*client]struct{}),
	}
}

func (h *Hub) Broadcast(environment string, event Event) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for c := range h.clients {
		if c.env == environment {
			select {
			case c.ch <- event:
			default:
				// Client too slow, skip
			}
		}
	}
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming unsupported", http.StatusInternalServerError)
		return
	}

	environment := chi.URLParam(r, "environment")
	if environment == "" {
		http.Error(w, "environment required", http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	c := &client{
		ch:  make(chan Event, 32),
		env: environment,
	}

	h.mu.Lock()
	h.clients[c] = struct{}{}
	h.mu.Unlock()

	defer func() {
		h.mu.Lock()
		delete(h.clients, c)
		h.mu.Unlock()
		close(c.ch)
		slog.Debug("SSE client disconnected", "environment", environment)
	}()

	slog.Debug("SSE client connected", "environment", environment)

	// Send initial ping
	fmt.Fprintf(w, "event: connected\ndata: {\"environment\":%q}\n\n", environment)
	flusher.Flush()

	// Keepalive ticker
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-r.Context().Done():
			return
		case event := <-c.ch:
			data, _ := json.Marshal(event.Data)
			fmt.Fprintf(w, "event: %s\ndata: %s\n\n", event.Type, data)
			flusher.Flush()
		case <-ticker.C:
			fmt.Fprintf(w, ": keepalive\n\n")
			flusher.Flush()
		}
	}
}
