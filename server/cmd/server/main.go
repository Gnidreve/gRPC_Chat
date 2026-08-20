package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
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

	// Reverse proxies in front of this server (Coolify/Traefik) may speak
	// plain HTTP/1.1 to the backend rather than HTTP/2. Serving the gRPC
	// server as h2c lets it accept both, instead of refusing HTTP/1.1
	// connections outright.
	h2Server := &http2.Server{}
	httpServer := &http.Server{
		Addr:    addr,
		Handler: h2c.NewHandler(grpcServer, h2Server),
	}

	log.Printf("chat server listening on %s", addr)
	if err := httpServer.Serve(lis); err != nil {
		log.Fatalf("serve: %v", err)
	}
}
