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

	"github.com/e2engine/cli/internal/service"
	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestListHandler_Run(t *testing.T) {
	ctrl := gomock.NewController(t)
	internalService := NewMockinternalService(ctrl)

	internalService.EXPECT().
		ListTests(
			gomock.Any(),
			&service.ListTestsParams{
				Limit:          10,
				OrderBy:        model.OrderByName,
				OrderDirection: model.OrderDirectionDesc,
			},
		).
		Return([]model.Test{
			{Name: "test-1"},
			{Name: "test-2"},
		}, nil)

	var writer bytes.Buffer

	params := &ListParams{
		Limit:          10,
		OrderBy:        string(model.OrderByName),
		OrderDirection: string(model.OrderDirectionDesc),
		Writer:         &writer,
	}
	params.Output.Format = render.FormatYAML

	err := NewListHandler(internalService).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(writer.String(), "test-1") {
		t.Errorf(
			"output does not contain test-1: %q",
			writer.String(),
		)
	}

	if !strings.Contains(writer.String(), "test-2") {
		t.Errorf(
			"output does not contain test-2: %q",
			writer.String(),
		)
	}
}

func TestListHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	internalService := NewMockinternalService(ctrl)

	internalService.EXPECT().
		ListTests(
			gomock.Any(),
			&service.ListTestsParams{},
		).
		Return([]model.Test{{Name: "test"}}, nil)

	var writer bytes.Buffer

	params := &ListParams{
		Writer: &writer,
	}
	params.Output.IsQuiet = true

	err := NewListHandler(internalService).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if writer.Len() != 0 {
		t.Errorf("output = %q, want empty", writer.String())
	}
}

func TestListHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	internalService := NewMockinternalService(ctrl)

	wantErr := errors.New("list failed")

	internalService.EXPECT().
		ListTests(gomock.Any(), gomock.Any()).
		Return(nil, wantErr)

	err := NewListHandler(internalService).
		Run(context.Background(), &ListParams{})

	if !errors.Is(err, wantErr) {
		t.Fatalf("Run() error = %v, want %v", err, wantErr)
	}
}

func TestListHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	internalService := NewMockinternalService(ctrl)

	internalService.EXPECT().
		ListTests(gomock.Any(), gomock.Any()).
		Return([]model.Test{{Name: "test"}}, nil)

	params := &ListParams{
		Writer: failingWriter{
			err: errors.New("write failed"),
		},
	}
	params.Output.Format = render.FormatYAML

	err := NewListHandler(internalService).
		Run(context.Background(), params)

	if !errors.Is(err, clierrors.ErrCannotRenderOutput) {
		t.Fatalf(
			"Run() error = %v, want ErrCannotRenderOutput",
			err,
		)
	}
}
