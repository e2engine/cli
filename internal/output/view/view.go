package view

import (
	"io"

	corekeys "github.com/e2engine/core/pkg/keys"
	"github.com/ygrebnov/errorc"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/pkg/errors"
	"github.com/e2engine/cli/pkg/keys"
)

func RenderObject(writer io.Writer, obj any, settings output.Settings) error {
	opts := make([]render.Option, 0, 1)
	if settings.IsNoHeaders {
		opts = append(opts, render.WithNoHeaders())
	}

	if err := render.Render(
		writer,
		obj,
		settings.Format,
		opts...,
	); err != nil {
		return errorc.With(
			errors.ErrCannotRenderOutput,
			errorc.String(keys.OutputFormat, string(settings.Format)),
			errorc.Error(corekeys.Cause, err),
		)
	}
	return nil
}
