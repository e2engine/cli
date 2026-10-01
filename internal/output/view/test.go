package view

import (
	"io"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func RenderTest(
	writer io.Writer,
	test *coremodel.Test,
	settings output.Settings,
) error {
	if test == nil {
		return nil // defensive
	}

	if settings.Format == render.FormatTable {
		return RenderTests(
			writer,
			[]coremodel.Test{*test},
			settings,
		)
	}

	return RenderObject(writer, test, settings)
}

func RenderTests(
	writer io.Writer,
	tests []coremodel.Test,
	settings output.Settings,
) error {
	if settings.Format == render.FormatTable {
		v := newTestsView(tests)
		return RenderObject(writer, v, settings)
	}

	return RenderObject(writer, tests, settings)
}

func newTestView(test *coremodel.Test) resource {
	return newResourceView(
		(coremodel.Resource[coremodel.TestSpec])(*test),
	)
}

func newTestsView(tests []coremodel.Test) []resource {
	result := make([]resource, len(tests))
	for i := range tests {
		result[i] = newTestView(&tests[i])
	}
	return result
}
