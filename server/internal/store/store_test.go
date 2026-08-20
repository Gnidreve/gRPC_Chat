package store

import (
	"context"
	"errors"
	"testing"
)

func newTestStore(t *testing.T) *Store {
	t.Helper()
	return New(NewMemoryHistory())
}

func TestJoinKeepsClientChosenNicknameAndColor(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)

	a := s.Join("", "Mara", "#12B76A")
	b := s.Join("", "Jonas", "#F04438")

	if a.Nickname != "Mara" || a.Color != "#12B76A" {
		t.Fatalf("expected Mara/#12B76A, got %q/%q", a.Nickname, a.Color)
	}
	if b.Nickname != "Jonas" || b.Color != "#F04438" {
		t.Fatalf("expected Jonas/#F04438, got %q/%q", b.Nickname, b.Color)
	}
	if a.ID == b.ID {
		t.Fatalf("expected distinct ids, both got %q", a.ID)
	}
}

func TestJoinReusesClientProvidedID(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)

	first := s.Join("device-abc", "Mara", "#12B76A")
	if first.ID != "device-abc" {
		t.Fatalf("expected id %q, got %q", "device-abc", first.ID)
	}

	// Simulate an app restart: same client id, Join called again.
	second := s.Join("device-abc", "Mara", "#12B76A")
	if second.ID != first.ID {
		t.Fatalf("expected the same id across rejoins, got %q then %q", first.ID, second.ID)
	}
}

func TestAddMessageUnknownUser(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)

	if _, err := s.AddMessage(context.Background(), "does-not-exist", "hi", 0, 0); !errors.Is(err, ErrUnknownUser) {
		t.Fatalf("expected ErrUnknownUser, got %v", err)
	}
}

func TestSubscribeReceivesMessage(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	user := s.Join("", "Mara", "#12B76A")

	_, events, unsubscribe, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer unsubscribe()

	// The subscriber's own join broadcasts a presence event; drain it.
	<-events

	if _, err := s.AddMessage(context.Background(), user.ID, "hey", 52.52, 13.405); err != nil {
		t.Fatalf("AddMessage: %v", err)
	}

	ev := <-events
	if ev.Message == nil {
		t.Fatalf("expected a message event, got %+v", ev)
	}
	if ev.Message.Text != "hey" {
		t.Fatalf("expected text %q, got %q", "hey", ev.Message.Text)
	}
	if ev.Message.User.ID != user.ID {
		t.Fatalf("expected sender %q, got %q", user.ID, ev.Message.User.ID)
	}
}

func TestSubscribeReplaysHistory(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	user := s.Join("", "Mara", "#12B76A")

	if _, err := s.AddMessage(context.Background(), user.ID, "before subscribing", 52.52, 13.405); err != nil {
		t.Fatalf("AddMessage: %v", err)
	}

	history, _, unsubscribe, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer unsubscribe()

	if len(history) != 1 || history[0].Text != "before subscribing" {
		t.Fatalf("expected history with 1 message, got %+v", history)
	}
}

func TestSubscribePresenceCount(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)

	_, eventsA, unsubscribeA, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer unsubscribeA()

	evA := <-eventsA
	if evA.OnlineCount == nil || *evA.OnlineCount != 1 {
		t.Fatalf("expected online count 1, got %+v", evA)
	}

	_, eventsB, unsubscribeB, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	// Both subscribers should be told the count is now 2.
	evA2 := <-eventsA
	if evA2.OnlineCount == nil || *evA2.OnlineCount != 2 {
		t.Fatalf("expected online count 2, got %+v", evA2)
	}
	evB := <-eventsB
	if evB.OnlineCount == nil || *evB.OnlineCount != 2 {
		t.Fatalf("expected online count 2, got %+v", evB)
	}

	unsubscribeB()

	evA3 := <-eventsA
	if evA3.OnlineCount == nil || *evA3.OnlineCount != 1 {
		t.Fatalf("expected online count 1 after unsubscribe, got %+v", evA3)
	}
}

func TestOnlineCount(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)

	if got := s.OnlineCount(); got != 0 {
		t.Fatalf("expected online count 0 before any subscriber, got %d", got)
	}

	_, _, unsubscribeA, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	_, _, unsubscribeB, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}

	if got := s.OnlineCount(); got != 2 {
		t.Fatalf("expected online count 2, got %d", got)
	}

	unsubscribeA()
	if got := s.OnlineCount(); got != 1 {
		t.Fatalf("expected online count 1 after one unsubscribe, got %d", got)
	}

	unsubscribeB()
	if got := s.OnlineCount(); got != 0 {
		t.Fatalf("expected online count 0 after all unsubscribed, got %d", got)
	}
}

func TestAddMessageStoresLocation(t *testing.T) {
	t.Parallel()

	s := newTestStore(t)
	user := s.Join("", "Mara", "#12B76A")

	msg, err := s.AddMessage(context.Background(), user.ID, "hey", 52.52, 13.405)
	if err != nil {
		t.Fatalf("AddMessage: %v", err)
	}
	if msg.Lat != 52.52 || msg.Lng != 13.405 {
		t.Fatalf("expected lat/lng 52.52/13.405, got %v/%v", msg.Lat, msg.Lng)
	}
	if !msg.HasLocation() {
		t.Fatalf("expected HasLocation true for a message with real coordinates")
	}

	history, _, unsubscribe, err := s.Subscribe(context.Background())
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	defer unsubscribe()

	if len(history) != 1 || history[0].Lat != 52.52 || history[0].Lng != 13.405 {
		t.Fatalf("expected replayed history to keep lat/lng, got %+v", history)
	}
}

func TestMessageHasLocation(t *testing.T) {
	t.Parallel()

	if (Message{}).HasLocation() {
		t.Fatalf("expected zero-value Message to have no location")
	}
	if !(Message{Lat: 52.52}).HasLocation() {
		t.Fatalf("expected a message with only Lat set to have a location")
	}
	if !(Message{Lng: 13.405}).HasLocation() {
		t.Fatalf("expected a message with only Lng set to have a location")
	}
}
