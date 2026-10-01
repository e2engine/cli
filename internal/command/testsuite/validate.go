package testsuite

import (
	"context"
	"fmt"
	"io"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/output/view"
)

type ValidateHandler struct {
	service internalService
}

func NewValidateHandler(service internalService) *ValidateHandler {
	return &ValidateHandler{
		service: service,
	}
}

type ValidateParams struct {
	SpecPath string
	Output   output.Settings
	Writer   io.Writer
}

func (h *ValidateHandler) Run(ctx context.Context, params ValidateParams) error {
	validated, err := h.service.ValidateTestSuite(ctx, params.SpecPath)

	if params.Output.IsQuiet {
		return err
	}

	if err != nil {
		return err
	}

	if !params.Output.IsVerbose {
		_, printErr := fmt.Fprintln(params.Writer, "Testsuite spec is valid")
		return printErr
	}

	return view.RenderTestSuite(params.Writer, validated, params.Output)
}
