package testsuite

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
	TestSuiteReference   string
	Output               output.Settings
	Writer               io.Writer
}

func (h *RunHandler) Run(ctx context.Context, params RunParams) error {
	testSuiteExecutionID, err := h.service.RunTestSuite(
		ctx,
		params.EnvironmentReference,
		params.TestSuiteReference,
	)
	if err != nil {
		return err
	}

	if params.Output.IsQuiet {
		_, err = fmt.Fprintln(params.Writer, testSuiteExecutionID)
		return err
	}

	_, err = fmt.Fprintln(params.Writer,
		"created testsuite execution with id:", testSuiteExecutionID,
	)
	// Verbose output is the same as default output for run commands.

	return err
}
