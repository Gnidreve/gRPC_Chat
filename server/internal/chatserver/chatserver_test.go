package chatserver

import (
	"context"
	"net"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	chatv1 "github.com/Gnidreve/gRPC_Chat/server/gen/chat/v1"
	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

func TestJoinSendSubscribe(t *testing.T) {
	t.Parallel()

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { _ = lis.Close() })

	grpcServer := grpc.NewServer()
	chatv1.RegisterChatServiceServer(grpcServer, New(store.New()))
	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(grpcServer.Stop)

	ctx := context.Background()
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

	client := chatv1.NewChatServiceClient(conn)

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}
	if joinResp.GetUser().GetColor() != "#12B76A" {
		t.Fatalf("expected color %q, got %q", "#12B76A", joinResp.GetUser().GetColor())
	}

	subCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	stream, err := client.Subscribe(subCtx, &chatv1.SubscribeRequest{UserId: joinResp.GetUser().GetId()})
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
}

func TestSubscribeReplaysHistoryBeforeLiveEvents(t *testing.T) {
	t.Parallel()

	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { _ = lis.Close() })

	grpcServer := grpc.NewServer()
	chatv1.RegisterChatServiceServer(grpcServer, New(store.New()))
	go func() { _ = grpcServer.Serve(lis) }()
	t.Cleanup(grpcServer.Stop)

	ctx := context.Background()
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

	client := chatv1.NewChatServiceClient(conn)

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara", Color: "#12B76A"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}

	if _, err := client.SendMessage(ctx, &chatv1.SendMessageRequest{
		UserId: joinResp.GetUser().GetId(),
		Text:   "before subscribing",
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	subCtx, cancel := context.WithCancel(ctx)
	t.Cleanup(cancel)
	stream, err := client.Subscribe(subCtx, &chatv1.SubscribeRequest{UserId: joinResp.GetUser().GetId()})
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
