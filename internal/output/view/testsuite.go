package view

import (
	"io"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func RenderTestSuite(
	writer io.Writer,
	ts *coremodel.TestSuite,
	settings output.Settings,
) error {
	if ts == nil {
		return nil // defensive
	}

	if settings.Format == render.FormatTable {
		return RenderTestSuites(
			writer,
			[]coremodel.TestSuite{*ts},
			settings,
		)
	}

	return RenderObject(writer, ts, settings)
}

func RenderTestSuites(
	writer io.Writer,
	tss []coremodel.TestSuite,
	settings output.Settings,
) error {
	if settings.Format == render.FormatTable {
		v := newTestSuitesView(tss)
		return RenderObject(writer, v, settings)
	}

	return RenderObject(writer, tss, settings)
}

func newTestSuiteView(ts *coremodel.TestSuite) resource {
	return newResourceView(
		(coremodel.Resource[coremodel.TestSuiteSpec])(*ts),
	)
}

func newTestSuitesView(tss []coremodel.TestSuite) []resource {
	result := make([]resource, len(tss))
	for i := range tss {
		result[i] = newTestSuiteView(&tss[i])
	}
	return result
}
