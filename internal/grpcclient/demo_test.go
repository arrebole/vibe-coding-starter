package grpcclient

import (
	"context"
	"errors"
	"log/slog"
	"net"
	"testing"
	"time"

	"github.com/arrebole/vibe-coding-starter/internal/config"
	demov1 "github.com/arrebole/vibe-coding-starter/proto/external/demo/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type demoServiceServer struct {
	demov1.UnimplementedDemoServiceServer
	ping func(context.Context, *demov1.PingRequest) (*demov1.PingResponse, error)
}

func (s *demoServiceServer) Ping(ctx context.Context, request *demov1.PingRequest) (*demov1.PingResponse, error) {
	return s.ping(ctx, request)
}

func TestDemoClientNotConfigured(t *testing.T) {
	client := NewDemoClient(config.GRPCClientConfig{Timeout: time.Second}, slog.Default())
	defer client.Close()

	_, err := client.Ping(context.Background(), "todo")
	if !errors.Is(err, ErrDemoClientNotConfigured) {
		t.Fatalf("Ping() error = %v, want %v", err, ErrDemoClientNotConfigured)
	}
}

func TestDemoClientPing(t *testing.T) {
	addr := startDemoService(t, func(_ context.Context, request *demov1.PingRequest) (*demov1.PingResponse, error) {
		if request.GetName() != "todo" {
			t.Errorf("Ping() name = %q, want %q", request.GetName(), "todo")
		}
		return &demov1.PingResponse{Message: "pong: " + request.GetName()}, nil
	})
	client := NewDemoClient(config.GRPCClientConfig{Addr: addr, Timeout: time.Second}, slog.Default())
	t.Cleanup(func() { _ = client.Close() })

	message, err := client.Ping(context.Background(), "todo")
	if err != nil {
		t.Fatalf("Ping() error = %v", err)
	}
	if message != "pong: todo" {
		t.Fatalf("Ping() = %q, want %q", message, "pong: todo")
	}
}

func TestDemoClientPingPropagatesRemoteError(t *testing.T) {
	addr := startDemoService(t, func(context.Context, *demov1.PingRequest) (*demov1.PingResponse, error) {
		return nil, status.Error(codes.PermissionDenied, "denied")
	})
	client := NewDemoClient(config.GRPCClientConfig{Addr: addr, Timeout: time.Second}, slog.Default())
	t.Cleanup(func() { _ = client.Close() })

	_, err := client.Ping(context.Background(), "todo")
	if status.Code(err) != codes.PermissionDenied {
		t.Fatalf("Ping() code = %v, want %v", status.Code(err), codes.PermissionDenied)
	}
}

func TestDemoClientPingTimesOut(t *testing.T) {
	addr := startDemoService(t, func(ctx context.Context, _ *demov1.PingRequest) (*demov1.PingResponse, error) {
		<-ctx.Done()
		return nil, status.FromContextError(ctx.Err()).Err()
	})
	client := NewDemoClient(config.GRPCClientConfig{Addr: addr, Timeout: 20 * time.Millisecond}, slog.Default())
	t.Cleanup(func() { _ = client.Close() })

	_, err := client.Ping(context.Background(), "todo")
	if status.Code(err) != codes.DeadlineExceeded {
		t.Fatalf("Ping() code = %v, want %v", status.Code(err), codes.DeadlineExceeded)
	}
}

func startDemoService(
	t *testing.T,
	ping func(context.Context, *demov1.PingRequest) (*demov1.PingResponse, error),
) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("listen: %v", err)
	}
	server := grpc.NewServer()
	demov1.RegisterDemoServiceServer(server, &demoServiceServer{ping: ping})
	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})
	return listener.Addr().String()
}
