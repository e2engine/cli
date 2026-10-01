package testsuite

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
	ts, err := h.service.CreateTestSuite(ctx, params.SpecPath)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		_, err = fmt.Fprintln(params.Writer, ts.ID)
		return err
	}

	if params.Output.IsVerbose {
		return view.RenderTestSuite(
			params.Writer,
			ts,
			params.Output,
		)
	}

	_, err = fmt.Fprintln(params.Writer,
		"created testsuite, name:", ts.Name,
		"id:", ts.ID,
		"version:", ts.Version,
	)

	return err
}
