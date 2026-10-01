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

func (s *Service) ValidateTestSuite(ctx context.Context, specPath string) (*model.TestSuite, error) {
	testSuite, err := fs.Load[model.TestSuite](specPath)
	if err != nil {
		return nil, err
	}

	err = s.core.ValidateTestSuite(ctx, testSuite)
	if err != nil {
		return nil, userFacingValidationError(err, "testsuite")
	}

	return testSuite, nil
}

type ListTestSuitesParams struct {
	Limit          int
	OrderBy        model.OrderBy
	OrderDirection model.OrderDirection
}

func (s *Service) ListTestSuites(ctx context.Context, params *ListTestSuitesParams) ([]model.TestSuite, error) {
	arg := &coreapi.GetTestSuitesPageParams{
		PageSize:       s.cfg.Output.ListMaxSize,
		OrderBy:        params.OrderBy,
		OrderDirection: params.OrderDirection,
	}

	if params.Limit > 0 && params.Limit < s.cfg.Output.ListMaxSize {
		arg.PageSize = params.Limit
	}

	page, err := s.core.GetTestSuitesPage(ctx, arg)
	if err != nil {
		return nil, err
	}

	return page.Items, nil
}

func (s *Service) CreateTestSuite(ctx context.Context, specPath string) (*model.TestSuite, error) {
	test, err := fs.Load[model.TestSuite](specPath)
	if err != nil {
		return nil, err
	}

	arg := &coreapi.CreateTestSuiteParams{
		Name:        test.Name,
		Description: test.Description,
		Spec:        test.Spec,
	}

	created, err := s.core.CreateTestSuite(ctx, arg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetTestSuite(ctx context.Context, ref string) (*model.TestSuite, error) {
	test, err := s.core.GetTestSuite(ctx, ref)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (s *Service) DeleteTestSuite(ctx context.Context, ref string) (*model.TestSuite, error) {
	test, err := s.core.DeleteTestSuite(ctx, ref)
	if err != nil {
		return nil, err
	}

	return test, nil
}

func (s *Service) RunTestSuite(ctx context.Context, environmentRef, testSuiteRef string) (string, error) {
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

	arg := coreapi.TestSuiteRunRequest{
		EnvironmentRef: environmentRef,
		TestSuiteRef:   testSuiteRef,
	}

	testSuiteExecutionID, err := s.core.RunTestSuite(runCtx, arg)
	if err != nil {
		return "", err
	}

	close(s.schedulerJobs)

	if s.runnerLauncher != nil {
		if err := <-producerErrors; err != nil {
			return testSuiteExecutionID, err
		}

		return testSuiteExecutionID, nil
	}

	if err := s.runner.Run(runCtx); err != nil {
		return testSuiteExecutionID, err
	}

	return testSuiteExecutionID, nil
}
