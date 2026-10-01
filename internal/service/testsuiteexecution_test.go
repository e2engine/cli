package service

import (
	"context"
	"errors"
	"testing"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/model"
	"go.uber.org/mock/gomock"

	"github.com/e2engine/cli/internal/config"
)

func TestService_ListTestSuiteExecutions(t *testing.T) {
	tests := []struct {
		name     string
		params   *ListTestSuiteExecutionsParams
		wantSize int
	}{
		{
			name: "default limit",
			params: &ListTestSuiteExecutionsParams{
				OrderBy:        model.OrderByID,
				OrderDirection: model.OrderDirectionAsc,
			},
			wantSize: 50,
		},
		{
			name: "smaller limit",
			params: &ListTestSuiteExecutionsParams{
				Limit:          10,
				OrderBy:        model.OrderByCreatedAt,
				OrderDirection: model.OrderDirectionDesc,
			},
			wantSize: 10,
		},
		{
			name: "zero limit",
			params: &ListTestSuiteExecutionsParams{
				Limit: 0,
			},
			wantSize: 50,
		},
		{
			name: "negative limit",
			params: &ListTestSuiteExecutionsParams{
				Limit: -1,
			},
			wantSize: 50,
		},
		{
			name: "limit equal to maximum",
			params: &ListTestSuiteExecutionsParams{
				Limit: 50,
			},
			wantSize: 50,
		},
		{
			name: "limit above maximum",
			params: &ListTestSuiteExecutionsParams{
				Limit: 100,
			},
			wantSize: 50,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			core := NewMockcoreService(ctrl)

			service := &Service{
				cfg: &config.Config{
					Output: &config.OutputConfig{
						ListMaxSize: 50,
					},
				},
				core: core,
			}

			want := []model.TestSuiteExecution{
				{
					ID: "testsuite-execution-1",
				},
			}

			core.EXPECT().
				GetTestSuiteExecutionsPage(
					gomock.Any(),
					&coreapi.GetTestSuiteExecutionsPageParams{
						PageSize:       test.wantSize,
						OrderBy:        test.params.OrderBy,
						OrderDirection: test.params.OrderDirection,
					},
				).
				Return(
					&model.TestSuiteExecutionsPage{
						Items: want,
					},
					nil,
				)

			got, err := service.ListTestSuiteExecutions(
				context.Background(),
				test.params,
			)
			if err != nil {
				t.Fatalf(
					"ListTestSuiteExecutions() error = %v",
					err,
				)
			}

			if len(got) != 1 {
				t.Fatalf(
					"ListTestSuiteExecutions() returned %d executions, want 1",
					len(got),
				)
			}

			if got[0].ID != want[0].ID {
				t.Errorf(
					"ListTestSuiteExecutions()[0].ID = %q, want %q",
					got[0].ID,
					want[0].ID,
				)
			}
		})
	}
}

func TestService_ListTestSuiteExecutions_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		cfg: &config.Config{
			Output: &config.OutputConfig{
				ListMaxSize: 50,
			},
		},
		core: core,
	}

	wantErr := errors.New("list failed")

	core.EXPECT().
		GetTestSuiteExecutionsPage(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.ListTestSuiteExecutions(
		context.Background(),
		&ListTestSuiteExecutionsParams{},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"ListTestSuiteExecutions() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetTestSuiteExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.TestSuiteExecution{
		ID: "testsuite-execution-1",
	}

	core.EXPECT().
		GetTestSuiteExecution(
			gomock.Any(),
			"testsuite-execution-ref",
		).
		Return(want, nil)

	got, err := service.GetTestSuiteExecution(
		context.Background(),
		"testsuite-execution-ref",
	)
	if err != nil {
		t.Fatalf(
			"GetTestSuiteExecution() error = %v",
			err,
		)
	}

	if got != want {
		t.Fatalf(
			"GetTestSuiteExecution() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_GetTestSuiteExecution_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get failed")

	core.EXPECT().
		GetTestSuiteExecution(
			gomock.Any(),
			"testsuite-execution-ref",
		).
		Return(nil, wantErr)

	_, err := service.GetTestSuiteExecution(
		context.Background(),
		"testsuite-execution-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetTestSuiteExecution() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetTestSuiteExecutionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := model.ExecutionStatusRunning

	core.EXPECT().
		GetTestSuiteExecutionStatus(
			gomock.Any(),
			"testsuite-execution-ref",
		).
		Return(want, nil)

	got, err := service.GetTestSuiteExecutionStatus(
		context.Background(),
		"testsuite-execution-ref",
	)
	if err != nil {
		t.Fatalf(
			"GetTestSuiteExecutionStatus() error = %v",
			err,
		)
	}

	if got != want {
		t.Fatalf(
			"GetTestSuiteExecutionStatus() = %v, want %v",
			got,
			want,
		)
	}
}

func TestService_GetTestSuiteExecutionStatus_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get failed")

	core.EXPECT().
		GetTestSuiteExecutionStatus(
			gomock.Any(),
			"testsuite-execution-ref",
		).
		Return(model.ExecutionStatus(""), wantErr)

	_, err := service.GetTestSuiteExecutionStatus(
		context.Background(),
		"testsuite-execution-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetTestSuiteExecutionStatus() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_DeleteTestSuiteExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.TestSuiteExecution{
		ID: "testsuite-execution-1",
	}

	core.EXPECT().
		DeleteTestSuiteExecution(
			gomock.Any(),
			"testsuite-execution-ref",
		).
		Return(want, nil)

	got, err := service.DeleteTestSuiteExecution(
		context.Background(),
		"testsuite-execution-ref",
	)
	if err != nil {
		t.Fatalf(
			"DeleteTestSuiteExecution() error = %v",
			err,
		)
	}

	if got != want {
		t.Fatalf(
			"DeleteTestSuiteExecution() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_DeleteTestSuiteExecution_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("delete failed")

	core.EXPECT().
		DeleteTestSuiteExecution(
			gomock.Any(),
			"testsuite-execution-ref",
		).
		Return(nil, wantErr)

	_, err := service.DeleteTestSuiteExecution(
		context.Background(),
		"testsuite-execution-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"DeleteTestSuiteExecution() error = %v, want %v",
			err,
			wantErr,
		)
	}
}
