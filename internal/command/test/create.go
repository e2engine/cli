package test

import (
	"context"
	"fmt"
	"io"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/output/view"
)

type CreateHandler struct {
	service internalService
}

func NewCreateHandler(service internalService) *CreateHandler {
	return &CreateHandler{
		service: service,
	}
}

type CreateParams struct {
	SpecPath string
	Output   output.Settings
	Writer   io.Writer
}

func (h *CreateHandler) Run(ctx context.Context, params CreateParams) error {
	test, err := h.service.CreateTest(ctx, params.SpecPath)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		_, err = fmt.Fprintln(params.Writer, test.ID)
		return err
	}

	if params.Output.IsVerbose {
		return view.RenderTest(
			params.Writer,
			test,
			params.Output,
		)
	}

	_, err = fmt.Fprintln(params.Writer,
		"created test, name:", test.Name,
		"id:", test.ID,
		"version:", test.Version,
	)

	return err
}
