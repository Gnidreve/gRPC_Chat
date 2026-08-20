package store

import "testing"

func TestJoinAssignsColor(t *testing.T) {
	s := New()

	a := s.Join("Mara")
	b := s.Join("Jonas")

	if a.Color == "" || b.Color == "" {
		t.Fatalf("expected colors to be assigned, got %q and %q", a.Color, b.Color)
	}
	if a.Color == b.Color {
		t.Fatalf("expected distinct colors, both got %q", a.Color)
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
	user := s.Join("Mara")

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
