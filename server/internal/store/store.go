// Package store holds the chat room state: users, presence and live
// broadcast stay in-memory (a single mutex-guarded struct); message history
// is delegated to a History implementation so it can survive process
// restarts (see history_redis.go) without the rest of the store caring.
package store

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"
)

// ErrUnknownUser is returned when an operation references a user id that
// never joined (or was evicted).
var ErrUnknownUser = errors.New("store: unknown user")

// maxHistory bounds message retention so it doesn't grow without limit.
// Once exceeded, the oldest messages are dropped; connected clients keep
// them, only a fresh Subscribe() loses the tail end of very old history.
const maxHistory = 500

// User is a chat participant: a server-assigned id, plus a nickname and
// color the user picked themselves.
type User struct {
	ID       string
	Nickname string
	Color    string
}

// Message is a single plain-text chat message.
type Message struct {
	User   User
	Text   string
	SentAt time.Time
}

// Event is one item delivered to a subscriber: either a Message or an
// OnlineCount change, never both.
type Event struct {
	Message     *Message
	OnlineCount *int
}

// History persists the message log. AddMessage/Subscribe delegate to it
// instead of holding messages in the Store itself, so storage backend and
// live presence/broadcast state can evolve independently.
type History interface {
	// Append adds msg to the end of the history.
	Append(ctx context.Context, msg Message) error
	// All returns the full retained history, oldest first.
	All(ctx context.Context) ([]Message, error)
}

// Store is the chat room: known users, live subscribers, and a pluggable
// message History.
type Store struct {
	mu          sync.RWMutex
	users       map[string]User
	subscribers map[string]chan Event

	history History
}

// New returns an empty Store backed by the given History.
func New(history History) *Store {
	return &Store{
		users:       make(map[string]User),
		subscribers: make(map[string]chan Event),
		history:     history,
	}
}

// Join registers (or re-registers, e.g. across app restarts) a user under
// the given id, nickname and color — all chosen by the client. Falls back
// to a server-generated id if the client didn't send one.
func (s *Store) Join(id, nickname, color string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	if id == "" {
		id = newID()
	}
	user := User{
		ID:       id,
		Nickname: nickname,
		Color:    color,
	}
	s.users[user.ID] = user
	return user
}

// AddMessage appends a message from userID to the history and delivers it
// to every active subscriber. Returns ErrUnknownUser if userID never joined.
func (s *Store) AddMessage(ctx context.Context, userID, text string) (Message, error) {
	s.mu.RLock()
	user, ok := s.users[userID]
	s.mu.RUnlock()
	if !ok {
		return Message{}, ErrUnknownUser
	}

	msg := Message{User: user, Text: text, SentAt: time.Now()}
	// Persist before broadcasting: a subscriber must never see a live event
	// for a message that then turns out to be missing from history on
	// reconnect.
	if err := s.history.Append(ctx, msg); err != nil {
		return Message{}, fmt.Errorf("append message to history: %w", err)
	}

	s.mu.RLock()
	recipients := s.recipientsLocked()
	s.mu.RUnlock()

	broadcast(recipients, Event{Message: &msg})
	return msg, nil
}

// Subscribe registers a new listener, bumps the online count (and notifies
// every subscriber, including this one, of the new count), and returns:
// the full message history so far, a channel that receives every event
// from this point on, and an unsubscribe func that must be called once the
// caller stops reading.
func (s *Store) Subscribe(ctx context.Context) ([]Message, <-chan Event, func(), error) {
	history, err := s.history.All(ctx)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("load message history: %w", err)
	}

	s.mu.Lock()
	id := newID()
	ch := make(chan Event, 32)
	s.subscribers[id] = ch
	recipients := s.recipientsLocked()
	count := len(s.subscribers)
	s.mu.Unlock()

	broadcast(recipients, Event{OnlineCount: &count})

	unsubscribe := func() {
		s.mu.Lock()
		if ch, ok := s.subscribers[id]; ok {
			delete(s.subscribers, id)
			close(ch)
		}
		recipients := s.recipientsLocked()
		count := len(s.subscribers)
		s.mu.Unlock()

		broadcast(recipients, Event{OnlineCount: &count})
	}
	return history, ch, unsubscribe, nil
}

// OnlineCount returns the current number of active Subscribe streams. Safe
// for concurrent use, including from an HTTP handler outside the gRPC layer.
func (s *Store) OnlineCount() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.subscribers)
}

// recipientsLocked snapshots the current subscriber channels. Callers must
// hold s.mu (for reading or writing).
func (s *Store) recipientsLocked() []chan Event {
	recipients := make([]chan Event, 0, len(s.subscribers))
	for _, ch := range s.subscribers {
		recipients = append(recipients, ch)
	}
	return recipients
}

// broadcast delivers ev to every recipient without blocking. A subscriber
// whose buffered channel is full (i.e. its Subscribe stream is stalled, e.g.
// a slow or dead network peer) has the event dropped for it rather than
// stalling AddMessage/Subscribe for every other caller — the dropped
// subscriber still recovers on reconnect via history replay.
func broadcast(recipients []chan Event, ev Event) {
	for _, ch := range recipients {
		select {
		case ch <- ev:
		default:
		}
	}
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand.Read only fails if the OS RNG is broken
	}
	return hex.EncodeToString(b)
}
