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

func TestService_ListTests(t *testing.T) {
	tests := []struct {
		name     string
		params   *ListTestsParams
		wantSize int
	}{
		{
			name: "default limit",
			params: &ListTestsParams{
				OrderBy:        model.OrderByID,
				OrderDirection: model.OrderDirectionAsc,
			},
			wantSize: 50,
		},
		{
			name: "smaller limit",
			params: &ListTestsParams{
				Limit:          10,
				OrderBy:        model.OrderByName,
				OrderDirection: model.OrderDirectionDesc,
			},
			wantSize: 10,
		},
		{
			name: "zero limit",
			params: &ListTestsParams{
				Limit: 0,
			},
			wantSize: 50,
		},
		{
			name: "negative limit",
			params: &ListTestsParams{
				Limit: -1,
			},
			wantSize: 50,
		},
		{
			name: "limit equal to maximum",
			params: &ListTestsParams{
				Limit: 50,
			},
			wantSize: 50,
		},
		{
			name: "limit above maximum",
			params: &ListTestsParams{
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

			want := []model.Test{
				{
					ID:   "test-1",
					Name: "test",
				},
			}

			core.EXPECT().
				GetTestsPage(
					gomock.Any(),
					&coreapi.GetTestsPageParams{
						PageSize:       test.wantSize,
						OrderBy:        test.params.OrderBy,
						OrderDirection: test.params.OrderDirection,
					},
				).
				Return(
					&model.TestsPage{
						Items: want,
					},
					nil,
				)

			got, err := service.ListTests(
				context.Background(),
				test.params,
			)
			if err != nil {
				t.Fatalf("ListTests() error = %v", err)
			}

			if len(got) != 1 {
				t.Fatalf(
					"ListTests() returned %d tests, want 1",
					len(got),
				)
			}

			if got[0].ID != want[0].ID {
				t.Errorf(
					"ListTests()[0].ID = %q, want %q",
					got[0].ID,
					want[0].ID,
				)
			}
		})
	}
}

func TestService_ListTests_Error(t *testing.T) {
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
		GetTestsPage(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.ListTests(
		context.Background(),
		&ListTestsParams{},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"ListTests() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_CreateTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeTestSpec(
		t,
		`
name: test-name
description: test description
spec: {}
`,
	)

	want := &model.Test{
		ID:   "test-1",
		Name: "test-name",
	}

	core.EXPECT().
		CreateTest(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(
			_ context.Context,
			arg *coreapi.CreateTestParams,
		) (*model.Test, error) {
			if arg.Name != "test-name" {
				t.Errorf(
					"CreateTestParams.Name = %q, want %q",
					arg.Name,
					"test-name",
				)
			}

			if arg.Description != "test description" {
				t.Errorf(
					"CreateTestParams.Description = %q, want %q",
					arg.Description,
					"test description",
				)
			}

			return want, nil
		})

	got, err := service.CreateTest(
		context.Background(),
		specPath,
	)
	if err != nil {
		t.Fatalf("CreateTest() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"CreateTest() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_CreateTest_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	_, err := service.CreateTest(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yml"),
	)
	if err == nil {
		t.Fatal("CreateTest() error = nil, want error")
	}
}

func TestService_CreateTest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeTestSpec(
		t,
		`
name: test-name
description: test description
spec: {}
`,
	)

	wantErr := errors.New("create failed")

	core.EXPECT().
		CreateTest(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.CreateTest(
		context.Background(),
		specPath,
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"CreateTest() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.Test{
		ID:   "test-1",
		Name: "test",
	}

	core.EXPECT().
		GetTest(
			gomock.Any(),
			"test-ref",
		).
		Return(want, nil)

	got, err := service.GetTest(
		context.Background(),
		"test-ref",
	)
	if err != nil {
		t.Fatalf("GetTest() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"GetTest() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_GetTest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get failed")

	core.EXPECT().
		GetTest(
			gomock.Any(),
			"test-ref",
		).
		Return(nil, wantErr)

	_, err := service.GetTest(
		context.Background(),
		"test-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetTest() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_DeleteTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.Test{
		ID:   "test-1",
		Name: "test",
	}

	core.EXPECT().
		DeleteTest(
			gomock.Any(),
			"test-ref",
		).
		Return(want, nil)

	got, err := service.DeleteTest(
		context.Background(),
		"test-ref",
	)
	if err != nil {
		t.Fatalf("DeleteTest() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"DeleteTest() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_DeleteTest_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("delete failed")

	core.EXPECT().
		DeleteTest(
			gomock.Any(),
			"test-ref",
		).
		Return(nil, wantErr)

	_, err := service.DeleteTest(
		context.Background(),
		"test-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"DeleteTest() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_ValidateTest(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeTestSpec(
		t,
		`
name: test-name
description: test description
spec: {}
`,
	)

	core.EXPECT().
		ValidateTest(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil)

	got, err := service.ValidateTest(
		context.Background(),
		specPath,
	)
	if err != nil {
		t.Fatalf("ValidateTest() error = %v", err)
	}

	if got == nil {
		t.Fatal("ValidateTest() = nil")
	}
}

func TestService_ValidateTest_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	_, err := service.ValidateTest(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yml"),
	)
	if err == nil {
		t.Fatal("ValidateTest() error = nil, want error")
	}
}

func writeTestSpec(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(
		t.TempDir(),
		"test.yml",
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
