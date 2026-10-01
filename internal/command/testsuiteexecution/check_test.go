package testsuiteexecution

import (
	"bytes"
	"errors"
	"strings"
	"testing"
	"time"

	"github.com/e2engine/core/model"
	"github.com/ygrebnov/render"
	"go.uber.org/mock/gomock"

	"github.com/e2engine/cli/internal/output"
	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestCheckHandler_Run(t *testing.T) {
	tests := []struct {
		name string

		status model.ExecutionStatus
		quiet  bool

		expectedError error
		checkOutput   func(*testing.T, string)
	}{
		{
			name:   "passed",
			status: model.ExecutionStatusPassed,
			checkOutput: func(t *testing.T, output string) {
				t.Helper()

				if !strings.Contains(
					output,
					string(model.ExecutionStatusPassed),
				) {
					t.Errorf(
						"output = %q, want status %q",
						output,
						model.ExecutionStatusPassed,
					)
				}
			},
		},
		{
			name:   "passed quiet",
			status: model.ExecutionStatusPassed,
			quiet:  true,
			checkOutput: func(t *testing.T, output string) {
				t.Helper()

				if output != "" {
					t.Errorf("output = %q, want empty", output)
				}
			},
		},
		{
			name:          "failed",
			status:        model.ExecutionStatusFailed,
			expectedError: clierrors.ErrExecutionNotPassed,
		},
		{
			name:          "error",
			status:        model.ExecutionStatusError,
			expectedError: clierrors.ErrExecutionNotPassed,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockinternalService(ctrl)

			service.EXPECT().
				GetTestSuiteExecutionStatus(
					gomock.Any(),
					"execution-ref",
				).
				Return(tt.status, nil)

			var writer bytes.Buffer

			params := CheckParams{
				Reference: "execution-ref",
				Timeout:   time.Second,
				Output: output.Settings{
					IsQuiet: tt.quiet,
					Format:  render.FormatYAML,
				},
				Writer: &writer,
			}

			err := NewCheckHandler(service).
				Run(t.Context(), params)

			if tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf(
						"Run() error = %v, want %v",
						err,
						tt.expectedError,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf("Run() error = %v", err)
			}

			if tt.checkOutput != nil {
				tt.checkOutput(t, writer.String())
			}
		})
	}
}

func TestCheckHandler_Run_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("get status failed")

	service.EXPECT().
		GetTestSuiteExecutionStatus(
			gomock.Any(),
			"execution-ref",
		).
		Return(model.ExecutionStatus(""), wantErr)

	err := NewCheckHandler(service).
		Run(t.Context(), CheckParams{
			Reference: "execution-ref",
			Timeout:   time.Second,
		})

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"Run() error = %v, want %v",
			err,
			wantErr,
		)
	}
}

func TestCheckHandler_Run_RenderError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		GetTestSuiteExecutionStatus(
			gomock.Any(),
			"execution-ref",
		).
		Return(model.ExecutionStatusPassed, nil)

	params := CheckParams{
		Reference: "execution-ref",
		Timeout:   time.Second,
		Writer: failingWriter{
			err: errors.New("write failed"),
		},
	}
	params.Output.Format = render.FormatYAML

	err := NewCheckHandler(service).
		Run(t.Context(), params)

	if !errors.Is(err, clierrors.ErrCannotRenderOutput) {
		t.Fatalf(
			"Run() error = %v, want ErrCannotRenderOutput",
			err,
		)
	}
}

func TestWaitForStatus(t *testing.T) {
	tests := []struct {
		name string

		statuses []model.ExecutionStatus

		expectedStatus model.ExecutionStatus
		expectedError  error
	}{
		{
			name: "passed immediately",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatusPassed,
			},
			expectedStatus: model.ExecutionStatusPassed,
		},
		{
			name: "failed immediately",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatusFailed,
			},
			expectedStatus: model.ExecutionStatusFailed,
		},
		{
			name: "error immediately",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatusError,
			},
			expectedStatus: model.ExecutionStatusError,
		},
		{
			name: "scheduled then passed",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatusScheduled,
				model.ExecutionStatusPassed,
			},
			expectedStatus: model.ExecutionStatusPassed,
		},
		{
			name: "running then passed",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatusRunning,
				model.ExecutionStatusPassed,
			},
			expectedStatus: model.ExecutionStatusPassed,
		},
		{
			name: "scheduled running then failed",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatusScheduled,
				model.ExecutionStatusRunning,
				model.ExecutionStatusFailed,
			},
			expectedStatus: model.ExecutionStatusFailed,
		},
		{
			name: "unexpected status",
			statuses: []model.ExecutionStatus{
				model.ExecutionStatus("unexpected"),
			},
			expectedError: clierrors.ErrUnexpectedExecutionStatus,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			service := NewMockinternalService(ctrl)

			for _, status := range tt.statuses {
				service.EXPECT().
					GetTestSuiteExecutionStatus(
						gomock.Any(),
						"execution-ref",
					).
					Return(status, nil)
			}

			status, err := waitForStatus(
				t.Context(),
				service,
				"execution-ref",
				time.Second,
			)

			if tt.expectedError != nil {
				if !errors.Is(err, tt.expectedError) {
					t.Fatalf(
						"waitForStatus() error = %v, want %v",
						err,
						tt.expectedError,
					)
				}

				if status != "" {
					t.Fatalf(
						"waitForStatus() status = %q, want empty",
						status,
					)
				}

				return
			}

			if err != nil {
				t.Fatalf(
					"waitForStatus() error = %v",
					err,
				)
			}

			if status != tt.expectedStatus {
				t.Errorf(
					"waitForStatus() status = %q, want %q",
					status,
					tt.expectedStatus,
				)
			}
		})
	}
}

func TestWaitForStatus_ServiceError(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	wantErr := errors.New("get status failed")

	service.EXPECT().
		GetTestSuiteExecutionStatus(
			gomock.Any(),
			"execution-ref",
		).
		Return(model.ExecutionStatus(""), wantErr)

	status, err := waitForStatus(
		t.Context(),
		service,
		"execution-ref",
		time.Second,
	)

	if !errors.Is(err, wantErr) {
		t.Fatalf(
			"waitForStatus() error = %v, want %v",
			err,
			wantErr,
		)
	}

	if status != "" {
		t.Errorf(
			"waitForStatus() status = %q, want empty",
			status,
		)
	}
}

func TestWaitForStatus_Timeout(t *testing.T) {
	ctrl := gomock.NewController(t)
	service := NewMockinternalService(ctrl)

	service.EXPECT().
		GetTestSuiteExecutionStatus(
			gomock.Any(),
			"execution-ref",
		).
		Return(model.ExecutionStatusRunning, nil).
		AnyTimes()

	status, err := waitForStatus(
		t.Context(),
		service,
		"execution-ref",
		20*time.Millisecond,
	)

	if !errors.Is(err, clierrors.ErrExecutionCheckTimeout) {
		t.Fatalf(
			"waitForStatus() error = %v, want ErrExecutionCheckTimeout",
			err,
		)
	}

	if status != "" {
		t.Errorf(
			"waitForStatus() status = %q, want empty",
			status,
		)
	}
}
