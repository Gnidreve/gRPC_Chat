// Command server runs the chat.v1.ChatService gRPC server.
package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	healthv1 "google.golang.org/grpc/health/grpc_health_v1"

	chatv1 "github.com/Gnidreve/gRPC_Chat/server/gen/chat/v1"
	"github.com/Gnidreve/gRPC_Chat/server/internal/chatserver"
	"github.com/Gnidreve/gRPC_Chat/server/internal/store"
)

const (
	readHeaderTimeout = 5 * time.Second
	shutdownTimeout   = 15 * time.Second
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	if err := run(logger); err != nil {
		logger.Error("server exited with error", "error", err)
		os.Exit(1)
	}
}

func run(logger *slog.Logger) error {
	addr := os.Getenv("CHAT_SERVER_ADDR")
	if addr == "" {
		addr = ":3000" // matches the Coolify deployment's configured exposed port
	}

	var lc net.ListenConfig
	lis, err := lc.Listen(context.Background(), "tcp", addr)
	if err != nil {
		return fmt.Errorf("listen on %s: %w", addr, err)
	}

	healthServer := health.NewServer()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(recoveryUnaryInterceptor(logger), loggingUnaryInterceptor(logger)),
		grpc.ChainStreamInterceptor(recoveryStreamInterceptor(logger), loggingStreamInterceptor(logger)),
	)
	chatv1.RegisterChatServiceServer(grpcServer, chatserver.New(store.New()))
	healthv1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)

	// Reverse proxies in front of this server (Coolify/Traefik) may speak
	// plain HTTP/1.1 to the backend rather than HTTP/2. Serving h2c
	// (unencrypted HTTP/2) alongside HTTP/1.1 lets it accept both, instead
	// of refusing HTTP/1.1 connections outright.
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           grpcServer,
		Protocols:         protocols,
		ReadHeaderTimeout: readHeaderTimeout,
	}

	serveErr := make(chan error, 1)
	go func() {
		logger.Info("chat server listening", "addr", addr)
		serveErr <- httpServer.Serve(lis)
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	select {
	case err := <-serveErr:
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			return err
		}
	case <-ctx.Done():
		logger.Info("shutdown signal received, draining connections")
		healthServer.SetServingStatus("", healthv1.HealthCheckResponse_NOT_SERVING)

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()
		if err := httpServer.Shutdown(shutdownCtx); err != nil {
			logger.Warn("graceful shutdown timed out, forcing close", "error", err)
			_ = httpServer.Close()
		}
	}

	logger.Info("chat server stopped")
	return nil
}
