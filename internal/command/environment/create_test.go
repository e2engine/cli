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

func TestCreateHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateEnvironment(gomock.Any(), "environment.yaml").
		Return(&model.Environment{
			ID:      "environment-id",
			Name:    "environment-name",
			Version: "1.0.0",
		}, nil)

	handler := NewCreateHandler(service)

	var writer bytes.Buffer

	err := handler.Run(context.Background(), CreateParams{
		SpecPath: "environment.yaml",
		Writer:   &writer,
	})
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	got := writer.String()

	if !strings.Contains(got, "created environment") {
		t.Errorf("output = %q, want creation message", got)
	}
	if !strings.Contains(got, "environment-name") {
		t.Errorf("output = %q, want environment name", got)
	}
	if !strings.Contains(got, "environment-id") {
		t.Errorf("output = %q, want environment ID", got)
	}
}

func TestCreateHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateEnvironment(gomock.Any(), "environment.yaml").
		Return(&model.Environment{
			ID:   "environment-id",
			Name: "environment-name",
		}, nil)

	handler := NewCreateHandler(service)

	var writer bytes.Buffer

	params := CreateParams{
		SpecPath: "environment.yaml",
		Writer:   &writer,
	}
	params.Output.IsQuiet = true

	err := handler.Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if got, want := writer.String(), "environment-id\n"; got != want {
		t.Errorf("output = %q, want %q", got, want)
	}
}

func TestCreateHandler_Run_Verbose(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateEnvironment(gomock.Any(), "environment.yaml").
		Return(&model.Environment{
			ID:          "environment-id",
			Name:        "environment-name",
			Description: "environment description",
		}, nil)

	handler := NewCreateHandler(service)

	var writer bytes.Buffer

	params := CreateParams{
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
		t.Errorf("output does not contain description: %q", writer.String())
	}
}

func TestCreateHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("create failed")

	service.EXPECT().
		CreateEnvironment(gomock.Any(), "environment.yaml").
		Return(nil, wantErr)

	handler := NewCreateHandler(service)

	err := handler.Run(context.Background(), CreateParams{
		SpecPath: "environment.yaml",
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}

func TestCreateHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		CreateEnvironment(gomock.Any(), "environment.yaml").
		Return(&model.Environment{Name: "environment"}, nil)

	handler := NewCreateHandler(service)

	params := CreateParams{
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
