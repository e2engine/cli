package service

import (
	"context"
	nativeerrors "errors"
	"os"
	"strings"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/execute"
	corecall "github.com/e2engine/core/execute/call"
	corepersist "github.com/e2engine/core/execute/persist"
	"github.com/e2engine/core/execute/runtime"
	corerouter "github.com/e2engine/core/execute/runtime/router"
	coreservice "github.com/e2engine/core/execute/runtime/service"
	coreworker "github.com/e2engine/core/execute/worker"
	"github.com/e2engine/core/model"
	"github.com/e2engine/core/pkg/keys"
	corelog "github.com/e2engine/core/pkg/log"
	corescheduler "github.com/e2engine/core/scheduler"
	"github.com/e2engine/core/transport"
	"github.com/e2engine/core/transport/socket"
	repository "github.com/e2engine/repository/sqlite"
	runnerapi "github.com/e2engine/runner-local/api"
	runnerconfig "github.com/e2engine/runner-local/pkg/config"
	"github.com/ygrebnov/errorc"
	"github.com/ygrebnov/log"

	"github.com/e2engine/cli/internal/config"
	"github.com/e2engine/cli/internal/runnerprocess"
	"github.com/e2engine/cli/internal/util/paths"
	"github.com/e2engine/cli/pkg/errors"
)

const appName = "e2engine"

type settings struct {
	buildRunner bool
}

type option func(*settings)

func withRunner() option {
	return func(s *settings) {
		s.buildRunner = true
	}
}

func build(ctx context.Context, opts ...option) (s *Service, buildErr error) {
	o := &settings{}
	for _, opt := range opts {
		opt(o)
	}

	serviceCloser := newCloser()

	defer func() {
		if buildErr != nil {
			buildErr = nativeerrors.Join(buildErr, serviceCloser.Close())
		}
	}()

	cfg, err := LoadConfig(ctx)
	if err != nil {
		buildErr = err
		return
	}

	logger, err := BuildLogger(cfg.Logger)
	if err != nil {
		buildErr = err
		return
	}
	serviceCloser.Add(logger)

	repositories, dbCloser, err := BuildRepositories(ctx, cfg)
	if err != nil {
		buildErr = err
		return
	}
	serviceCloser.Add(dbCloser)

	var (
		scheduler     schedulerService
		runner        RunnerService
		producer      *transport.Producer
		rLauncher     runnerLauncher
		schedulerJobs chan execute.TestJob
		runnerJobs    chan execute.TestJob
	)

	if o.buildRunner {
		schedulerJobs = make(chan execute.TestJob, cfg.Scheduler.QueueSize)
		runnerJobs = schedulerJobs

		if cfg.Transport != nil && cfg.Transport.Kind == transport.KindSocket {
			sender := socket.NewSender(
				cfg.Transport.Address,
			)

			producer = transport.NewProducer(
				schedulerJobs,
				sender,
			)

			rLauncher = runnerprocess.NewLauncher(cfg.Transport.Address, logger)
		}

		scheduler, err = buildSchedulerService(schedulerJobs)
		if err != nil {
			buildErr = err
			return
		}

		if cfg.Transport == nil || cfg.Transport.Kind == transport.KindDirect {
			runner, err = BuildRunnerService(
				cfg,
				logger,
				repositories,
				runnerJobs,
			)
			if err != nil {
				buildErr = err
				return
			}
		}
	}

	cService, err := buildCoreService(
		cfg,
		logger,
		repositories,
		scheduler,
	)
	if err != nil {
		buildErr = err
		return
	}

	return newService(
		cfg,
		logger,
		serviceCloser,
		cService,
		runner,
		rLauncher,
		producer,
		schedulerJobs,
	), nil
}

func LoadConfig(ctx context.Context) (*config.Config, error) {
	// TODO: capture log message and add to log after it is initialized.
	configLogger, _ := corelog.NewSilentLogger()
	defer configLogger.Close()

	cfg, err := config.Load(ctx, configLogger)
	if err != nil {
		return nil, err
	}
	return cfg, nil
}

