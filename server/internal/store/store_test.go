package store

import "testing"

func TestJoinKeepsClientChosenNicknameAndColor(t *testing.T) {
	s := New()

	a := s.Join("Mara", "#12B76A")
	b := s.Join("Jonas", "#F04438")

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

func TestAddMessageUnknownUser(t *testing.T) {
	s := New()

	if _, err := s.AddMessage("does-not-exist", "hi"); err != ErrUnknownUser {
		t.Fatalf("expected ErrUnknownUser, got %v", err)
	}
}

func TestSubscribeReceivesMessage(t *testing.T) {
	s := New()
	user := s.Join("Mara", "#12B76A")

	messages, unsubscribe := s.Subscribe()
	defer unsubscribe()

	if _, err := s.AddMessage(user.ID, "hey"); err != nil {
		t.Fatalf("AddMessage: %v", err)
	}

	select {
	case msg := <-messages:
		if msg.Text != "hey" {
			t.Fatalf("expected text %q, got %q", "hey", msg.Text)
		}
		if msg.User.ID != user.ID {
			t.Fatalf("expected sender %q, got %q", user.ID, msg.User.ID)
		}
	default:
		t.Fatal("expected a message to be waiting on the subscriber channel")
	}
}
