package service

import (
	"context"
	"errors"
	"testing"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/execute"
	"github.com/e2engine/core/pkg/log"
	"go.uber.org/mock/gomock"

	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestService_RunTestSuite_Direct(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	runner := NewMockRunnerService(ctrl)

	jobs := make(chan execute.TestJob)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	service := &Service{
		logger:        logger,
		core:          core,
		runner:        runner,
		schedulerJobs: jobs,
	}

	runner.EXPECT().
		StartTestExecution(gomock.Any()).
		Return(nil)

	core.EXPECT().
		RunTestSuite(
			gomock.Any(),
			coreapi.TestSuiteRunRequest{
				EnvironmentRef: "environment-ref",
				TestSuiteRef:   "testsuite-ref",
			},
		).
		Return("testsuite-execution-id", nil)

	runner.EXPECT().
		Run(gomock.Any()).
		DoAndReturn(func(context.Context) error {
			select {
			case _, ok := <-jobs:
				if ok {
					t.Fatal("schedulerJobs is not closed")
				}
			default:
				t.Fatal("schedulerJobs is not closed")
			}

			return nil
		})

	got, err := service.RunTestSuite(
		context.Background(),
		"environment-ref",
		"testsuite-ref",
	)
	if err != nil {
		t.Fatalf("RunTestSuite() error = %v", err)
	}

	if got != "testsuite-execution-id" {
		t.Errorf(
			"RunTestSuite() = %q, want %q",
			got,
			"testsuite-execution-id",
		)
	}
}

func TestService_RunTestSuite_StartRunnerError(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	runner := NewMockRunnerService(ctrl)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	service := &Service{
		logger:        logger,
		core:          core,
		runner:        runner,
		schedulerJobs: make(chan execute.TestJob),
	}

	wantErr := errors.New("start runner failed")

	runner.EXPECT().
		StartTestExecution(gomock.Any()).
		Return(wantErr)

	got, err := service.RunTestSuite(
		context.Background(),
		"environment-ref",
		"testsuite-ref",
	)

	if got != "" {
		t.Errorf("RunTestSuite() = %q, want empty ID", got)
	}

	if !errors.Is(err, clierrors.ErrFailedToStartRunner) {
		t.Fatalf(
			"RunTestSuite() error = %v, want ErrFailedToStartRunner",
			err,
		)
	}
}

func TestService_RunTestSuite_CoreError(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	runner := NewMockRunnerService(ctrl)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	jobs := make(chan execute.TestJob)

	service := &Service{
		logger:        logger,
		core:          core,
		runner:        runner,
		schedulerJobs: jobs,
	}

	runner.EXPECT().
		StartTestExecution(gomock.Any()).
		Return(nil)

	wantErr := errors.New("run testsuite failed")

	core.EXPECT().
		RunTestSuite(
			gomock.Any(),
			coreapi.TestSuiteRunRequest{
				EnvironmentRef: "environment-ref",
				TestSuiteRef:   "testsuite-ref",
			},
		).
		Return("", wantErr)

	got, err := service.RunTestSuite(
		context.Background(),
		"environment-ref",
		"testsuite-ref",
	)

	if got != "" {
		t.Errorf("RunTestSuite() = %q, want empty ID", got)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"RunTestSuite() error = %v, want %v",
			err,
			wantErr,
		)
	}

	select {
	case _, ok := <-jobs:
		if !ok {
			t.Fatal("schedulerJobs was closed on core error")
		}
	default:
	}
}

func TestService_RunTestSuite_RunnerError(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	runner := NewMockRunnerService(ctrl)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	service := &Service{
		logger:        logger,
		core:          core,
		runner:        runner,
		schedulerJobs: make(chan execute.TestJob),
	}

	runner.EXPECT().
		StartTestExecution(gomock.Any()).
		Return(nil)

	core.EXPECT().
		RunTestSuite(
			gomock.Any(),
			coreapi.TestSuiteRunRequest{
				EnvironmentRef: "environment-ref",
				TestSuiteRef:   "testsuite-ref",
			},
		).
		Return("testsuite-execution-id", nil)

	wantErr := errors.New("runner failed")

	runner.EXPECT().
		Run(gomock.Any()).
		Return(wantErr)

	got, err := service.RunTestSuite(
		context.Background(),
		"environment-ref",
		"testsuite-ref",
	)

	if got != "testsuite-execution-id" {
		t.Errorf(
			"RunTestSuite() = %q, want %q",
			got,
			"testsuite-execution-id",
		)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"RunTestSuite() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_RunTestSuite_LauncherError(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	launcher := NewMockrunnerLauncher(ctrl)

	logger, err := log.NewSilentLogger()
	if err != nil {
		t.Fatalf("Failed to create logger: %v", err)
	}

	service := &Service{
		logger:         logger,
		core:           core,
		runnerLauncher: launcher,
		schedulerJobs:  make(chan execute.TestJob),
	}

	wantErr := errors.New("launcher failed")

	launcher.EXPECT().
		Start(gomock.Any()).
		Return(wantErr)

	got, err := service.RunTestSuite(
		context.Background(),
		"environment-ref",
		"testsuite-ref",
	)

	if got != "" {
		t.Errorf("RunTestSuite() = %q, want empty ID", got)
	}

	if !errors.Is(err, clierrors.ErrFailedToStartRunner) {
		t.Fatalf(
			"RunTestSuite() error = %v, want ErrFailedToStartRunner",
			err,
		)
	}
}
