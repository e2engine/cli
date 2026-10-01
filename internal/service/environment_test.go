package service

import (
	"context"
	"errors"
	"os"
	"path/filepath"
	"testing"

	coreapi "github.com/e2engine/core/api"
	"github.com/e2engine/core/model"
	"go.uber.org/mock/gomock"

	"github.com/e2engine/cli/internal/config"
)

func TestService_ListEnvironments(t *testing.T) {
	tests := []struct {
		name     string
		params   *ListEnvironmentsParams
		wantSize int
	}{
		{
			name: "default limit",
			params: &ListEnvironmentsParams{
				OrderBy:        model.OrderByID,
				OrderDirection: model.OrderDirectionAsc,
			},
			wantSize: 50,
		},
		{
			name: "smaller limit",
			params: &ListEnvironmentsParams{
				Limit:          10,
				OrderBy:        model.OrderByName,
				OrderDirection: model.OrderDirectionDesc,
			},
			wantSize: 10,
		},
		{
			name: "zero limit",
			params: &ListEnvironmentsParams{
				Limit: 0,
			},
			wantSize: 50,
		},
		{
			name: "negative limit",
			params: &ListEnvironmentsParams{
				Limit: -1,
			},
			wantSize: 50,
		},
		{
			name: "limit equal to maximum",
			params: &ListEnvironmentsParams{
				Limit: 50,
			},
			wantSize: 50,
		},
		{
			name: "limit above maximum",
			params: &ListEnvironmentsParams{
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

			want := []model.Environment{
				{
					ID:   "environment-1",
					Name: "environment",
				},
			}

			core.EXPECT().
				GetEnvironmentsPage(
					gomock.Any(),
					&coreapi.GetEnvironmentsPageParams{
						PageSize:       test.wantSize,
						OrderBy:        test.params.OrderBy,
						OrderDirection: test.params.OrderDirection,
					},
				).
				Return(
					&model.EnvironmentsPage{
						Items: want,
					},
					nil,
				)

			got, err := service.ListEnvironments(
				context.Background(),
				test.params,
			)
			if err != nil {
				t.Fatalf("ListEnvironments() error = %v", err)
			}

			if len(got) != 1 {
				t.Fatalf(
					"ListEnvironments() returned %d environments, want 1",
					len(got),
				)
			}

			if got[0].ID != want[0].ID {
				t.Errorf(
					"ListEnvironments()[0].ID = %q, want %q",
					got[0].ID,
					want[0].ID,
				)
			}
		})
	}
}

func TestService_ListEnvironments_Error(t *testing.T) {
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
		GetEnvironmentsPage(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.ListEnvironments(
		context.Background(),
		&ListEnvironmentsParams{},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"ListEnvironments() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.Environment{
		ID:   "environment-1",
		Name: "environment",
	}

	core.EXPECT().
		GetEnvironment(
			gomock.Any(),
			"environment-ref",
		).
		Return(want, nil)

	got, err := service.GetEnvironment(
		context.Background(),
		"environment-ref",
	)
	if err != nil {
		t.Fatalf("GetEnvironment() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"GetEnvironment() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_GetEnvironment_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get failed")

	core.EXPECT().
		GetEnvironment(
			gomock.Any(),
			"environment-ref",
		).
		Return(nil, wantErr)

	_, err := service.GetEnvironment(
		context.Background(),
		"environment-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetEnvironment() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_DeleteEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.Environment{
		ID:   "environment-1",
		Name: "environment",
	}

	core.EXPECT().
		DeleteEnvironment(
			gomock.Any(),
			"environment-ref",
		).
		Return(want, nil)

	got, err := service.DeleteEnvironment(
		context.Background(),
		"environment-ref",
	)
	if err != nil {
		t.Fatalf("DeleteEnvironment() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"DeleteEnvironment() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_DeleteEnvironment_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("delete failed")

	core.EXPECT().
		DeleteEnvironment(
			gomock.Any(),
			"environment-ref",
		).
		Return(nil, wantErr)

	_, err := service.DeleteEnvironment(
		context.Background(),
		"environment-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"DeleteEnvironment() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func writeEnvironmentSpec(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(
		t.TempDir(),
		"environment.yml",
	)

	if err := os.WriteFile(
		path,
		[]byte(content),
		0o600,
	); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	return path
}

func TestService_CreateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeEnvironmentSpec(
		t,
		`
name: test-environment
description: test description
spec: {}
`,
	)

	want := &model.Environment{
		ID:   "environment-1",
		Name: "test-environment",
	}

	core.EXPECT().
		CreateEnvironment(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(
			_ context.Context,
			arg *coreapi.CreateEnvironmentParams,
		) (*model.Environment, error) {
			if arg.Name != "test-environment" {
				t.Errorf(
					"CreateEnvironmentParams.Name = %q, want %q",
					arg.Name,
					"test-environment",
				)
			}

			if arg.Description != "test description" {
				t.Errorf(
					"CreateEnvironmentParams.Description = %q, want %q",
					arg.Description,
					"test description",
				)
			}

			return want, nil
		})

	got, err := service.CreateEnvironment(
		context.Background(),
		specPath,
	)
	if err != nil {
		t.Fatalf("CreateEnvironment() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"CreateEnvironment() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_CreateEnvironment_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	_, err := service.CreateEnvironment(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yml"),
	)
	if err == nil {
		t.Fatal("CreateEnvironment() error = nil, want error")
	}
}

func TestService_ValidateEnvironment(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeEnvironmentSpec(
		t,
		`
name: test-environment
description: test description
spec: {}
`,
	)

	core.EXPECT().
		ValidateEnvironment(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil)

	got, err := service.ValidateEnvironment(
		context.Background(),
		specPath,
	)
	if err != nil {
		t.Fatalf("ValidateEnvironment() error = %v", err)
	}

	if got == nil {
		t.Fatal("ValidateEnvironment() = nil")
	}
}

func TestService_ValidateEnvironment_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	_, err := service.ValidateEnvironment(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yml"),
	)
	if err == nil {
		t.Fatal("ValidateEnvironment() error = nil, want error")
	}
}
