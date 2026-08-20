// Package store holds the chat room state in memory: users and the
// message history. It is deliberately simple — a single mutex-guarded
// struct — so it can later be swapped for a Redis-backed implementation
// without touching the gRPC layer.
package store

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"sync"
	"time"
)

// ErrUnknownUser is returned when an operation references a user id that
// never joined (or was evicted).
var ErrUnknownUser = errors.New("store: unknown user")

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

// Store is the in-memory chat room: known users, message history, and the
// live subscribers currently streaming events.
type Store struct {
	mu          sync.RWMutex
	users       map[string]User
	messages    []Message
	subscribers map[string]chan Event
}

func New() *Store {
	return &Store{
		users:       make(map[string]User),
		subscribers: make(map[string]chan Event),
	}
}

// Join registers a new user with the given nickname and color, both
// chosen by the user themselves.
func (s *Store) Join(nickname, color string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:       newID(),
		Nickname: nickname,
		Color:    color,
	}
	s.users[user.ID] = user
	return user
}

// AddMessage appends a message from userID to the history and delivers it
// to every active subscriber. Returns ErrUnknownUser if userID never joined.
func (s *Store) AddMessage(userID, text string) (Message, error) {
	s.mu.Lock()
	user, ok := s.users[userID]
	if !ok {
		s.mu.Unlock()
		return Message{}, ErrUnknownUser
	}

	msg := Message{User: user, Text: text, SentAt: time.Now()}
	s.messages = append(s.messages, msg)
	recipients := s.recipientsLocked()
	s.mu.Unlock()

	broadcast(recipients, Event{Message: &msg})
	return msg, nil
}

// Subscribe registers a new listener, bumps the online count (and notifies
// every subscriber, including this one, of the new count), and returns:
// the full message history so far, a channel that receives every event
// from this point on, and an unsubscribe func that must be called once the
// caller stops reading.
func (s *Store) Subscribe() ([]Message, <-chan Event, func()) {
	s.mu.Lock()
	history := make([]Message, len(s.messages))
	copy(history, s.messages)

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
	return history, ch, unsubscribe
}

// recipientsLocked snapshots the current subscriber channels. Callers must
// hold s.mu.
func (s *Store) recipientsLocked() []chan Event {
	recipients := make([]chan Event, 0, len(s.subscribers))
	for _, ch := range s.subscribers {
		recipients = append(recipients, ch)
	}
	return recipients
}

func broadcast(recipients []chan Event, ev Event) {
	for _, ch := range recipients {
		ch <- ev
	}
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand.Read only fails if the OS RNG is broken
	}
	return hex.EncodeToString(b)
}
