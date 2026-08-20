package store

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// newTestRedisHistory returns a redisHistory pointed at a scratch key on a
// local/CI Redis instance. Skips (not fails) if no Redis is reachable, so
// `go test ./...` still works on a machine without Redis — CI runs a real
// Redis service container so the test suite exercises this for real there.
func newTestRedisHistory(t *testing.T) History {
	t.Helper()

	addr := os.Getenv("TEST_REDIS_ADDR")
	if addr == "" {
		addr = "localhost:6379"
	}
	client := redis.NewClient(&redis.Options{Addr: addr})
	t.Cleanup(func() { _ = client.Close() })

	pingCtx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		t.Skipf("redis not reachable at %s: %v", addr, err)
	}

	key := "test:history:" + newID()
	t.Cleanup(func() { _ = client.Del(context.Background(), key).Err() })

	return NewRedisHistory(client, key)
}

func TestRedisHistoryAppendAndAll(t *testing.T) {
	t.Parallel()

	h := newTestRedisHistory(t)
	ctx := context.Background()

	msg1 := Message{User: User{ID: "u1", Nickname: "Mara", Color: "#12B76A"}, Text: "hi", SentAt: time.Now().UTC()}
	msg2 := Message{User: User{ID: "u1", Nickname: "Mara", Color: "#12B76A"}, Text: "there", SentAt: time.Now().UTC()}

	if err := h.Append(ctx, msg1); err != nil {
		t.Fatalf("Append: %v", err)
	}
	if err := h.Append(ctx, msg2); err != nil {
		t.Fatalf("Append: %v", err)
	}

	got, err := h.All(ctx)
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(got) != 2 {
		t.Fatalf("expected 2 messages, got %d: %+v", len(got), got)
	}
	if got[0].Text != "hi" || got[1].Text != "there" {
		t.Fatalf("expected [hi, there] in order, got %+v", got)
	}
	if !got[0].SentAt.Equal(msg1.SentAt) {
		t.Fatalf("expected SentAt %v, got %v", msg1.SentAt, got[0].SentAt)
	}
	if got[0].User != msg1.User {
		t.Fatalf("expected user %+v, got %+v", msg1.User, got[0].User)
	}
}

func TestRedisHistoryCapsAtMaxHistory(t *testing.T) {
	t.Parallel()

	h := newTestRedisHistory(t)
	ctx := context.Background()

	const overflow = 10
	for i := range maxHistory + overflow {
		msg := Message{User: User{ID: "u1"}, Text: "msg", SentAt: time.Now()}
		if err := h.Append(ctx, msg); err != nil {
			t.Fatalf("Append #%d: %v", i, err)
		}
	}

	got, err := h.All(ctx)
	if err != nil {
		t.Fatalf("All: %v", err)
	}
	if len(got) != maxHistory {
		t.Fatalf("expected history capped at %d, got %d", maxHistory, len(got))
	}
}
