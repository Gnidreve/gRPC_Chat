package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	chatv1 "github.com/Gnidreve/gRPC_Chat/server/gen/chat/v1"
	"github.com/Gnidreve/gRPC_Chat/server/internal/chatserver"
	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

func main() {
	addr := os.Getenv("CHAT_SERVER_ADDR")
	if addr == "" {
		addr = ":3000" // matches the Coolify deployment's configured exposed port
	}

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen on %s: %v", addr, err)
	}

	grpcServer := grpc.NewServer()
	chatv1.RegisterChatServiceServer(grpcServer, chatserver.New(store.New()))

	log.Printf("chat server listening on %s", addr)
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
