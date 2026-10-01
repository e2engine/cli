package testexecution

import (
	"context"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/service"
)

type internalService interface {
	ListTestExecutions(ctx context.Context, params *service.ListTestExecutionsParams) ([]model.TestExecution, error)
	GetTestExecution(ctx context.Context, ref string) (*model.TestExecution, error)
	GetTestExecutionStatus(ctx context.Context, ref string) (model.ExecutionStatus, error)
	DeleteTestExecution(ctx context.Context, ref string) (*model.TestExecution, error)
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=service.go -destination=mock_service_test.go -package=testexecution
