package view

import (
	"io"
	"time"

	coremodel "github.com/e2engine/core/model"
	idpkg "github.com/e2engine/core/pkg/id"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func RenderTestExecution(
	writer io.Writer,
	te *coremodel.TestExecution,
	settings output.Settings,
) error {
	if te == nil {
		return nil // defensive
	}

	if settings.Format == render.FormatTable {
		return RenderTestExecutions(
			writer,
			[]coremodel.TestExecution{*te},
			settings,
		)
	}

	return RenderObject(writer, te, settings)
}

func RenderTestExecutions(
	writer io.Writer,
	tes []coremodel.TestExecution,
	settings output.Settings,
) error {
	if settings.Format == render.FormatTable {
		v := newTestExecutionsView(tes)
		return RenderObject(writer, v, settings)
	}

	return RenderObject(writer, tes, settings)
}

type testExecution struct {
	ID         string
	Status     coremodel.ExecutionStatus
	StartedAt  time.Time
	FinishedAt time.Time

	EnvironmentID   string
	EnvironmentName string

	TestID   string
	TestName string
}

// TableHeaders returns headers for table-format test executions list rendering.
func (t testExecution) TableHeaders() []string {
	return []string{
		"ID",
		"STATUS",
		"STARTED",
		"FINISHED",
		"ENVIRONMENT_ID",
		"ENVIRONMENT_NAME",
		"TEST_ID",
		"TEST_NAME",
	}
}

// TableRows constructs rows for table-format test executions list rendering.
func (t testExecution) TableRows() [][]string {
	return [][]string{{
		idpkg.TruncateID(t.ID),
		string(t.Status),
		t.StartedAt.Format(time.RFC3339Nano),
		t.FinishedAt.Format(time.RFC3339Nano),
		idpkg.TruncateID(t.EnvironmentID),
		t.EnvironmentName,
		idpkg.TruncateID(t.TestID),
		t.TestName,
	}}
}

func newTestExecutionView(te *coremodel.TestExecution) testExecution {
	return testExecution{
		ID:              te.ID,
		Status:          te.Status,
		StartedAt:       te.StartedAt,
		FinishedAt:      te.FinishedAt,
		EnvironmentID:   te.EnvironmentID,
		EnvironmentName: te.EnvironmentName,
		TestID:          te.TestID,
		TestName:        te.TestName,
	}
}

func newTestExecutionsView(testExecutions []coremodel.TestExecution) []testExecution {
	result := make([]testExecution, len(testExecutions))
	for i := range testExecutions {
		result[i] = newTestExecutionView(&testExecutions[i])
	}
	return result
}
