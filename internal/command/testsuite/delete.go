package testsuite

import (
	"context"
	"fmt"
	"io"

	"github.com/e2engine/cli/internal/output"
)

type DeleteHandler struct {
	service internalService
}

func NewDeleteHandler(service internalService) *DeleteHandler {
	return &DeleteHandler{
		service: service,
	}
}

type DeleteParams struct {
	Reference string
	Output    output.Settings
	Writer    io.Writer
}

func (h *DeleteHandler) Run(ctx context.Context, params DeleteParams) error {
	ts, err := h.service.DeleteTestSuite(ctx, params.Reference)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		return nil
	}

	_, err = fmt.Fprintln(params.Writer,
		"deleted testsuite, name:", ts.Name,
		"id:", ts.ID,
		"version:", ts.Version,
	)

	return err
}
