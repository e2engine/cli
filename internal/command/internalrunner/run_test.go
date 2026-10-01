package internalrunner

import (
	"context"
	"net"
	"testing"
	"time"

	"github.com/e2engine/core/transport"
	"github.com/e2engine/core/transport/socket"
)

func TestRun_ContextCanceled(t *testing.T) {
	t.Setenv("E2ENGINE_DATA_DIR", t.TempDir())
	t.Setenv("E2ENGINE_TRANSPORT_KIND", string(transport.KindSocket))

	address := freeAddress(t)
	t.Setenv("E2ENGINE_TRANSPORT_ADDRESS", address)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- Run(ctx)
	}()

	waitForListener(t, ctx, address)

	cancel()

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}

	case <-time.After(time.Second):
		t.Fatal("Run() did not stop after context cancellation")
	}
}

func TestRun_AddressAlreadyInUse(t *testing.T) {
	t.Setenv("E2ENGINE_DATA_DIR", t.TempDir())
	t.Setenv("E2ENGINE_TRANSPORT_KIND", string(transport.KindSocket))

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		t.Fatalf("net.Listen() error = %v", err)
	}
	defer listener.Close()

	t.Setenv(
		"E2ENGINE_TRANSPORT_ADDRESS",
		listener.Addr().String(),
	)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	if err := Run(ctx); err == nil {
		t.Fatal("Run() error = nil, want non-nil")
	}
}

func TestRun_EndMessage(t *testing.T) {
	t.Setenv("E2ENGINE_DATA_DIR", t.TempDir())
	t.Setenv("E2ENGINE_TRANSPORT_KIND", string(transport.KindSocket))

	address := freeAddress(t)
	t.Setenv("E2ENGINE_TRANSPORT_ADDRESS", address)

	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Second,
	)
	defer cancel()

	errCh := make(chan error, 1)

	go func() {
		errCh <- Run(ctx)
	}()

	waitForListener(t, ctx, address)

	sender := socket.NewSender(address)

	if err := sender.Send(
		ctx,
		transport.NewEndMessage(),
	); err != nil {
		t.Fatalf("Send() error = %v", err)
	}

	select {
	case err := <-errCh:
		if err != nil {
			t.Fatalf("Run() error = %v", err)
		}

	case <-ctx.Done():
		t.Fatalf(
			"Run() did not stop after end message: %v",
			ctx.Err(),
		)
	}
}

func freeAddress(t *testing.T) string {
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

func waitForListener(
	t *testing.T,
	ctx context.Context,
	address string,
) {
	t.Helper()

	ticker := time.NewTicker(10 * time.Millisecond)
	defer ticker.Stop()

	for {
		conn, err := (&net.Dialer{}).DialContext(
			ctx,
			"tcp",
			address,
		)
		if err == nil {
			if err := conn.Close(); err != nil {
				t.Fatalf("conn.Close() error = %v", err)
			}

			return
		}

		select {
		case <-ctx.Done():
			t.Fatalf(
				"listener %q did not become ready: %v",
				address,
				ctx.Err(),
			)

		case <-ticker.C:
		}
	}
}
