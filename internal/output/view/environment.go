package view

import (
	"io"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func RenderEnvironment(
	writer io.Writer,
	environment *coremodel.Environment,
	settings output.Settings,
) error {
	if environment == nil {
		return nil // defensive
	}

	if settings.Format == render.FormatTable {
		return RenderEnvironments(
			writer,
			[]coremodel.Environment{*environment},
			settings,
		)
	}

	return RenderObject(writer, environment, settings)
}

func RenderEnvironments(
	writer io.Writer,
	environments []coremodel.Environment,
	settings output.Settings,
) error {
	if settings.Format == render.FormatTable {
		v := newEnvironmentsView(environments)
		return RenderObject(writer, v, settings)
	}

	return RenderObject(writer, environments, settings)
}

func newEnvironmentView(env *coremodel.Environment) resource {
	return newResourceView(
		(coremodel.Resource[coremodel.EnvironmentSpec])(*env),
	)
}

func newEnvironmentsView(environments []coremodel.Environment) []resource {
	result := make([]resource, len(environments))
	for i := range environments {
		result[i] = newEnvironmentView(&environments[i])
	}
	return result
}
