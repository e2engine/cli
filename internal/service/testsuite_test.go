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

func TestService_ListTestSuites(t *testing.T) {
	tests := []struct {
		name     string
		params   *ListTestSuitesParams
		wantSize int
	}{
		{
			name: "default limit",
			params: &ListTestSuitesParams{
				OrderBy:        model.OrderByID,
				OrderDirection: model.OrderDirectionAsc,
			},
			wantSize: 50,
		},
		{
			name: "smaller limit",
			params: &ListTestSuitesParams{
				Limit:          10,
				OrderBy:        model.OrderByName,
				OrderDirection: model.OrderDirectionDesc,
			},
			wantSize: 10,
		},
		{
			name: "zero limit",
			params: &ListTestSuitesParams{
				Limit: 0,
			},
			wantSize: 50,
		},
		{
			name: "negative limit",
			params: &ListTestSuitesParams{
				Limit: -1,
			},
			wantSize: 50,
		},
		{
			name: "limit equal to maximum",
			params: &ListTestSuitesParams{
				Limit: 50,
			},
			wantSize: 50,
		},
		{
			name: "limit above maximum",
			params: &ListTestSuitesParams{
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

			want := []model.TestSuite{
				{
					ID:   "testsuite-1",
					Name: "testsuite",
				},
			}

			core.EXPECT().
				GetTestSuitesPage(
					gomock.Any(),
					&coreapi.GetTestSuitesPageParams{
						PageSize:       test.wantSize,
						OrderBy:        test.params.OrderBy,
						OrderDirection: test.params.OrderDirection,
					},
				).
				Return(
					&model.TestSuitesPage{
						Items: want,
					},
					nil,
				)

			got, err := service.ListTestSuites(
				context.Background(),
				test.params,
			)
			if err != nil {
				t.Fatalf("ListTestSuites() error = %v", err)
			}

			if len(got) != 1 {
				t.Fatalf(
					"ListTestSuites() returned %d test suites, want 1",
					len(got),
				)
			}

			if got[0].ID != want[0].ID {
				t.Errorf(
					"ListTestSuites()[0].ID = %q, want %q",
					got[0].ID,
					want[0].ID,
				)
			}
		})
	}
}

func TestService_ListTestSuites_Error(t *testing.T) {
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
		GetTestSuitesPage(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.ListTestSuites(
		context.Background(),
		&ListTestSuitesParams{},
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"ListTestSuites() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_CreateTestSuite(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeTestSuiteSpec(
		t,
		`
name: testsuite-name
description: testsuite description
spec: {}
`,
	)

	want := &model.TestSuite{
		ID:   "testsuite-1",
		Name: "testsuite-name",
	}

	core.EXPECT().
		CreateTestSuite(
			gomock.Any(),
			gomock.Any(),
		).
		DoAndReturn(func(
			_ context.Context,
			arg *coreapi.CreateTestSuiteParams,
		) (*model.TestSuite, error) {
			if arg.Name != "testsuite-name" {
				t.Errorf(
					"CreateTestSuiteParams.Name = %q, want %q",
					arg.Name,
					"testsuite-name",
				)
			}

			if arg.Description != "testsuite description" {
				t.Errorf(
					"CreateTestSuiteParams.Description = %q, want %q",
					arg.Description,
					"testsuite description",
				)
			}

			return want, nil
		})

	got, err := service.CreateTestSuite(
		context.Background(),
		specPath,
	)
	if err != nil {
		t.Fatalf("CreateTestSuite() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"CreateTestSuite() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_CreateTestSuite_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	_, err := service.CreateTestSuite(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yml"),
	)
	if err == nil {
		t.Fatal("CreateTestSuite() error = nil, want error")
	}
}

func TestService_CreateTestSuite_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeTestSuiteSpec(
		t,
		`
name: testsuite-name
description: testsuite description
spec: {}
`,
	)

	wantErr := errors.New("create failed")

	core.EXPECT().
		CreateTestSuite(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil, wantErr)

	_, err := service.CreateTestSuite(
		context.Background(),
		specPath,
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"CreateTestSuite() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_GetTestSuite(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.TestSuite{
		ID:   "testsuite-1",
		Name: "testsuite",
	}

	core.EXPECT().
		GetTestSuite(
			gomock.Any(),
			"testsuite-ref",
		).
		Return(want, nil)

	got, err := service.GetTestSuite(
		context.Background(),
		"testsuite-ref",
	)
	if err != nil {
		t.Fatalf("GetTestSuite() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"GetTestSuite() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_GetTestSuite_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("get failed")

	core.EXPECT().
		GetTestSuite(
			gomock.Any(),
			"testsuite-ref",
		).
		Return(nil, wantErr)

	_, err := service.GetTestSuite(
		context.Background(),
		"testsuite-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"GetTestSuite() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_DeleteTestSuite(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	want := &model.TestSuite{
		ID:   "testsuite-1",
		Name: "testsuite",
	}

	core.EXPECT().
		DeleteTestSuite(
			gomock.Any(),
			"testsuite-ref",
		).
		Return(want, nil)

	got, err := service.DeleteTestSuite(
		context.Background(),
		"testsuite-ref",
	)
	if err != nil {
		t.Fatalf("DeleteTestSuite() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"DeleteTestSuite() = %p, want %p",
			got,
			want,
		)
	}
}

func TestService_DeleteTestSuite_Error(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	wantErr := errors.New("delete failed")

	core.EXPECT().
		DeleteTestSuite(
			gomock.Any(),
			"testsuite-ref",
		).
		Return(nil, wantErr)

	_, err := service.DeleteTestSuite(
		context.Background(),
		"testsuite-ref",
	)
	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"DeleteTestSuite() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestService_ValidateTestSuite(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	specPath := writeTestSuiteSpec(
		t,
		`
name: testsuite-name
description: testsuite description
spec: {}
`,
	)

	core.EXPECT().
		ValidateTestSuite(
			gomock.Any(),
			gomock.Any(),
		).
		Return(nil)

	got, err := service.ValidateTestSuite(
		context.Background(),
		specPath,
	)
	if err != nil {
		t.Fatalf("ValidateTestSuite() error = %v", err)
	}

	if got == nil {
		t.Fatal("ValidateTestSuite() = nil")
	}
}

func TestService_ValidateTestSuite_LoadError(t *testing.T) {
	ctrl := gomock.NewController(t)
	core := NewMockcoreService(ctrl)

	service := &Service{
		core: core,
	}

	_, err := service.ValidateTestSuite(
		context.Background(),
		filepath.Join(t.TempDir(), "missing.yml"),
	)
	if err == nil {
		t.Fatal("ValidateTestSuite() error = nil, want error")
	}
}

func writeTestSuiteSpec(t *testing.T, content string) string {
	t.Helper()

	path := filepath.Join(
		t.TempDir(),
		"testsuite.yml",
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
