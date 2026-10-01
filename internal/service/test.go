package service

import (
	"context"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/model"
	"github.com/e2engine/core/pkg/keys"
	"github.com/ygrebnov/errorc"

	"github.com/e2engine/cli/internal/fs"
	"github.com/e2engine/cli/pkg/errors"
)

func (s *Service) ValidateTest(ctx context.Context, specPath string) (*model.Test, error) {
	test, err := fs.Load[model.Test](specPath)
	if err != nil {
		return nil, err
	}

	err = s.core.ValidateTest(ctx, test)
	if err != nil {
		return nil, userFacingValidationError(err, "test")
	}

	return test, nil
}

type ListTestsParams struct {
	Limit          int
	OrderBy        model.OrderBy
	OrderDirection model.OrderDirection
}

func (s *Service) ListTests(ctx context.Context, params *ListTestsParams) ([]model.Test, error) {
	arg := &coreapi.GetTestsPageParams{
		PageSize:       s.cfg.Output.ListMaxSize,
		OrderBy:        params.OrderBy,
		OrderDirection: params.OrderDirection,
	}

	if params.Limit > 0 && params.Limit < s.cfg.Output.ListMaxSize {
		arg.PageSize = params.Limit
	}

	page, err := s.core.GetTestsPage(ctx, arg)
	if err != nil {
		return nil, err
	}

	return page.Items, nil
}

func (s *Service) CreateTest(ctx context.Context, specPath string) (*model.Test, error) {
	test, err := fs.Load[model.Test](specPath)
	if err != nil {
		return nil, err
	}

	arg := &coreapi.CreateTestParams{
		Name:        test.Name,
		Description: test.Description,
		Spec:        test.Spec,
	}

	created, err := s.core.CreateTest(ctx, arg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetTest(ctx context.Context, ref string) (*model.Test, error) {
	test, err := s.core.GetTest(ctx, ref)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (s *Service) DeleteTest(ctx context.Context, ref string) (*model.Test, error) {
	test, err := s.core.DeleteTest(ctx, ref)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (s *Service) RunTest(
	ctx context.Context,
	environmentRef string,
	testRef string,
) (string, error) {
	runCtx, cancel := context.WithCancel(ctx)
	defer cancel()

	var producerErrors chan error

	if s.runnerLauncher != nil {
		if err := s.runnerLauncher.Start(runCtx); err != nil {
			return "", errorc.With(
				errors.ErrFailedToStartRunner,
				errorc.Error(keys.Cause, err),
			)
		}

		s.logger.Debug("Started runner process")

		producerErrors = make(chan error, 1)

		go func() {
			err := s.producer.Run(runCtx)
			if err != nil {
				cancel()
			}

			producerErrors <- err
		}()
	} else {
		if err := s.runner.StartTestExecution(runCtx); err != nil {
			return "", errorc.With(
				errors.ErrFailedToStartRunner,
				errorc.Error(keys.Cause, err),
			)
		}

		s.logger.Debug("Started runner service")
	}

	arg := coreapi.TestRunRequest{
		EnvironmentRef: environmentRef,
		TestRef:        testRef,
	}

	testExecutionID, err := s.core.RunTest(runCtx, arg)
	if err != nil {
		return "", err
	}

	close(s.schedulerJobs)

	if s.runnerLauncher != nil {
		if err := <-producerErrors; err != nil {
			return testExecutionID, err
		}

		return testExecutionID, nil
	}

	if err := s.runner.Run(runCtx); err != nil {
		return testExecutionID, err
	}

	return testExecutionID, nil
}
