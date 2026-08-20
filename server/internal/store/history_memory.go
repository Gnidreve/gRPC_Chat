package store

import (
	"context"
	"sync"
)

// memoryHistory is a History that keeps messages in a plain slice — used by
// tests and as the fallback when no Redis is configured. Lost on restart.
type memoryHistory struct {
	mu   sync.RWMutex
	msgs []Message
}

// NewMemoryHistory returns a History that keeps messages in process memory
// only — simple, but lost on every restart. Prefer NewRedisHistory in
// production so history survives redeploys.
func NewMemoryHistory() History {
	return &memoryHistory{}
}

func (h *memoryHistory) Append(_ context.Context, msg Message) error {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.msgs = append(h.msgs, msg)
	if len(h.msgs) > maxHistory {
		h.msgs = h.msgs[len(h.msgs)-maxHistory:]
	}
	return nil
}

func (h *memoryHistory) All(_ context.Context) ([]Message, error) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	out := make([]Message, len(h.msgs))
	copy(out, h.msgs)
	return out, nil
}
