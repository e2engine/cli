package view

import (
	"bytes"
	"errors"
	"testing"
	"time"

	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
	clierrors "github.com/e2engine/cli/pkg/errors"
)

func TestRenderObject(t *testing.T) {
	createdAt := time.Date(
		2026, time.September, 28,
		10, 15, 30, 123456789,
		time.UTC,
	)
	updatedAt := time.Date(
		2026, time.September, 28,
		10, 20, 45, 987654321,
		time.UTC,
	)

	obj := resource{
		ID:        "environment:1234567890abcdef",
		Name:      "payment-demo",
		Version:   "1.0.0",
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	}

	const (
		tableOutput = "" +
			"ID            NAME          VERSION  CREATED                         UPDATED                       \n" +
			"------------  ------------  -------  ------------------------------  ------------------------------\n" +
			"1234567890ab  payment-demo  1.0.0    2026-09-28T10:15:30.123456789Z  2026-09-28T10:20:45.987654321Z\n\n"

		tableNoHeadersOutput = "" +
			"1234567890ab  payment-demo  1.0.0  2026-09-28T10:15:30.123456789Z  2026-09-28T10:20:45.987654321Z\n\n"

		jsonOutput = "{\n" +
			"\t\"ID\": \"environment:1234567890abcdef\",\n" +
			"\t\"Name\": \"payment-demo\",\n" +
			"\t\"Version\": \"1.0.0\",\n" +
			"\t\"CreatedAt\": \"2026-09-28T10:15:30.123456789Z\",\n" +
			"\t\"UpdatedAt\": \"2026-09-28T10:20:45.987654321Z\"\n" +
			"}\n"

		yamlOutput = "" +
			"id: environment:1234567890abcdef\n" +
			"name: payment-demo\n" +
			"version: 1.0.0\n" +
			"createdat: 2026-09-28T10:15:30.123456789Z\n" +
			"updatedat: 2026-09-28T10:20:45.987654321Z\n\n"
	)

	tests := []struct {
		name       string
		settings   output.Settings
		want       string
		wantErr    error
		wantErrStr string
	}{
		{
			name: "table",
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: tableOutput,
		},
		{
			name: "table quiet",
			settings: output.Settings{
				Format:  render.FormatTable,
				IsQuiet: true,
			},
			want: tableOutput,
		},
		{
			name: "table verbose",
			settings: output.Settings{
				Format:    render.FormatTable,
				IsVerbose: true,
			},
			want: tableOutput,
		},
		{
			name: "table no headers",
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: tableNoHeadersOutput,
		},
		{
			name: "table quiet verbose no headers",
			settings: output.Settings{
				Format:      render.FormatTable,
				IsQuiet:     true,
				IsVerbose:   true,
				IsNoHeaders: true,
			},
			want: tableNoHeadersOutput,
		},
		{
			name: "json",
			settings: output.Settings{
				Format: render.FormatJSON,
			},
			want: jsonOutput,
		},
		{
			name: "json quiet",
			settings: output.Settings{
				Format:  render.FormatJSON,
				IsQuiet: true,
			},
			want: jsonOutput,
		},
		{
			name: "json verbose",
			settings: output.Settings{
				Format:    render.FormatJSON,
				IsVerbose: true,
			},
			want: jsonOutput,
		},
		{
			name: "json no headers",
			settings: output.Settings{
				Format:      render.FormatJSON,
				IsNoHeaders: true,
			},
			want: jsonOutput,
		},
		{
			name: "json quiet verbose no headers",
			settings: output.Settings{
				Format:      render.FormatJSON,
				IsQuiet:     true,
				IsVerbose:   true,
				IsNoHeaders: true,
			},
			want: jsonOutput,
		},
		{
			name: "yaml",
			settings: output.Settings{
				Format: render.FormatYAML,
			},
			want: yamlOutput,
		},
		{
			name: "yaml quiet",
			settings: output.Settings{
				Format:  render.FormatYAML,
				IsQuiet: true,
			},
			want: yamlOutput,
		},
		{
			name: "yaml verbose",
			settings: output.Settings{
				Format:    render.FormatYAML,
				IsVerbose: true,
			},
			want: yamlOutput,
		},
		{
			name: "yaml no headers",
			settings: output.Settings{
				Format:      render.FormatYAML,
				IsNoHeaders: true,
			},
			want: yamlOutput,
		},
		{
			name: "yaml quiet verbose no headers",
			settings: output.Settings{
				Format:      render.FormatYAML,
				IsQuiet:     true,
				IsVerbose:   true,
				IsNoHeaders: true,
			},
			want: yamlOutput,
		},
		{
			name: "unsupported format",
			settings: output.Settings{
				Format: render.Format("xml"),
			},
			wantErr:    clierrors.ErrCannotRenderOutput,
			wantErrStr: "cannot render output, output.format: xml, cause: unsupported format",
		},
		{
			name: "empty format",
			settings: output.Settings{
				Format: "",
			},
			wantErr:    clierrors.ErrCannotRenderOutput,
			wantErrStr: "cannot render output, output.format: , cause: unsupported format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderObject(&buf, obj, tt.settings)

			if tt.wantErr != nil {
				if err == nil {
					t.Fatalf("RenderObject() error = nil, want %v", tt.wantErr)
				}
				if !errors.Is(err, tt.wantErr) {
					t.Errorf(
						"RenderObject() error = %v, want errors.Is(_, %v)",
						err,
						tt.wantErr,
					)
				}
				if err.Error() != tt.wantErrStr {
					t.Errorf(
						"RenderObject() error = %q, want %q",
						err.Error(),
						tt.wantErrStr,
					)
				}
				if got := buf.String(); got != "" {
					t.Errorf("RenderObject() output = %q, want empty", got)
				}
				return
			}

			if err != nil {
				t.Fatalf("RenderObject() unexpected error: %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderObject() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}
