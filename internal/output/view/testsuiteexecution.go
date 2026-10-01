package view

import (
	"io"
	"strconv"
	"time"

	coremodel "github.com/e2engine/core/model"
	idpkg "github.com/e2engine/core/pkg/id"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func RenderTestSuiteExecution(
	writer io.Writer,
	tse *coremodel.TestSuiteExecution,
	settings output.Settings,
) error {
	if tse == nil {
		return nil // defensive
	}

	if settings.Format == render.FormatTable {
		return RenderTestSuiteExecutions(
			writer,
			[]coremodel.TestSuiteExecution{*tse},
			settings,
		)
	}

	return RenderObject(writer, tse, settings)
}

func RenderTestSuiteExecutions(
	writer io.Writer,
	tses []coremodel.TestSuiteExecution,
	settings output.Settings,
) error {
	if settings.Format == render.FormatTable {
		v := newTestSuiteExecutionsView(tses)
		return RenderObject(writer, v, settings)
	}

	return RenderObject(writer, tses, settings)
}

type testSuiteExecution struct {
	ID         string
	Status     coremodel.ExecutionStatus
	StartedAt  time.Time
	FinishedAt time.Time

	TestSuiteID   string
	TestSuiteName string

	TestsCount int

	EnvironmentID   string
	EnvironmentName string
}

// TableHeaders returns headers for table-format test suite executions list rendering.
func (t testSuiteExecution) TableHeaders() []string {
	return []string{
		"ID",
		"STATUS",
		"STARTED",
		"FINISHED",
		"TESTSUITE_ID",
		"TESTSUITE_NAME",
		"TESTS",
		"ENVIRONMENT_ID",
		"ENVIRONMENT_NAME",
	}
}

// TableRows constructs rows for table-format test suite executions list rendering.
func (t testSuiteExecution) TableRows() [][]string {
	return [][]string{{
		idpkg.TruncateID(t.ID),
		string(t.Status),
		t.StartedAt.Format(time.RFC3339Nano),
		t.FinishedAt.Format(time.RFC3339Nano),
		idpkg.TruncateID(t.TestSuiteID),
		t.TestSuiteName,
		strconv.Itoa(t.TestsCount),
		idpkg.TruncateID(t.EnvironmentID),
		t.EnvironmentName,
	}}
}

func newTestSuiteExecutionView(tse *coremodel.TestSuiteExecution) testSuiteExecution {
	return testSuiteExecution{
		ID:              tse.ID,
		Status:          tse.Status,
		StartedAt:       tse.StartedAt,
		FinishedAt:      tse.FinishedAt,
		TestSuiteID:     tse.TestSuiteID,
		TestSuiteName:   tse.TestSuiteName,
		TestsCount:      tse.TestsCount,
		EnvironmentID:   tse.EnvironmentID,
		EnvironmentName: tse.EnvironmentName,
	}
}

func newTestSuiteExecutionsView(testSuiteExecutions []coremodel.TestSuiteExecution) []testSuiteExecution {
	result := make([]testSuiteExecution, len(testSuiteExecutions))
	for i := range testSuiteExecutions {
		result[i] = newTestSuiteExecutionView(&testSuiteExecutions[i])
	}
	return result
}
