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

func TestGetHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		GetTest(gomock.Any(), "test-ref").
		Return(&model.Test{
			Name:        "test-name",
			Description: "test description",
		}, nil)

	var writer bytes.Buffer

	params := GetParams{
		Reference: "test-ref",
		Writer:    &writer,
	}
	params.Output.Format = render.FormatYAML

	err := NewGetHandler(service).
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
}

func TestGetHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		GetTest(gomock.Any(), "test-ref").
		Return(&model.Test{Name: "test"}, nil)

	var writer bytes.Buffer

	params := GetParams{
		Reference: "test-ref",
		Writer:    &writer,
	}
	params.Output.IsQuiet = true

	err := NewGetHandler(service).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if writer.Len() != 0 {
		t.Errorf("output = %q, want empty", writer.String())
	}
}

func TestGetHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("get failed")

	service.EXPECT().
		GetTest(gomock.Any(), "test-ref").
		Return(nil, wantErr)

	err := NewGetHandler(service).
		Run(context.Background(), GetParams{
			Reference: "test-ref",
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}

func TestGetHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		GetTest(gomock.Any(), "test-ref").
		Return(&model.Test{Name: "test"}, nil)

	params := GetParams{
		Reference: "test-ref",
		Writer: failingWriter{
			err: errors.New("write failed"),
		},
	}
	params.Output.Format = render.FormatYAML

	err := NewGetHandler(service).
		Run(context.Background(), params)

	if !errors.Is(err, clierrors.ErrCannotRenderOutput) {
		t.Fatalf(
			"Run() error = %v, want ErrCannotRenderOutput",
			err,
		)
	}
}
