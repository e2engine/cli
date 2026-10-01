package test

import (
	"context"
	"fmt"
	"io"

	"github.com/e2engine/cli/internal/output"
)

type RunHandler struct {
	service internalService
}

func NewRunHandler(service internalService) *RunHandler {
	return &RunHandler{
		service: service,
	}
}

type RunParams struct {
	EnvironmentReference string
	TestReference        string
	Output               output.Settings
	Writer               io.Writer
}

func (h *RunHandler) Run(ctx context.Context, params RunParams) error {
	testExecutionID, err := h.service.RunTest(ctx, params.EnvironmentReference, params.TestReference)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		_, err = fmt.Fprintln(params.Writer, testExecutionID)
		return err
	}

	_, err = fmt.Fprintln(params.Writer,
		"created test execution with id:", testExecutionID,
	)
	// Verbose output is the same as default output for run commands.

	return err
}
