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
	lis := bufconn.Listen(1024 * 1024)
	t.Cleanup(func() { lis.Close() })

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
	t.Cleanup(func() { conn.Close() })

	client := chatv1.NewChatServiceClient(conn)

	joinResp, err := client.Join(ctx, &chatv1.JoinRequest{Nickname: "Mara"})
	if err != nil {
		t.Fatalf("Join: %v", err)
	}
	if joinResp.GetUser().GetColor() == "" {
		t.Fatal("expected a color to be assigned")
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

	if _, err := client.SendMessage(ctx, &chatv1.SendMessageRequest{
		UserId: joinResp.GetUser().GetId(),
		Text:   "hallo",
	}); err != nil {
		t.Fatalf("SendMessage: %v", err)
	}

	msg, err := stream.Recv()
	if err != nil {
		t.Fatalf("stream.Recv: %v", err)
	}
	if msg.GetText() != "hallo" {
		t.Fatalf("expected text %q, got %q", "hallo", msg.GetText())
	}
	if msg.GetUser().GetNickname() != "Mara" {
		t.Fatalf("expected nickname %q, got %q", "Mara", msg.GetUser().GetNickname())
	}
}
