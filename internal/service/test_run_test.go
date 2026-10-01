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

func TestService_RunTest_Direct(t *testing.T) {
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
		RunTest(
			gomock.Any(),
			coreapi.TestRunRequest{
				EnvironmentRef: "environment-ref",
				TestRef:        "test-ref",
			},
		).
		Return("execution-id", nil)

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

	got, err := service.RunTest(
		context.Background(),
		"environment-ref",
		"test-ref",
	)
	if err != nil {
		t.Fatalf("RunTest() error = %v", err)
	}

	if got != "execution-id" {
		t.Errorf(
			"RunTest() = %q, want %q",
			got,
			"execution-id",
		)
	}
}

func TestService_RunTest_StartRunnerError(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	runner := NewMockRunnerService(ctrl)

	service := &Service{
		core:          core,
		runner:        runner,
		schedulerJobs: make(chan execute.TestJob),
	}

	wantErr := errors.New("start runner failed")

	runner.EXPECT().
		StartTestExecution(gomock.Any()).
		Return(wantErr)

	got, err := service.RunTest(
		context.Background(),
		"environment-ref",
		"test-ref",
	)

	if got != "" {
		t.Errorf("RunTest() = %q, want empty ID", got)
	}

	if !errors.Is(err, clierrors.ErrFailedToStartRunner) {
		t.Fatalf(
			"RunTest() error = %v, want ErrFailedToStartRunner",
			err,
		)
	}
}

func TestService_RunTest_CoreError(t *testing.T) {
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

	wantErr := errors.New("run test failed")

	core.EXPECT().
		RunTest(
			gomock.Any(),
			coreapi.TestRunRequest{
				EnvironmentRef: "environment-ref",
				TestRef:        "test-ref",
			},
		).
		Return("", wantErr)

	got, err := service.RunTest(
		context.Background(),
		"environment-ref",
		"test-ref",
	)

	if got != "" {
		t.Errorf("RunTest() = %q, want empty ID", got)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"RunTest() error = %v, want %v",
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

func TestService_RunTest_RunnerError(t *testing.T) {
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
		RunTest(
			gomock.Any(),
			coreapi.TestRunRequest{
				EnvironmentRef: "environment-ref",
				TestRef:        "test-ref",
			},
		).
		Return("execution-id", nil)

	wantErr := errors.New("runner failed")

	runner.EXPECT().
		Run(gomock.Any()).
		Return(wantErr)

	got, err := service.RunTest(
		context.Background(),
		"environment-ref",
		"test-ref",
	)

	if got != "execution-id" {
		t.Errorf(
			"RunTest() = %q, want %q",
			got,
			"execution-id",
		)
	}

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"RunTest() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_RunTest_LauncherError(t *testing.T) {
	ctrl := gomock.NewController(t)

	core := NewMockcoreService(ctrl)
	launcher := NewMockrunnerLauncher(ctrl)

	service := &Service{
		core:           core,
		runnerLauncher: launcher,
		schedulerJobs:  make(chan execute.TestJob),
	}

	wantErr := errors.New("launcher failed")

	launcher.EXPECT().
		Start(gomock.Any()).
		Return(wantErr)

	got, err := service.RunTest(
		context.Background(),
		"environment-ref",
		"test-ref",
	)

	if got != "" {
		t.Errorf("RunTest() = %q, want empty ID", got)
	}

	if !errors.Is(err, clierrors.ErrFailedToStartRunner) {
		t.Fatalf(
			"RunTest() error = %v, want ErrFailedToStartRunner",
			err,
		)
	}
}
