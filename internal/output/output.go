package output

import (
	"io"
	"strings"

	"github.com/ygrebnov/model/validation"
	"github.com/ygrebnov/render"
)

type Settings struct {
	Format      render.Format
	IsQuiet     bool
	IsVerbose   bool
	IsNoHeaders bool
}

func NewSettings(format string, quiet, verbose, noHeaders bool) Settings {
	return Settings{
		Format:      render.Format(strings.TrimSpace(format)),
		IsQuiet:     quiet,
		IsVerbose:   verbose,
		IsNoHeaders: noHeaders,
	}
}

type ValidationErrorOutput struct {
	Valid   bool   `json:"valid" yaml:"valid" table:"-"`
	Field   string `json:"field,omitempty" yaml:"field,omitempty" table:"FIELD"`
	Rule    string `json:"rule,omitempty" yaml:"rule,omitempty" table:"RULE"`
	Message string `json:"message" yaml:"message" table:"MESSAGE"`
}

func (v ValidationErrorOutput) TableHeaders() []string {
	return []string{"FIELD", "RULE", "MESSAGE"}
}

func (v ValidationErrorOutput) TableRows() [][]string {
	return [][]string{{v.Field, v.Rule, v.Message}}
}

func RenderValidationError(
	w io.Writer,
	err *validation.Error,
	s Settings,
) error {
	if err == nil || err.Len() == 0 {
		return nil
	}

	issues := err.ByField()

	records := make([]ValidationErrorOutput, 0, err.Len())
	for _, field := range err.Fields() {
		for _, issue := range issues[field] {
			records = append(records, ValidationErrorOutput{
				Valid:   false,
				Field:   issue.Path,
				Rule:    issue.Rule,
				Message: issue.Err.Error(),
			})
		}
	}

	opts := make([]render.Option, 0, 1)
	if s.IsNoHeaders {
		opts = append(
			opts,
			render.WithNoHeaders(),
		)
	}

	switch len(records) {
	case 0:
		return nil
	case 1:
		return render.Render(w, records[0], s.Format, opts...)
	default:
		return render.Render(w, records, s.Format, opts...)
	}
}
