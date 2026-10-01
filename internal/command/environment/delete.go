package environment

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
	environment, err := h.service.DeleteEnvironment(ctx, params.Reference)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		return nil
	}

	_, err = fmt.Fprintln(params.Writer,
		"deleted environment, name:", environment.Name,
		"id:", environment.ID,
		"version:", environment.Version,
	)

	return err
}
