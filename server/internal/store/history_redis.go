package store

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
)

// redisHistory is a History backed by a single capped Redis list, so message
// history survives process restarts/redeploys. The rest of the store (users,
// presence, live pub/sub) stays in-memory — see Store.
type redisHistory struct {
	client *redis.Client
	key    string
}

// NewRedisHistory returns a History backed by a Redis list at key on client.
// client must already be configured; NewRedisHistory does not ping it —
// callers should verify connectivity once at startup instead.
func NewRedisHistory(client *redis.Client, key string) History {
	return &redisHistory{client: client, key: key}
}

func (h *redisHistory) Append(ctx context.Context, msg Message) error {
	data, err := json.Marshal(msg)
	if err != nil {
		return fmt.Errorf("marshal message: %w", err)
	}

	pipe := h.client.TxPipeline()
	pipe.RPush(ctx, h.key, data)
	// Keep only the newest maxHistory entries — negative indices count from
	// the tail, so this trims everything before the last maxHistory items.
	pipe.LTrim(ctx, h.key, -maxHistory, -1)
	if _, err := pipe.Exec(ctx); err != nil {
		return fmt.Errorf("append message to redis: %w", err)
	}
	return nil
}

func (h *redisHistory) All(ctx context.Context) ([]Message, error) {
	raw, err := h.client.LRange(ctx, h.key, 0, -1).Result()
	if err != nil {
		return nil, fmt.Errorf("load history from redis: %w", err)
	}

	msgs := make([]Message, 0, len(raw))
	for _, r := range raw {
		var msg Message
		if err := json.Unmarshal([]byte(r), &msg); err != nil {
			return nil, fmt.Errorf("unmarshal message: %w", err)
		}
		msgs = append(msgs, msg)
	}
	return msgs, nil
}