func BuildLogger(cfg *log.Config) (corelog.Logger, error) {
	if cfg == nil {
		var err error
		cfg, err = config.GetDefaultLoggerConfig()
		if err != nil {
			return nil, err
		}
	}

	logger, err := corelog.NewLogger(cfg)
	if err != nil {
		if logger != nil {
			err = nativeerrors.Join(err, logger.Close())
		}

		return nil, errorc.With(
			errors.ErrCannotInitializeLogger,
			errorc.Error(keys.Cause, err),
		)
	}

	return logger, nil
}

func BuildRunnerService(
	cfg *config.Config,
	logger corelog.Logger,
	repositories coreapi.Repositories,
	tests <-chan execute.TestJob,
) (RunnerService, error) {
	stateUpdater := buildStateUpdaterService(
		repositories.TestExecution,
		repositories.TestSuiteExecution,
	)

	routerProvider := corerouter.NewProvider()
	callStoreFactory := corecall.NewStoreFactory()
	serviceProvider := coreservice.NewProvider(cfg.Runtime.ServiceProvider, logger)

	testWorker, err := coreworker.NewWorker(
		cfg.Runner.Worker,
		logger,
		serviceProvider.GetGRPCMethodResolver(),
	)
	if err != nil {
		return nil, err
	}

	environment := runtime.NewEnvironmentController(
		cfg.Runtime,
		logger,
		serviceProvider,
		routerProvider,
		callStoreFactory,
	)

	return buildRunnerService(
		cfg.Runner,
		logger,
		testWorker,
		tests,
		stateUpdater,
		environment,
	)
}

func BuildRepositories(ctx context.Context, cfg *config.Config) (coreapi.Repositories, closer, error) {
	dataPath, err := getDataDir()
	if err != nil {
		return coreapi.Repositories{}, nil, err
	}

	dbPath := repository.GetPath(dataPath, cfg.DB)

	dbConn, err := repository.Open(ctx, dbPath, cfg.DB)
	if err != nil {
		return coreapi.Repositories{}, nil, err
	}

	repositories := repository.NewRepositories(dbConn)

	return repositories, dbConn, nil
}

func getDataDir() (string, error) {
	if path := strings.TrimSpace(os.Getenv("E2ENGINE_DATA_DIR")); path != "" {
		return path, nil
	}

	return paths.GetDataDir(appName)
}

func buildSchedulerService(tests chan<- execute.TestJob) (schedulerService, error) {
	return corescheduler.NewScheduler(tests)
}

func buildStateUpdaterService(
	testExecution corepersist.TestExecutionRepository,
	testSuiteExecution corepersist.TestSuiteExecutionRepository,
) stateUpdaterService {
	return corepersist.NewUpdater(testExecution, testSuiteExecution)
}

type environmentController interface {
	Acquire(
		ctx context.Context,
		instanceID string,
		env *model.Environment,
	) (*runtime.EnvironmentInstance, bool, error)

	Release(
		ctx context.Context,
		instanceID string,
	) error
	ReleaseAll(ctx context.Context) error
}

var _ environmentController = (*runtime.EnvironmentController)(nil)

func buildRunnerService(
	cfg *runnerconfig.Config,
	logger corelog.Logger,
	testWorker execute.Worker[execute.TestJob, execute.TestExecutionResult],
	tests <-chan execute.TestJob,
	stateUpdater stateUpdaterService,
	environment environmentController,
) (RunnerService, error) {
	return runnerapi.NewService(
		cfg,
		runnerapi.WithLogger(logger),
		runnerapi.WithTestWorker(testWorker),
		runnerapi.WithTestJobsChannel(tests),
		runnerapi.WithStateUpdater(stateUpdater),
		runnerapi.WithEnvironmentController(environment),
	)
}

func buildCoreService(
	cfg *config.Config,
	logger corelog.Logger,
	repositories coreapi.Repositories,
	scheduler schedulerService,
) (*coreapi.Service, error) {
	opts := []coreapi.Option{
		coreapi.WithLogger(logger),
		coreapi.WithRuntimeConfig(cfg.Runtime),
	}

	if scheduler != nil {
		opts = append(opts, coreapi.WithScheduler(scheduler))
	}

	return coreapi.NewService(
		repositories,
		opts...,
	)
}
