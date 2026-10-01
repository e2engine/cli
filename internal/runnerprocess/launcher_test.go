package runnerprocess

import (
	"context"
	"errors"
	"net"
	"testing"
	"time"

	"github.com/e2engine/core/pkg/log"
)

func TestLauncher_waitReady_Ready(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("log.NewSilentLogger() error = %v", err)
	}

	launcher := NewLauncher(
		listener.Addr().String(),
		logger,
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	if err := launcher.waitReady(ctx); err != nil {
		t.Fatalf("waitReady() error = %v", err)
	}
}

func TestLauncher_waitReady_WaitsUntilReady(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}

	address := listener.Addr().String()

	if err := listener.Close(); err != nil {
		t.Fatalf("listener.Close() error = %v", err)
	}

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("log.NewSilentLogger() error = %v", err)
	}

	launcher := NewLauncher(address, logger)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		time.Second,
	)
	defer cancel()

	result := make(chan error, 1)

	go func() {
		result <- launcher.waitReady(ctx)
	}()

	time.Sleep(50 * time.Millisecond)

	listener, err = net.Listen("tcp", address)
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}
	defer func() {
		_ = listener.Close()
	}()

	select {
	case err := <-result:
		if err != nil {
			t.Fatalf("waitReady() error = %v", err)
		}

	case <-ctx.Done():
		t.Fatal("waitReady() did not detect listener readiness")
	}
}

func TestLauncher_waitReady_ContextCanceled(t *testing.T) {
	address := unusedAddress(t)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("log.NewSilentLogger() error = %v", err)
	}

	launcher := NewLauncher(address, logger)

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err = launcher.waitReady(ctx)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf(
			"waitReady() error = %v, want %v",
			err,
			context.Canceled,
		)
	}
}

func TestLauncher_waitReady_ContextDeadlineExceeded(t *testing.T) {
	address := unusedAddress(t)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("log.NewSilentLogger() error = %v", err)
	}

	launcher := NewLauncher(address, logger)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		50*time.Millisecond,
	)
	defer cancel()

	err = launcher.waitReady(ctx)
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf(
			"waitReady() error = %v, want %v",
			err,
			context.DeadlineExceeded,
		)
	}
}

func unusedAddress(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}

	address := listener.Addr().String()

	if err := listener.Close(); err != nil {
		t.Fatalf("listener.Close() error = %v", err)
	}

	return address
}
