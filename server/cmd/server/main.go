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

	"github.com/redis/go-redis/v9"
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
	redisPingTimeout  = 5 * time.Second
	redisHistoryKey   = "chat:history"
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

	history, closeHistory, err := newHistory(logger)
	if err != nil {
		return err
	}
	defer closeHistory()

	st := store.New(history)
	healthServer := health.NewServer()
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(recoveryUnaryInterceptor(logger), loggingUnaryInterceptor(logger)),
		grpc.ChainStreamInterceptor(recoveryStreamInterceptor(logger), loggingStreamInterceptor(logger)),
	)
	chatv1.RegisterChatServiceServer(grpcServer, chatserver.New(st))
	healthv1.RegisterHealthServer(grpcServer, healthServer)
	healthServer.SetServingStatus("", healthv1.HealthCheckResponse_SERVING)

	// GET /api/online-count is the one plain-HTTP, publicly reachable route
	// alongside the gRPC service — served over Coolify's regular HTTPS proxy
	// domain (unlike the chat traffic itself, which bypasses that proxy on a
	// separate plaintext port, see application/lib/config/server_config.dart).
	// Everything else falls through to the gRPC handler.
	mux := http.NewServeMux()
	mux.HandleFunc("GET /api/online-count", onlineCountHandler(st))
	mux.Handle("/", grpcServer)

	// Reverse proxies in front of this server (Coolify/Traefik) may speak
	// plain HTTP/1.1 to the backend rather than HTTP/2. Serving h2c
	// (unencrypted HTTP/2) alongside HTTP/1.1 lets it accept both, instead
	// of refusing HTTP/1.1 connections outright.
	protocols := new(http.Protocols)
	protocols.SetHTTP1(true)
	protocols.SetUnencryptedHTTP2(true)
	httpServer := &http.Server{
		Addr:              addr,
		Handler:           mux,
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

// newHistory sets up the message history backend from CHAT_REDIS_URL. If
// unset, message history falls back to in-memory (fine for local dev, but
// lost on every restart). If set, connectivity is verified once here and
// failure aborts startup — a configured-but-unreachable Redis fails loud
// rather than silently degrading to in-memory and quietly losing history.
// The returned closer must be deferred by the caller.
func newHistory(logger *slog.Logger) (store.History, func(), error) {
	redisURL := os.Getenv("CHAT_REDIS_URL")
	if redisURL == "" {
		logger.Warn("CHAT_REDIS_URL not set, message history is in-memory only and will not survive restarts")
		return store.NewMemoryHistory(), func() {}, nil
	}

	opts, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, nil, fmt.Errorf("parse CHAT_REDIS_URL: %w", err)
	}
	client := redis.NewClient(opts)

	pingCtx, cancel := context.WithTimeout(context.Background(), redisPingTimeout)
	defer cancel()
	if err := client.Ping(pingCtx).Err(); err != nil {
		_ = client.Close()
		return nil, nil, fmt.Errorf("connect to redis: %w", err)
	}

	logger.Info("message history backed by redis")
	return store.NewRedisHistory(client, redisHistoryKey), func() { _ = client.Close() }, nil
}
