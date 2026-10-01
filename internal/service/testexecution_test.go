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

func TestService_ListTestExecutions(t *testing.T) {
	tests := []struct {
		name     string
		params   *ListTestExecutionsParams
		wantSize int
	}{
		{
			name: "default limit",
			params: &ListTestExecutionsParams{
				OrderBy:        model.OrderByID,
				OrderDirection: model.OrderDirectionAsc,
			},
			wantSize: 50,
		},
		{
			name: "smaller limit",
			params: &ListTestExecutionsParams{
				Limit:          10,
				OrderBy:        model.OrderByCreatedAt,
				OrderDirection: model.OrderDirectionDesc,
			},
			wantSize: 10,
		},
		{
			name: "zero limit",
			params: &ListTestExecutionsParams{
				Limit: 0,
			},
			wantSize: 50,
		},
		{
			name: "negative limit",
			params: &ListTestExecutionsParams{
				Limit: -1,
			},
			wantSize: 50,
		},
		{
			name: "limit equal to maximum",
			params: &ListTestExecutionsParams{
				Limit: 50,
			},
			wantSize: 50,
		},
		{
			name: "limit above maximum",
			params: &ListTestExecutionsParams{
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

			want := []model.TestExecution{
				{
					ID: "test-execution-1",
				},
			}

			core.EXPECT().
				GetTestExecutionsPage(
					gomock.Any(),
					&coreapi.GetTestExecutionsPageParams{
						PageSize:       test.wantSize,
						OrderBy:        test.params.OrderBy,
						OrderDirection: test.params.OrderDirection,
					},
				).
				Return(
					&model.TestExecutionsPage{
						Items: want,
					},
					nil,
				)

			got, err := service.ListTestExecutions(
				context.Background(),
				test.params,
			)
			if err != nil {
				t.Fatalf("ListTestExecutions() error = %v", err)
			}

			if len(got) != 1 {
				t.Fatalf(
					"ListTestExecutions() returned %d executions, want 1",
					len(got),
				)
			}

			if got[0].ID != want[0].ID {
				t.Errorf(
					"ListTestExecutions()[0].ID = %q, want %q",
					got[0].ID,
					want[0].ID,
				)
			}
		})
	}
}

func TestService_ListTestExecutions_Error(t *testing.T) {
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
		GetTestExecutionsPage(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.ListTestExecutions(
		context.Background(),
		&ListTestExecutionsParams{},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"ListTestExecutions() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetTestExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.TestExecution{
		ID: "test-execution-1",
	}

	core.EXPECT().
		GetTestExecution(
			gomock.Any(),
			"test-execution-ref",
		).
		Return(want, nil)

	got, err := service.GetTestExecution(
		context.Background(),
		"test-execution-ref",
	)
	if err != nil {
		t.Fatalf("GetTestExecution() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"GetTestExecution() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_GetTestExecution_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get failed")

	core.EXPECT().
		GetTestExecution(
			gomock.Any(),
			"test-execution-ref",
		).
		Return(nil, wantErr)

	_, err := service.GetTestExecution(
		context.Background(),
		"test-execution-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetTestExecution() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetTestExecutionStatus(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := model.ExecutionStatusPassed

	core.EXPECT().
		GetTestExecutionStatus(
			gomock.Any(),
			"test-execution-ref",
		).
		Return(want, nil)

	got, err := service.GetTestExecutionStatus(
		context.Background(),
		"test-execution-ref",
	)
	if err != nil {
		t.Fatalf("GetTestExecutionStatus() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"GetTestExecutionStatus() = %v, want %v",
			got,
			want,
		)
	}
}

func TestService_GetTestExecutionStatus_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get status failed")

	core.EXPECT().
		GetTestExecutionStatus(
			gomock.Any(),
			"test-execution-ref",
		).
		Return(model.ExecutionStatus(""), wantErr)

	_, err := service.GetTestExecutionStatus(
		context.Background(),
		"test-execution-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetTestExecutionStatus() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_DeleteTestExecution(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.TestExecution{
		ID: "test-execution-1",
	}

	core.EXPECT().
		DeleteTestExecution(
			gomock.Any(),
			"test-execution-ref",
		).
		Return(want, nil)

	got, err := service.DeleteTestExecution(
		context.Background(),
		"test-execution-ref",
	)
	if err != nil {
		t.Fatalf("DeleteTestExecution() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"DeleteTestExecution() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_DeleteTestExecution_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("delete failed")

	core.EXPECT().
		DeleteTestExecution(
			gomock.Any(),
			"test-execution-ref",
		).
		Return(nil, wantErr)

	_, err := service.DeleteTestExecution(
		context.Background(),
		"test-execution-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"DeleteTestExecution() error = %v, want %v",
			err,
			wantErr,
		)
	}
}
