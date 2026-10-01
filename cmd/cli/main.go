package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/e2engine/cli/internal/service"
)

func main() {
	os.Exit(run())
}

func run() int {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	streams := basicStreams{
		in:     os.Stdin,
		out:    os.Stdout,
		errOut: os.Stderr,
	}

	p := providers{
		base:       service.NewProvider(),
		withRunner: service.NewProviderWithRunner(),
	}
	defer func() {
		if err := p.base.Close(); err != nil {
			_, _ = fmt.Fprintln(streams.ErrOut(), err)
		}
		if err := p.withRunner.Close(); err != nil {
			_, _ = fmt.Fprintln(streams.ErrOut(), err)
		}
	}()

	return executeRootCommand(ctx, streams, p)
}

type serviceProvider interface {
	Get(ctx context.Context) (*service.Service, error)
	Close() error
}
