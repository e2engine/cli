package testsuiteexecution

import (
	"context"
	"io"
	"time"

	"github.com/e2engine/core/model"
	"github.com/e2engine/core/pkg/keys"
	"github.com/ygrebnov/errorc"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/output/view"
	"github.com/e2engine/cli/pkg/errors"
)

type CheckHandler struct {
	service internalService
}

func NewCheckHandler(service internalService) *CheckHandler {
	return &CheckHandler{
		service: service,
	}
}

type CheckParams struct {
	Reference string
	Timeout   time.Duration
	Output    output.Settings
	Writer    io.Writer
}

type CheckResult struct {
	Status model.ExecutionStatus `json:"status" yaml:"status"`
}

func (h *CheckHandler) Run(ctx context.Context, params CheckParams) error {
	status, err := waitForStatus(ctx, h.service, params.Reference, params.Timeout)
	if err != nil {
		return err
	}

	switch status {
	case model.ExecutionStatusPassed:
		if params.Output.IsQuiet {
			return nil
		}

		// default and verbose return similar output
		return view.RenderObject(
			params.Writer,
			CheckResult{Status: status},
			params.Output,
		)

	default:
		return errorc.With(
			errors.ErrExecutionNotPassed,
			errorc.String(keys.ExecutionStatus, string(status)),
		)
	}
}

func waitForStatus(
	ctx context.Context,
	service internalService,
	reference string,
	timeout time.Duration,
) (model.ExecutionStatus, error) {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		status, err := service.GetTestSuiteExecutionStatus(ctx, reference)
		if err != nil {
			return "", err
		}

		switch status {
		case model.ExecutionStatusPassed,
			model.ExecutionStatusFailed,
			model.ExecutionStatusError:
			return status, nil

		case model.ExecutionStatusScheduled,
			model.ExecutionStatusRunning:
			// check again after the ticker interval

		default:
			return "", errorc.With(
				errors.ErrUnexpectedExecutionStatus,
				errorc.String(keys.ExecutionStatus, string(status)),
			)
		}

		select {
		case <-ctx.Done():
			return "", errors.ErrExecutionCheckTimeout
		case <-ticker.C:
		}
	}
}
