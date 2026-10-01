package testsuiteexecution

import (
	"context"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/service"
)

type internalService interface {
	ListTestSuiteExecutions(
		ctx context.Context,
		params *service.ListTestSuiteExecutionsParams,
	) ([]model.TestSuiteExecution, error)
	GetTestSuiteExecution(
		ctx context.Context,
		ref string,
	) (*model.TestSuiteExecution, error)
	GetTestSuiteExecutionStatus(
		ctx context.Context,
		ref string,
	) (model.ExecutionStatus, error)
	DeleteTestSuiteExecution(
		ctx context.Context,
		ref string,
	) (*model.TestSuiteExecution, error)
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=service.go -destination=mock_service_test.go -package=testsuiteexecution
