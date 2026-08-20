package chatserver

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"google.golang.org/grpc/test/bufconn"
	"google.golang.org/protobuf/proto"

	chatv1 "github.com/Gnidreve/gRPC_Chat/server/gen/chat/v1"
	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

// berlin and munich are real, well-separated coordinates used across these
// tests as stand-ins for two different users' locations.
var (
	berlin = struct{ lat, lng float64 }{52.5200, 13.4050}
	munich = struct{ lat, lng float64 }{48.1351, 11.5820}
)

func newTestClient(t *testing.T) chatv1.ChatServiceClient {
	t.Helper()

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { _ = lis.Close() })

	grpcServer := grpc.NewServer()
	chatv1.RegisterChatServiceServer(grpcServer, New(store.New(store.NewMemoryHistory())))
	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(grpcServer.Stop)

	conn, err := grpc.NewClient("passthrough:///bufnet",
		grpc.WithContextDialer(func(ctx context.Context, _ string) (net.Conn, error) {
			return lis.DialContext(ctx)
		}),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	t.Cleanup(func() { _ = conn.Close() })

	return chatv1.NewChatServiceClient(conn)
}

func TestJoinSendSubscribe(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}
	if joinResp.GetUser().GetColor() != "#12B76A" {
		t.Fatalf("expected color %q, got %q", "#12B76A", joinResp.GetUser().GetColor())
	}

	subCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	stream, err := client.Subscribe(subCtx, &chatv1.SubscribeRequest{
		UserId: joinResp.GetUser().GetId(),
		Lat:    proto.Float64(berlin.lat),
		Lng:    proto.Float64(berlin.lng),
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	// Wait for the server to confirm the subscription is registered before
	// sending, otherwise the message could race the store registration.
	if _, err := stream.Header(); err != nil {
		t.Fatalf("stream.Header: %v", err)
	}

	// No history yet, so the first event is this subscriber's own presence
	// update from joining.
	presenceEv, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv (presence): %v", err)
	}
	if presenceEv.GetPresence().GetOnlineCount() != 1 {
		t.Fatalf("expected online count 1, got %+v", presenceEv)
	}

	if _, err := client.SendMessage(ctx, &chatv1.SendMessageRequest{
		UserId: joinResp.GetUser().GetId(),
		Text:   "hallo",
		Lat:    proto.Float64(berlin.lat),
		Lng:    proto.Float64(berlin.lng),
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	msgEv, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv (message): %v", err)
	}
	msg := msgEv.GetMessage()
	if msg.GetText() != "hallo" {
		t.Fatalf("expected text %q, got %q", "hallo", msg.GetText())
	}
	if msg.GetUser().GetNickname() != "Mara" {
		t.Fatalf("expected nickname %q, got %q", "Mara", msg.GetUser().GetNickname())
	}
	if msg.DistanceKm != nil {
		t.Fatalf("expected no distance on the sender's own message, got %d", msg.GetDistanceKm())
	}
}

func TestSubscribeReplaysHistoryBeforeLiveEvents(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}

	if _, err := client.SendMessage(ctx, &chatv1.SendMessageRequest{
		UserId: joinResp.GetUser().GetId(),
		Text:   "before subscribing",
		Lat:    proto.Float64(berlin.lat),
		Lng:    proto.Float64(berlin.lng),
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	subCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	stream, err := client.Subscribe(subCtx, &chatv1.SubscribeRequest{
		UserId: joinResp.GetUser().GetId(),
		Lat:    proto.Float64(berlin.lat),
		Lng:    proto.Float64(berlin.lng),
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if _, err := stream.Header(); err != nil {
		t.Fatalf("stream.Header: %v", err)
	}

	historyEv, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv (history): %v", err)
	}
	if historyEv.GetMessage().GetText() != "before subscribing" {
		t.Fatalf("expected replayed history message, got %+v", historyEv)
	}

	presenceEv, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv (presence): %v", err)
	}
	if presenceEv.GetPresence().GetOnlineCount() != 1 {
		t.Fatalf("expected online count 1, got %+v", presenceEv)
	}
}

func TestSubscribeComputesDistanceForOthersMessages(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()

	sender, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join (sender): %v", err)
	}
	recipient, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Jonas", Color: "#F04438"})
	if err != nil {
		t.Fatalf("Join (recipient): %v", err)
	}

	subCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	stream, err := client.Subscribe(subCtx, &chatv1.SubscribeRequest{
		UserId: recipient.GetUser().GetId(),
		Lat:    proto.Float64(munich.lat),
		Lng:    proto.Float64(munich.lng),
	})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if _, err := stream.Header(); err != nil {
		t.Fatalf("stream.Header: %v", err)
	}

	// Recipient's own join presence update.
	if _, err := stream.Recv(); err != nil {
		t.Fatalf("stream.Recv (presence 1): %v", err)
	}

	if _, err := client.SendMessage(ctx, &chatv1.SendMessageRequest{
		UserId: sender.GetUser().GetId(),
		Text:   "hallo aus Berlin",
		Lat:    proto.Float64(berlin.lat),
		Lng:    proto.Float64(berlin.lng),
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	// The sender only Joins (no Subscribe), and Join doesn't broadcast
	// presence, so the message is the very next event on the stream.
	msgEv, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv (message): %v", err)
	}
	msg := msgEv.GetMessage()
	if msg.DistanceKm == nil {
		t.Fatalf("expected a distance on another user's message, got none")
	}

	want := int32(haversineKm(munich.lat, munich.lng, berlin.lat, berlin.lng) + 0.5)
	if msg.GetDistanceKm() != want {
		t.Fatalf("expected distance_km %d (Berlin-Munich), got %d", want, msg.GetDistanceKm())
	}
	// Sanity-check it's in the right ballpark rather than trusting only the
	// self-referential computation above.
	if msg.GetDistanceKm() < 400 || msg.GetDistanceKm() > 600 {
		t.Fatalf("expected Berlin-Munich distance roughly 400-600km, got %d", msg.GetDistanceKm())
	}
}

func TestSubscribeRequiresLocation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}

	stream, err := client.Subscribe(ctx, &chatv1.SubscribeRequest{UserId: joinResp.GetUser().GetId()})
	if err != nil {
		t.Fatalf("Subscribe: %v", err)
	}
	if _, err := stream.Recv(); status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument for missing location, got %v", err)
	}
}

func TestSendMessageRequiresLocation(t *testing.T) {
	t.Parallel()

	client := newTestClient(t)
	ctx := context.Background()

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}

	_, err = client.SendMessage(ctx, &chatv1.SendMessageRequest{
		UserId: joinResp.GetUser().GetId(),
		Text:   "hallo",
	})
	if status.Code(err) != codes.InvalidArgument {
		t.Fatalf("expected InvalidArgument for missing location, got %v", err)
	}
}
