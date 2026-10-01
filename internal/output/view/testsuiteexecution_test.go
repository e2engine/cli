package view

import (
	"bytes"
	"testing"
	"time"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func TestRenderTestSuiteExecution(t *testing.T) {
	tse := newTestTestSuiteExecution(
		"test-suite-execution:1234567890abcdef",
		coremodel.ExecutionStatus("passed"),
		time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
		time.Date(2026, 9, 28, 10, 15, 35, 0, time.UTC),
		"test-suite:abcdef1234567890",
		"smoke",
		3,
		"environment:fedcba0987654321",
		"payment-demo",
	)

	tests := []struct {
		name     string
		tse      *coremodel.TestSuiteExecution
		settings output.Settings
		want     string
	}{
		{
			name: "nil test suite execution",
			tse:  nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "",
		},
		{
			name: "table",
			tse:  tse,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            STATUS  STARTED               FINISHED              TESTSUITE_ID  TESTSUITE_NAME  TESTS  ENVIRONMENT_ID  ENVIRONMENT_NAME\n" +
				"------------  ------  --------------------  --------------------  ------------  --------------  -----  --------------  ----------------\n" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456  smoke           3      fedcba098765    payment-demo    \n\n",
		},
		{
			name: "table without headers",
			tse:  tse,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456  smoke  3  fedcba098765  payment-demo\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTestSuiteExecution(
				&buf,
				tt.tse,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTestSuiteExecution() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTestSuiteExecution() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestRenderTestSuiteExecutions(t *testing.T) {
	testSuiteExecutions := []coremodel.TestSuiteExecution{
		*newTestTestSuiteExecution(
			"test-suite-execution:1234567890abcdef",
			coremodel.ExecutionStatus("passed"),
			time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
			time.Date(2026, 9, 28, 10, 15, 35, 0, time.UTC),
			"test-suite:abcdef1234567890",
			"smoke",
			3,
			"environment:fedcba0987654321",
			"payment-demo",
		),
		*newTestTestSuiteExecution(
			"test-suite-execution:abcdef1234567890",
			coremodel.ExecutionStatus("failed"),
			time.Date(2026, 9, 28, 11, 20, 0, 0, time.UTC),
			time.Date(2026, 9, 28, 11, 20, 8, 0, time.UTC),
			"test-suite:0123456789abcdef",
			"regression",
			12,
			"environment:1234567890abcdef",
			"staging",
		),
	}

	tests := []struct {
		name                string
		testSuiteExecutions []coremodel.TestSuiteExecution
		settings            output.Settings
		want                string
	}{
		{
			name:                "empty",
			testSuiteExecutions: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "ID  STATUS  STARTED  FINISHED  TESTSUITE_ID  TESTSUITE_NAME  TESTS  ENVIRONMENT_ID  ENVIRONMENT_NAME\n--  ------  -------  --------  ------------  --------------  -----  --------------  ----------------\n\n",
		},
		{
			name:                "table",
			testSuiteExecutions: testSuiteExecutions,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            STATUS  STARTED               FINISHED              TESTSUITE_ID  TESTSUITE_NAME  TESTS  ENVIRONMENT_ID  ENVIRONMENT_NAME\n" +
				"------------  ------  --------------------  --------------------  ------------  --------------  -----  --------------  ----------------\n" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456  smoke           3      fedcba098765    payment-demo    \n" +
				"abcdef123456  failed  2026-09-28T11:20:00Z  2026-09-28T11:20:08Z  0123456789ab  regression      12     1234567890ab    staging         \n\n",
		},
		{
			name:                "table without headers",
			testSuiteExecutions: testSuiteExecutions,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  passed  2026-09-28T10:15:30Z  2026-09-28T10:15:35Z  abcdef123456  smoke  3  fedcba098765  payment-demo\n" +
				"abcdef123456  failed  2026-09-28T11:20:00Z  2026-09-28T11:20:08Z  0123456789ab  regression  12  1234567890ab  staging\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTestSuiteExecutions(
				&buf,
				tt.testSuiteExecutions,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTestSuiteExecutions() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTestSuiteExecutions() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func newTestTestSuiteExecution(
	id string,
	status coremodel.ExecutionStatus,
	startedAt time.Time,
	finishedAt time.Time,
	testSuiteID string,
	testSuiteName string,
	testsCount int,
	environmentID string,
	environmentName string,
) *coremodel.TestSuiteExecution {
	return &coremodel.TestSuiteExecution{
		ID:              id,
		Status:          status,
		StartedAt:       startedAt,
		FinishedAt:      finishedAt,
		TestSuiteID:     testSuiteID,
		TestSuiteName:   testSuiteName,
		TestsCount:      testsCount,
		EnvironmentID:   environmentID,
		EnvironmentName: environmentName,
	}
}
