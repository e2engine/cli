package testexecution

import (
	"context"
	"io"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/output/view"
	"github.com/e2engine/cli/internal/service"
)

type ListHandler struct {
	service internalService
}

func NewListHandler(s internalService) *ListHandler {
	return &ListHandler{
		service: s,
	}
}

type ListParams struct {
	Limit          int
	OrderBy        string
	OrderDirection string
	Output         output.Settings
	Writer         io.Writer
}

func (h *ListHandler) Run(ctx context.Context, params *ListParams) error {
	arg := &service.ListTestExecutionsParams{
		Limit:          params.Limit,
		OrderBy:        model.OrderBy(params.OrderBy),
		OrderDirection: model.OrderDirection(params.OrderDirection),
	}

	tes, err := h.service.ListTestExecutions(ctx, arg)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		return nil
	}

	return view.RenderTestExecutions(params.Writer, tes, params.Output)
}
