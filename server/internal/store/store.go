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

// colors is the fixed palette users are assigned from, round-robin, on Join.
var colors = []string{
	"#12B76A", "#F04438", "#F79009", "#7A5AF8",
	"#0BA5EC", "#EE46BC", "#84CC16", "#F97066",
}

// User is a chat participant: a server-assigned id and color, plus a
// nickname the user picked themselves.
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

// Store is the in-memory chat room: known users, message history, and the
// live subscribers currently streaming new messages.
type Store struct {
	mu          sync.RWMutex
	users       map[string]User
	messages    []Message
	subscribers map[string]chan Message
	nextColor   int
}

func New() *Store {
	return &Store{
		users:       make(map[string]User),
		subscribers: make(map[string]chan Message),
	}
}

// Join registers a new user with the given nickname and assigns them the
// next color from the palette.
func (s *Store) Join(nickname string) User {
	s.mu.Lock()
	defer s.mu.Unlock()

	user := User{
		ID:       newID(),
		Nickname: nickname,
		Color:    colors[s.nextColor%len(colors)],
	}
	s.nextColor++
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

	recipients := make([]chan Message, 0, len(s.subscribers))
	for _, ch := range s.subscribers {
		recipients = append(recipients, ch)
	}
	s.mu.Unlock()

	for _, ch := range recipients {
		ch <- msg
	}
	return msg, nil
}

// Subscribe registers a new listener and returns a channel that receives
// every message posted from this point on, plus an unsubscribe func that
// must be called once the caller stops reading.
func (s *Store) Subscribe() (<-chan Message, func()) {
	s.mu.Lock()
	defer s.mu.Unlock()

	id := newID()
	ch := make(chan Message, 16)
	s.subscribers[id] = ch

	unsubscribe := func() {
		s.mu.Lock()
		defer s.mu.Unlock()
		if ch, ok := s.subscribers[id]; ok {
			delete(s.subscribers, id)
			close(ch)
		}
	}
	return ch, unsubscribe
}

func newID() string {
	b := make([]byte, 16)
	if _, err := rand.Read(b); err != nil {
		panic(err) // crypto/rand.Read only fails if the OS RNG is broken
	}
	return hex.EncodeToString(b)
}
