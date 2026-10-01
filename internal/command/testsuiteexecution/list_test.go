package testsuiteexecution

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
		ListTestSuiteExecutions(
			gomock.Any(),
			&service.ListTestSuiteExecutionsParams{
				Limit:          10,
				OrderBy:        model.OrderByCreatedAt,
				OrderDirection: model.OrderDirectionDesc,
			},
		).
		Return([]model.TestSuiteExecution{
			{ID: "execution-1"},
			{ID: "execution-2"},
		}, nil)

	var writer bytes.Buffer

	params := &ListParams{
		Limit:          10,
		OrderBy:        string(model.OrderByCreatedAt),
		OrderDirection: string(model.OrderDirectionDesc),
		Writer:         &writer,
	}
	params.Output.Format = render.FormatYAML

	err := NewListHandler(internalService).
		Run(context.Background(), params)
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	if !strings.Contains(writer.String(), "execution-1") {
		t.Errorf(
			"output does not contain execution-1: %q",
			writer.String(),
		)
	}

	if !strings.Contains(writer.String(), "execution-2") {
		t.Errorf(
			"output does not contain execution-2: %q",
			writer.String(),
		)
	}
}

func TestListHandler_Run_Quiet(t *testing.T) {
	ctrl := gomock.NewController(t)
	internalService := NewMockinternalService(ctrl)

	internalService.EXPECT().
		ListTestSuiteExecutions(
			gomock.Any(),
			&service.ListTestSuiteExecutionsParams{},
		).
		Return([]model.TestSuiteExecution{
			{ID: "execution-id"},
		}, nil)

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
		ListTestSuiteExecutions(gomock.Any(), gomock.Any()).
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
		ListTestSuiteExecutions(gomock.Any(), gomock.Any()).
		Return([]model.TestSuiteExecution{
			{ID: "execution-id"},
		}, nil)

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
