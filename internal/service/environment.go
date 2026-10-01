package service

import (
	"context"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/fs"
)

func (s *Service) ValidateEnvironment(ctx context.Context, specPath string) (*model.Environment, error) {
	env, err := fs.Load[model.Environment](specPath)
	if err != nil {
		return nil, err
	}

	err = s.core.ValidateEnvironment(ctx, env)
	if err != nil {
		return nil, userFacingValidationError(err, "environment")
	}

	return env, nil
}

type ListEnvironmentsParams struct {
	Limit          int
	OrderBy        model.OrderBy
	OrderDirection model.OrderDirection
}

func (s *Service) ListEnvironments(ctx context.Context, params *ListEnvironmentsParams) ([]model.Environment, error) {
	arg := &coreapi.GetEnvironmentsPageParams{
		PageSize:       s.cfg.Output.ListMaxSize,
		OrderBy:        params.OrderBy,
		OrderDirection: params.OrderDirection,
	}

	if params.Limit > 0 && params.Limit < s.cfg.Output.ListMaxSize {
		arg.PageSize = params.Limit
	}

	page, err := s.core.GetEnvironmentsPage(ctx, arg)
	if err != nil {
		return nil, err
	}

	return page.Items, nil
}

func (s *Service) CreateEnvironment(ctx context.Context, specPath string) (*model.Environment, error) {
	env, err := fs.Load[model.Environment](specPath)
	if err != nil {
		return nil, err
	}

	arg := &coreapi.CreateEnvironmentParams{
		Name:        env.Name,
		Description: env.Description,
		Spec:        env.Spec,
	}

	created, err := s.core.CreateEnvironment(ctx, arg)
	if err != nil {
		return nil, err
	}

	return created, nil
}

func (s *Service) GetEnvironment(ctx context.Context, ref string) (*model.Environment, error) {
	environment, err := s.core.GetEnvironment(ctx, ref)
	if err != nil {
		return nil, err
	}

	return environment, nil
}

func (s *Service) DeleteEnvironment(ctx context.Context, ref string) (*model.Environment, error) {
	environment, err := s.core.DeleteEnvironment(ctx, ref)
	if err != nil {
		return nil, err
	}

	return environment, nil
}
