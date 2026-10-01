package internalrunner

import (
	"context"
	"errors"

	"github.com/e2engine/core/execute"
	"github.com/e2engine/core/transport"
	"github.com/e2engine/core/transport/socket"

	"github.com/e2engine/cli/internal/service"
)

// Run starts the internal runner process and serves test jobs received over
// the configured socket transport until the job stream ends or the context is canceled.
func Run(ctx context.Context) error {
	cfg, err := service.LoadConfig(ctx)
	if err != nil {
		return err
	}

	logger, err := service.BuildLogger(cfg.Logger)
	if err != nil {
		return err
	}
	defer logger.Close()

	repositories, closer, err := service.BuildRepositories(ctx, cfg)
	if err != nil {
		return err
	}
	defer closer.Close()

	runnerJobs := make(chan execute.TestJob, cfg.Scheduler.QueueSize)

	runner, err := service.BuildRunnerService(
		cfg,
		logger,
		repositories,
		runnerJobs,
	)
	if err != nil {
		return err
	}

	receiver := socket.NewReceiver(
		cfg.Transport.Address,
	)

	if err := receiver.Listen(); err != nil {
		return err
	}
	defer receiver.Close()

	consumer := transport.NewConsumer(
		runnerJobs,
		receiver,
	)

	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-runCtx.Done()
		_ = receiver.Close()
	}()

	// Important: initialize the worker stream before Runner.Run().
	if err := runner.StartTestExecution(runCtx); err != nil {
		return err
	}

	errorsCh := make(chan error, 2)

	go func() {
		err := consumer.Run(runCtx)
		if err != nil {
			cancel()
		}

		errorsCh <- err
	}()

	go func() {
		err := runner.Run(runCtx)
		if err != nil {
			cancel()
		}

		errorsCh <- err
	}()

	var runErr error

	for range 2 {
		err := <-errorsCh
		if err != nil &&
			!errors.Is(err, context.Canceled) &&
			runErr == nil {
			runErr = err
		}
	}

	return runErr
}
