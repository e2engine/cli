package service

import (
	"context"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/model"
)

type ListTestExecutionsParams struct {
	Limit          int
	OrderBy        model.OrderBy
	OrderDirection model.OrderDirection
}

func (s *Service) ListTestExecutions(
	ctx context.Context,
	params *ListTestExecutionsParams,
) ([]model.TestExecution, error) {
	arg := &coreapi.GetTestExecutionsPageParams{
		PageSize:       s.cfg.Output.ListMaxSize,
		OrderBy:        params.OrderBy,
		OrderDirection: params.OrderDirection,
	}

	if params.Limit > 0 && params.Limit < s.cfg.Output.ListMaxSize {
		arg.PageSize = params.Limit
	}

	page, err := s.core.GetTestExecutionsPage(ctx, arg)
	if err != nil {
		return nil, err
	}

	return page.Items, nil
}

func (s *Service) GetTestExecution(ctx context.Context, ref string) (*model.TestExecution, error) {
	testExecution, err := s.core.GetTestExecution(ctx, ref)
	if err != nil {
		return nil, err
	}

	return testExecution, nil
}

func (s *Service) GetTestExecutionStatus(ctx context.Context, ref string) (model.ExecutionStatus, error) {
	status, err := s.core.GetTestExecutionStatus(ctx, ref)
	if err != nil {
		return "", err
	}

	return status, nil
}

func (s *Service) DeleteTestExecution(ctx context.Context, ref string) (*model.TestExecution, error) {
	testExecution, err := s.core.DeleteTestExecution(ctx, ref)
	if err != nil {
		return nil, err
	}

	return testExecution, nil
}
