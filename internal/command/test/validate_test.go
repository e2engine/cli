package test

import (
	"bytes"
	"context"
	"errors"
	"strings"
	"testing"

	"github.com/e2engine/core/model"
	"github.com/ygrebnov/render"
	"go.uber.org/mock/gomock"

	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestValidateHandler_Run(t *testing.T) {
	tests := []struct {
		name       string
		quiet      bool
		serviceErr error
		wantOutput string
		wantErr    bool
	}{
		{
			name:       "valid",
			wantOutput: "Test spec is valid\n",
		},
		{
			name:  "quiet",
			quiet: true,
		},
		{
			name:       "validation error",
			serviceErr: errors.New("validation failed"),
			wantErr:    true,
		},
		{
			name:       "quiet validation error",
			quiet:      true,
			serviceErr: errors.New("validation failed"),
			wantErr:    true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockinternalService(ctrl)

			service.EXPECT().
				ValidateTest(gomock.Any(), "test.yaml").
				Return(&model.Test{}, tt.serviceErr)

			var writer bytes.Buffer

			params := ValidateParams{
				SpecPath: "test.yaml",
				Writer:   &writer,
			}
			params.Output.IsQuiet = tt.quiet

			err := NewValidateHandler(service).
				Run(context.Background(), params)

			if tt.wantErr {
				if !errors.Is(err, tt.serviceErr) {
					t.Fatalf(
						"Run() error = %v, want %v",
						err,
						tt.serviceErr,
					)
				}
			} else if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if got := writer.String(); got != tt.wantOutput {
				t.Errorf(
					"output = %q, want %q",
					got,
					tt.wantOutput,
				)
			}
		})
	}
}

func TestValidateHandler_Run_Verbose(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		ValidateTest(gomock.Any(), "test.yaml").
		Return(&model.Test{
			Name:        "test-name",
			Description: "test description",
		}, nil)

	var writer bytes.Buffer

	params := ValidateParams{
		SpecPath: "test.yaml",
		Writer:   &writer,
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatYAML

	err := NewValidateHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(writer.String(), "test-name") {
		t.Errorf(
			"output does not contain test name: %q",
			writer.String(),
		)
	}

	if !strings.Contains(writer.String(), "test description") {
		t.Errorf(
			"output does not contain description: %q",
			writer.String(),
		)
	}
}

func TestValidateHandler_Run_VerboseTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		ValidateTest(gomock.Any(), "test.yaml").
		Return(&model.Test{}, nil)

	var writer bytes.Buffer

	params := ValidateParams{
		SpecPath: "test.yaml",
		Writer:   &writer,
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatTable

	err := NewValidateHandler(service).
		Run(context.Background(), params)

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestValidateHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		ValidateTest(gomock.Any(), "test.yaml").
		Return(&model.Test{Name: "test"}, nil)

	params := ValidateParams{
		SpecPath: "test.yaml",
		Writer: failingWriter{
			err: errors.New("write failed"),
		},
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatYAML

	err := NewValidateHandler(service).
		Run(context.Background(), params)

	if !errors.Is(err, clierrors.ErrCannotRenderOutput) {
		t.Fatalf(
			"Run() error = %v, want ErrCannotRenderOutput",
			err,
		)
	}
}
