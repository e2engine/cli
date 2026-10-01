package service

import (
	"context"
	nativeerrors "errors"
	"testing"

	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestProvider_Close_NilProvider(t *testing.T) {
	var provider *Provider

	if err := provider.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestProvider_Close_NotInitialized(t *testing.T) {
	provider := NewProvider()

	if err := provider.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
}

func TestProvider_Get_ReturnsCachedError(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	provider := NewProvider()

	_, err1 := provider.Get(ctx)
	if err1 == nil {
		t.Fatal("first Get() error = nil, want error")
	}

	_, err2 := provider.Get(context.Background())
	if err2 == nil {
		t.Fatal("second Get() error = nil, want error")
	}

	if err1 != err2 {
		t.Fatalf(
			"Get() errors differ: first = %v, second = %v",
			err1,
			err2,
		)
	}
}

func TestProvider_Close_WrapsServiceError(t *testing.T) {
	closeErr := nativeerrors.New("close failed")

	provider := &Provider{
		app: &Service{
			closer: &closerStack{
				stack: []closer{
					&providerTestCloser{err: closeErr},
				},
			},
		},
	}

	err := provider.Close()
	if err == nil {
		t.Fatal("Close() error = nil, want error")
	}

	if !nativeerrors.Is(err, clierrors.ErrFailedToCloseService) {
		t.Fatalf(
			"Close() error = %v, want errors.Is(..., ErrFailedToCloseService)",
			err,
		)
	}
}

type providerTestCloser struct {
	err error
}

func (c *providerTestCloser) Close() error {
	return c.err
}
