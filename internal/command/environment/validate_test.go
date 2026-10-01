package environment

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
		output     func() ValidateParams
		serviceErr error
		wantOutput string
		wantErr    error
	}{
		{
			name: "valid",
			output: func() ValidateParams {
				return ValidateParams{
					SpecPath: "environment.yaml",
				}
			},
			wantOutput: "Environment spec is valid\n",
		},
		{
			name: "quiet",
			output: func() ValidateParams {
				params := ValidateParams{
					SpecPath: "environment.yaml",
				}
				params.Output.IsQuiet = true
				return params
			},
		},
		{
			name: "validation error",
			output: func() ValidateParams {
				return ValidateParams{
					SpecPath: "environment.yaml",
				}
			},
			serviceErr: errors.New("validation failed"),
			wantErr:    errors.New("validation failed"),
		},
		{
			name: "quiet validation error",
			output: func() ValidateParams {
				params := ValidateParams{
					SpecPath: "environment.yaml",
				}
				params.Output.IsQuiet = true
				return params
			},
			serviceErr: errors.New("validation failed"),
			wantErr:    errors.New("validation failed"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockinternalService(ctrl)

			params := tt.output()
			var writer bytes.Buffer
			params.Writer = &writer

			var serviceErr error
			if tt.serviceErr != nil {
				serviceErr = tt.serviceErr
			}

			service.EXPECT().
				ValidateEnvironment(gomock.Any(), "environment.yaml").
				Return(&model.Environment{}, serviceErr)

			handler := NewValidateHandler(service)

			err := handler.Run(context.Background(), params)

			if tt.wantErr != nil {
				if err == nil || err.Error() != tt.wantErr.Error() {
					t.Fatalf("Run() error = %v, want %v", err, tt.wantErr)
				}
			} else if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if got := writer.String(); got != tt.wantOutput {
				t.Errorf("output = %q, want %q", got, tt.wantOutput)
			}
		})
	}
}

func TestValidateHandler_Run_Verbose(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	environment := &model.Environment{
		Name:        "environment-name",
		Description: "environment description",
	}

	service.EXPECT().
		ValidateEnvironment(gomock.Any(), "environment.yaml").
		Return(environment, nil)

	var writer bytes.Buffer

	handler := NewValidateHandler(service)
	params := ValidateParams{
		SpecPath: "environment.yaml",
		Writer:   &writer,
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatYAML

	err := handler.Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(writer.String(), "environment-name") {
		t.Errorf("output does not contain environment name: %q", writer.String())
	}

	if !strings.Contains(writer.String(), "environment description") {
		t.Errorf("output does not contain environment description: %q", writer.String())
	}
}

func TestValidateHandler_Run_VerboseTable(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		ValidateEnvironment(gomock.Any(), "environment.yaml").
		Return(&model.Environment{}, nil)

	handler := NewValidateHandler(service)

	var writer bytes.Buffer

	params := ValidateParams{
		SpecPath: "environment.yaml",
		Writer:   &writer,
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatTable

	err := handler.Run(context.Background(), params)

	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
}

func TestValidateHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		ValidateEnvironment(gomock.Any(), "environment.yaml").
		Return(&model.Environment{Name: "environment"}, nil)

	handler := NewValidateHandler(service)

	params := ValidateParams{
		SpecPath: "environment.yaml",
		Writer: failingWriter{
			err: errors.New("write failed"),
		},
	}
	params.Output.IsVerbose = true
	params.Output.Format = render.FormatYAML

	err := handler.Run(context.Background(), params)

	if !errors.Is(err, clierrors.ErrCannotRenderOutput) {
		t.Fatalf(
			"Run() error = %v, want ErrCannotRenderOutput",
			err,
		)
	}
}
