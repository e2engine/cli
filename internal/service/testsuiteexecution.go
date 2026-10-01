package service

import (
	"context"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/model"
)

type ListTestSuiteExecutionsParams struct {
	Limit          int
	OrderBy        model.OrderBy
	OrderDirection model.OrderDirection
}

func (s *Service) ListTestSuiteExecutions(
	ctx context.Context,
	params *ListTestSuiteExecutionsParams,
) ([]model.TestSuiteExecution, error) {
	arg := &coreapi.GetTestSuiteExecutionsPageParams{
		PageSize:       s.cfg.Output.ListMaxSize,
		OrderBy:        params.OrderBy,
		OrderDirection: params.OrderDirection,
	}

	if params.Limit > 0 && params.Limit < s.cfg.Output.ListMaxSize {
		arg.PageSize = params.Limit
	}

	page, err := s.core.GetTestSuiteExecutionsPage(ctx, arg)
	if err != nil {
		return nil, err
	}

	return page.Items, nil
}

func (s *Service) GetTestSuiteExecution(ctx context.Context, ref string) (*model.TestSuiteExecution, error) {
	testSuiteExecution, err := s.core.GetTestSuiteExecution(ctx, ref)
	if err != nil {
		return nil, err
	}

	return testSuiteExecution, nil
}

func (s *Service) GetTestSuiteExecutionStatus(ctx context.Context, ref string) (model.ExecutionStatus, error) {
	status, err := s.core.GetTestSuiteExecutionStatus(ctx, ref)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (s *Service) DeleteTestSuiteExecution(ctx context.Context, ref string) (*model.TestSuiteExecution, error) {
	testSuiteExecution, err := s.core.DeleteTestSuiteExecution(ctx, ref)
	if err != nil {
		return nil, err
	}

	return testSuiteExecution, nil
}
