// Package chatserver implements the chat.v1.ChatService gRPC contract on
// top of a store.Store.
package chatserver

import (
	"context"
	"errors"
	"math"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	chatv1 "github.com/Gnidreve/gRPC_Chat/server/gen/chat/v1"
	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

// Server implements the chat.v1.ChatService gRPC contract on top of a
// store.Store.
type Server struct {
	chatv1.UnimplementedChatServiceServer

	store *store.Store
}

// New returns a Server backed by the given store.
func New(s *store.Store) *Server {
	return &Server{store: s}
}

// Join registers or re-registers a user with the store.
func (s *Server) Join(_ context.Context, req *chatv1.JoinRequest) (*chatv1.JoinResponse, error) {
	user := s.store.Join(req.GetId(), req.GetNickname(), req.GetColor())
	return &chatv1.JoinResponse{User: toProtoUser(user)}, nil
}

// SendMessage appends a message from the requesting user and broadcasts it
// to subscribers. The sender's location is required — it's stored with the
// message so other participants' distance to it can be computed later.
func (s *Server) SendMessage(ctx context.Context, req *chatv1.SendMessageRequest) (*chatv1.SendMessageResponse, error) {
	if req.Lat == nil || req.Lng == nil {
		return nil, status.Error(codes.InvalidArgument, "location is required")
	}
	msg, err := s.store.AddMessage(ctx, req.GetUserId(), req.GetText(), *req.Lat, *req.Lng)
	if err != nil {
		if errors.Is(err, store.ErrUnknownUser) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Unavailable, "failed to store message")
	}
	// This echoes the sender's own message back to them, so no distance is
	// computed (toProtoMessage never sets DistanceKm).
	return &chatv1.SendMessageResponse{Message: toProtoMessage(msg)}, nil
}

// Subscribe streams message history followed by live events to the caller.
// The caller's location is required and fixed for the lifetime of this
// call: every DistanceKm on this stream, for history and live events
// alike, is computed against it (see toProtoMessageWithDistance).
func (s *Server) Subscribe(req *chatv1.SubscribeRequest, stream chatv1.ChatService_SubscribeServer) error {
	if req.Lat == nil || req.Lng == nil {
		return status.Error(codes.InvalidArgument, "location is required")
	}
	requesterID, refLat, refLng := req.GetUserId(), *req.Lat, *req.Lng

	history, events, unsubscribe, err := s.store.Subscribe(stream.Context())
	if err != nil {
		return status.Error(codes.Unavailable, "failed to load message history")
	}
	defer unsubscribe()

	// Send headers now that the subscription is registered with the store,
	// so a client that waits for them (via stream.Header()) is guaranteed
	// not to miss an event it triggers right after Subscribe returns.
	if err := stream.SendHeader(metadata.MD{}); err != nil {
		return err
	}

	for _, msg := range history {
		if err := stream.Send(messageEvent(msg, requesterID, refLat, refLng)); err != nil {
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
			if err := stream.Send(toProtoEvent(ev, requesterID, refLat, refLng)); err != nil {
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

// toProtoMessageWithDistance is toProtoMessage plus DistanceKm: the
// Haversine distance in kilometers, rounded, from (refLat, refLng) to
// where m was sent. Omitted for the requester's own messages and for
// messages with no stored location (m.HasLocation() false).
func toProtoMessageWithDistance(m store.Message, requesterID string, refLat, refLng float64) *chatv1.ChatMessage {
	pm := toProtoMessage(m)
	if m.User.ID != requesterID && m.HasLocation() {
		km := int32(math.Round(haversineKm(refLat, refLng, m.Lat, m.Lng)))
		pm.DistanceKm = &km
	}
	return pm
}

func messageEvent(m store.Message, requesterID string, refLat, refLng float64) *chatv1.ChatEvent {
	return &chatv1.ChatEvent{
		Event: &chatv1.ChatEvent_Message{Message: toProtoMessageWithDistance(m, requesterID, refLat, refLng)},
	}
}

func toProtoEvent(ev store.Event, requesterID string, refLat, refLng float64) *chatv1.ChatEvent {
	if ev.Message != nil {
		return messageEvent(*ev.Message, requesterID, refLat, refLng)
	}
	return &chatv1.ChatEvent{
		Event: &chatv1.ChatEvent_Presence{
			Presence: &chatv1.PresenceUpdate{OnlineCount: clampToInt32(*ev.OnlineCount)},
		},
	}
}

// clampToInt32 converts n to int32, clamping instead of silently wrapping if
// it ever exceeds the range (subscriber counts stay tiny in practice, but a
// silent wraparound to a negative online count would be a confusing bug).
func clampToInt32(n int) int32 {
	if n > math.MaxInt32 {
		return math.MaxInt32
	}
	return int32(n) //nolint:gosec // bounds-checked above; gosec's G115 can't see the guard
}
