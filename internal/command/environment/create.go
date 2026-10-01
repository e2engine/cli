package environment

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
	environment, err := h.service.CreateEnvironment(ctx, params.SpecPath)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		_, err = fmt.Fprintln(params.Writer, environment.ID)
		return err
	}

	if params.Output.IsVerbose {
		return view.RenderEnvironment(
			params.Writer,
			environment,
			params.Output,
		)
	}

	_, err = fmt.Fprintln(params.Writer,
		"created environment, name:", environment.Name,
		"id:", environment.ID,
		"version:", environment.Version,
	)

	return err
}
