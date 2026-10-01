package testsuite

import (
	"context"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/service"
)

type internalService interface {
	ValidateTestSuite(ctx context.Context, specPath string) (*model.TestSuite, error)
	ListTestSuites(ctx context.Context, params *service.ListTestSuitesParams) ([]model.TestSuite, error)
	CreateTestSuite(ctx context.Context, specPath string) (*model.TestSuite, error)
	GetTestSuite(ctx context.Context, ref string) (*model.TestSuite, error)
	DeleteTestSuite(ctx context.Context, ref string) (*model.TestSuite, error)
	RunTestSuite(
		ctx context.Context,
		environmentRef string,
		testSuiteRef string,
	) (string, error)
}

//go:generate go run -mod=mod go.uber.org/mock/mockgen -source=service.go -destination=mock_service_test.go -package=testsuite
