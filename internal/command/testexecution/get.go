package testexecution

import (
	"context"
	"io"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/output/view"
)

type GetHandler struct {
	service internalService
}

func NewGetHandler(service internalService) *GetHandler {
	return &GetHandler{
		service: service,
	}
}

type GetParams struct {
	Reference string
	Output    output.Settings
	Writer    io.Writer
}

func (h *GetHandler) Run(ctx context.Context, params GetParams) error {
	te, err := h.service.GetTestExecution(ctx, params.Reference)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		return nil
	}

	return view.RenderTestExecution(params.Writer, te, params.Output)
}
