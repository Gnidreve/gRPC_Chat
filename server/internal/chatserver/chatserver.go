// Package chatserver implements the chat.v1.ChatService gRPC contract on
// top of an in-memory store.Store.
package chatserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatv1 "github.com/Gnidreve/gRPC_Chat/server/gen/chat/v1"
	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

type Server struct {
	chatv1.UnimplementedChatServiceServer

	store *store.Store
}

func New(s *store.Store) *Server {
	return &Server{store: s}
}

func (s *Server) Join(_ context.Context, req *chatv1.JoinRequest) (*chatv1.JoinResponse, error) {
	user := s.store.Join(req.GetId(), req.GetNickname(), req.GetColor())
	return &chatv1.JoinResponse{User: toProtoUser(user)}, nil
}

func (s *Server) SendMessage(_ context.Context, req *chatv1.SendMessageRequest) (*chatv1.SendMessageResponse, error) {
	msg, err := s.store.AddMessage(req.GetUserId(), req.GetText())
	if err != nil {
		return nil, status.Error(codes.NotFound, err.Error())
	}
	return &chatv1.SendMessageResponse{Message: toProtoMessage(msg)}, nil
}

func (s *Server) Subscribe(_ *chatv1.SubscribeRequest, stream chatv1.ChatService_SubscribeServer) error {
	history, events, unsubscribe := s.store.Subscribe()
	defer unsubscribe()

	// Send headers now that the subscription is registered with the store,
	// so a client that waits for them (via stream.Header()) is guaranteed
	// not to miss an event it triggers right after Subscribe returns.
	if err := stream.SendHeader(metadata.MD{}); err != nil {
		return err
	}

	for _, msg := range history {
		if err := stream.Send(messageEvent(msg)); err != nil {
			return err
		}
	}

	for {
		select {
		case <-stream.Context().Done():
			return stream.Context().Err()
		case ev, ok := <-events:
			if !ok {
				return nil
			}
			if err := stream.Send(toProtoEvent(ev)); err != nil {
				return err
			}
		}
	}
}

func toProtoUser(u store.User) *chatv1.User {
	return &chatv1.User{
		Id:       u.ID,
		Nickname: u.Nickname,
		Color:    u.Color,
	}
}

func toProtoMessage(m store.Message) *chatv1.ChatMessage {
	return &chatv1.ChatMessage{
		User:   toProtoUser(m.User),
		Text:   m.Text,
		SentAt: timestamppb.New(m.SentAt),
	}
}

func messageEvent(m store.Message) *chatv1.ChatEvent {
	return &chatv1.ChatEvent{
		Event: &chatv1.ChatEvent_Message{Message: toProtoMessage(m)},
	}
}

func toProtoEvent(ev store.Event) *chatv1.ChatEvent {
	if ev.Message != nil {
		return messageEvent(*ev.Message)
	}
	return &chatv1.ChatEvent{
		Event: &chatv1.ChatEvent_Presence{
			Presence: &chatv1.PresenceUpdate{OnlineCount: int32(*ev.OnlineCount)},
		},
	}
}
