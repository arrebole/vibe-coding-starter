package grpcclient

import (
	"context"
	"errors"
	"log/slog"
	"testing"
	"time"

	"github.com/arrebole/vibe-coding-starter/internal/config"
)

func TestDemoClientNotConfigured(t *testing.T) {
	client := NewDemoClient(config.GRPCClientConfig{Timeout: time.Second}, slog.Default())
	defer client.Close()

	_, err := client.Ping(context.Background(), "todo")
	if !errors.Is(err, ErrDemoClientNotConfigured) {
		t.Fatalf("Ping() error = %v, want %v", err, ErrDemoClientNotConfigured)
	}
}
