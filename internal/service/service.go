package service

import (
	"context"
	"errors"
	"strings"
	"time"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/execute"
	"github.com/e2engine/core/model"
	coreerrors "github.com/e2engine/core/pkg/errors"
	"github.com/e2engine/core/pkg/log"
	"github.com/e2engine/core/transport"
	"github.com/ygrebnov/model/validation"

	"github.com/e2engine/cli/internal/config"
	internalerrors "github.com/e2engine/cli/internal/errors"
)

type coreService interface {
	GetEnvironmentsPage(
		ctx context.Context,
		params *coreapi.GetEnvironmentsPageParams,
	) (*model.EnvironmentsPage, error)
	ValidateEnvironment(ctx context.Context, env *model.Environment) error
	CreateEnvironment(ctx context.Context, params *coreapi.CreateEnvironmentParams) (*model.Environment, error)
	GetEnvironment(ctx context.Context, ref string) (*model.Environment, error)
	DeleteEnvironment(ctx context.Context, prefix string) (*model.Environment, error)

	GetTestsPage(ctx context.Context, params *coreapi.GetTestsPageParams) (*model.TestsPage, error)
	ValidateTest(ctx context.Context, test *model.Test) error
	CreateTest(ctx context.Context, params *coreapi.CreateTestParams) (*model.Test, error)
	GetTest(ctx context.Context, ref string) (*model.Test, error)
	DeleteTest(ctx context.Context, prefix string) (*model.Test, error)
	RunTest(ctx context.Context, params coreapi.TestRunRequest) (string, error)

	GetTestSuitesPage(ctx context.Context, params *coreapi.GetTestSuitesPageParams) (*model.TestSuitesPage, error)
	ValidateTestSuite(ctx context.Context, ts *model.TestSuite) error
	CreateTestSuite(
		ctx context.Context,
		params *coreapi.CreateTestSuiteParams,
	) (*model.TestSuite, error)
	GetTestSuite(ctx context.Context, ref string) (*model.TestSuite, error)
	DeleteTestSuite(ctx context.Context, prefix string) (*model.TestSuite, error)
	RunTestSuite(ctx context.Context, params coreapi.TestSuiteRunRequest) (string, error)

	GetTestExecutionsPage(
		ctx context.Context,
		params *coreapi.GetTestExecutionsPageParams,
	) (*model.TestExecutionsPage, error)
	GetTestExecution(ctx context.Context, ref string) (*model.TestExecution, error)
	GetTestExecutionStatus(ctx context.Context, ref string) (model.ExecutionStatus, error)
	DeleteTestExecution(ctx context.Context, prefix string) (*model.TestExecution, error)

	GetTestSuiteExecutionsPage(
		ctx context.Context,
		params *coreapi.GetTestSuiteExecutionsPageParams,
	) (*model.TestSuiteExecutionsPage, error)
	GetTestSuiteExecution(ctx context.Context, ref string) (*model.TestSuiteExecution, error)
	GetTestSuiteExecutionStatus(
		ctx context.Context,
		ref string,
	) (model.ExecutionStatus, error)
	DeleteTestSuiteExecution(ctx context.Context, prefix string) (*model.TestSuiteExecution, error)
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=service.go -destination=mock_service_test.go -package=service RunnerService,runnerLauncher

type RunnerService interface {
	StartTestExecution(ctx context.Context) error
	Run(ctx context.Context) error
}

type schedulerService interface {
	PublishTestJob(ctx context.Context, job execute.TestJob) error
}

type stateUpdaterService interface {
	SetTestRunning(ctx context.Context, executionID string, startedAt time.Time) error
	SetTestSuiteRunning(ctx context.Context, executionID string, startedAt time.Time) error
	SetTestCompleted(ctx context.Context, result execute.TestExecutionResult) (bool, error)
	SetTestSuiteCompleted(ctx context.Context, result execute.TestSuiteExecutionResult) (bool, error)
}

type runnerLauncher interface {
	Start(ctx context.Context) error
}

type Service struct {
	cfg    *config.Config
	logger log.Logger
	closer *closerStack

	core           coreService
	runner         RunnerService
	runnerLauncher runnerLauncher

	producer *transport.Producer

	schedulerJobs chan execute.TestJob
}

func newService(
	cfg *config.Config,
	logger log.Logger,
	closer *closerStack,
	core coreService,
	runner RunnerService,
	runnerLauncher runnerLauncher,
	producer *transport.Producer,
	schedulerJobs chan execute.TestJob,
) *Service {
	return &Service{
		cfg:            cfg,
		logger:         logger,
		closer:         closer,
		core:           core,
		runner:         runner,
		runnerLauncher: runnerLauncher,
		producer:       producer,
		schedulerJobs:  schedulerJobs,
	}
}

func (s *Service) GetConfig() *config.Config {
	return s.cfg
}

func (s *Service) GetLogger() log.Logger {
	return s.logger
}

func (s *Service) Close() error {
	if s.closer == nil {
		return nil
	}

	return s.closer.Close()
}

type closerStack struct {
	stack []closer
}

type closer interface {
	Close() error
}

func newCloser(stack ...closer) *closerStack {
	return &closerStack{
		stack: stack,
	}
}

func (c *closerStack) Add(resource closer) {
	c.stack = append(c.stack, resource)
}

func (c *closerStack) Close() error {
	var closeErr error
	for i := len(c.stack) - 1; i >= 0; i-- {
		closeErr = errors.Join(
			closeErr,
			c.stack[i].Close(),
		)
	}
	return closeErr
}

func userFacingValidationError(
	err error,
	kind string,
) error {
	if !errors.Is(err, coreerrors.ErrInvalidSpec) {
		return err
	}

	vErr, ok := errors.AsType[*validation.Error](err)
	if !ok {
		return err
	}

	return internalerrors.NewUserError(
		"invalid "+kind+" spec, fields: "+strings.Join(vErr.Fields(), ", "),
		err,
	)
}
