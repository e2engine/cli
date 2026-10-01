package view

import (
	"bytes"
	"testing"
	"time"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func TestRenderTestExecution(t *testing.T) {
	te := newTestTestExecution(
		"test-execution:1234567890abcdef",
		coremodel.ExecutionStatus("passed"),
		time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
		time.Date(2026, 9, 28, 10, 15, 35, 0, time.UTC),
		"environment:abcdef1234567890",
		"payment-demo",
		"test:fedcba0987654321",
		"successful-payment",
	)

	tests := []struct {
		name     string
		te       *coremodel.TestExecution
		settings output.Settings
		want     string
	}{
		{
			name: "nil test execution",
			te:   nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "",
		},
		{
			name: "table",
			te:   te,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            STATUS  STARTED               FINISHED              ENVIRONMENT_ID  ENVIRONMENT_NAME  TEST_ID       TEST_NAME         \n" +
				"------------  ------  --------------------  --------------------  --------------  ----------------  ------------  ------------------\n" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456    payment-demo      fedcba098765  successful-payment\n\n",
		},
		{
			name: "table without headers",
			te:   te,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456  payment-demo  fedcba098765  successful-payment\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTestExecution(
				&buf,
				tt.te,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTestExecution() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTestExecution() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestRenderTestExecutions(t *testing.T) {
	testExecutions := []coremodel.TestExecution{
		*newTestTestExecution(
			"test-execution:1234567890abcdef",
			coremodel.ExecutionStatus("passed"),
			time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
			time.Date(2026, 9, 28, 10, 15, 35, 0, time.UTC),
			"environment:abcdef1234567890",
			"payment-demo",
			"test:fedcba0987654321",
			"successful-payment",
		),
		*newTestTestExecution(
			"test-execution:abcdef1234567890",
			coremodel.ExecutionStatus("failed"),
			time.Date(2026, 9, 28, 11, 20, 0, 0, time.UTC),
			time.Date(2026, 9, 28, 11, 20, 3, 0, time.UTC),
			"environment:1234567890abcdef",
			"payment-demo",
			"test:0123456789abcdef",
			"fraud-rejection",
		),
	}

	tests := []struct {
		name           string
		testExecutions []coremodel.TestExecution
		settings       output.Settings
		want           string
	}{
		{
			name:           "empty",
			testExecutions: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "ID  STATUS  STARTED  FINISHED  ENVIRONMENT_ID  ENVIRONMENT_NAME  TEST_ID  TEST_NAME\n--  ------  -------  --------  --------------  ----------------  -------  ---------\n\n",
		},
		{
			name:           "table",
			testExecutions: testExecutions,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            STATUS  STARTED               FINISHED              ENVIRONMENT_ID  ENVIRONMENT_NAME  TEST_ID       TEST_NAME         \n" +
				"------------  ------  --------------------  --------------------  --------------  ----------------  ------------  ------------------\n" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456    payment-demo      fedcba098765  successful-payment\n" +
				"abcdef123456  failed  2026-09-28T11:20:00Z  2026-09-28T11:20:03Z  1234567890ab    payment-demo      0123456789ab  fraud-rejection   \n\n",
		},
		{
			name:           "table without headers",
			testExecutions: testExecutions,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456  payment-demo  fedcba098765  successful-payment\n" +
				"abcdef123456  failed  2026-09-28T11:20:00Z  2026-09-28T11:20:03Z  1234567890ab  payment-demo  0123456789ab  fraud-rejection\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTestExecutions(
				&buf,
				tt.testExecutions,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTestExecutions() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTestExecutions() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func newTestTestExecution(
	id string,
	status coremodel.ExecutionStatus,
	startedAt time.Time,
	finishedAt time.Time,
	environmentID string,
	environmentName string,
	testID string,
	testName string,
) *coremodel.TestExecution {
	return &coremodel.TestExecution{
		ID:              id,
		Status:          status,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		EnvironmentID:   environmentID,
		EnvironmentName: environmentName,
		TestID:          testID,
		TestName:        testName,
	}
}
