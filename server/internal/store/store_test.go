package store

import "testing"

func TestJoinKeepsClientChosenNicknameAndColor(t *testing.T) {
	s := New()

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
	s := New()

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
	s := New()

	if _, err := s.AddMessage("does-not-exist", "hi"); err != ErrUnknownUser {
		t.Fatalf("expected ErrUnknownUser, got %v", err)
	}
}

func TestSubscribeReceivesMessage(t *testing.T) {
	s := New()
	user := s.Join("", "Mara", "#12B76A")

	_, events, unsubscribe := s.Subscribe()
	defer unsubscribe()

	// The subscriber's own join broadcasts a presence event; drain it.
	<-events

	if _, err := s.AddMessage(user.ID, "hey"); err != nil {
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
	s := New()
	user := s.Join("", "Mara", "#12B76A")

	if _, err := s.AddMessage(user.ID, "before subscribing"); err != nil {
		t.Fatalf("AddMessage: %v", err)
	}

	history, _, unsubscribe := s.Subscribe()
	defer unsubscribe()

	if len(history) != 1 || history[0].Text != "before subscribing" {
		t.Fatalf("expected history with 1 message, got %+v", history)
	}
}

func TestSubscribePresenceCount(t *testing.T) {
	s := New()

	_, eventsA, unsubscribeA := s.Subscribe()
	defer unsubscribeA()

	evA := <-eventsA
	if evA.OnlineCount == nil || *evA.OnlineCount != 1 {
		t.Fatalf("expected online count 1, got %+v", evA)
	}

	_, eventsB, unsubscribeB := s.Subscribe()

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
