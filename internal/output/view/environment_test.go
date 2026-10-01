package view

import (
	"bytes"
	"testing"
	"time"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func TestRenderEnvironment(t *testing.T) {
	environment := newTestEnvironment(
		"environment:1234567890abcdef",
		"payment-demo",
		"1.0.0",
		time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
		time.Date(2026, 9, 28, 10, 20, 45, 0, time.UTC),
	)

	tests := []struct {
		name        string
		environment *coremodel.Environment
		settings    output.Settings
		want        string
	}{
		{
			name:        "nil environment",
			environment: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "",
		},
		{
			name:        "table",
			environment: environment,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            NAME          VERSION  CREATED               UPDATED             \n" +
				"------------  ------------  -------  --------------------  --------------------\n" +
				"1234567890ab  payment-demo  1.0.0    2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n\n",
		},
		{
			name:        "table without headers",
			environment: environment,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  payment-demo  1.0.0  2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n\n",
		},
		{
			name:        "json",
			environment: environment,
			settings: output.Settings{
				Format: render.FormatJSON,
			},
			want: "{\n" +
				"\t\"kind\": \"\",\n" +
				"\t\"id\": \"environment:1234567890abcdef\",\n" +
				"\t\"name\": \"payment-demo\",\n" +
				"\t\"version\": \"1.0.0\",\n" +
				"\t\"spec\": {\n\t\t\"services\": null\n\t},\n" +
				"\t\"created_at\": \"2026-09-28T10:15:30Z\",\n" +
				"\t\"updated_at\": \"2026-09-28T10:20:45Z\"\n" +
				"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderEnvironment(
				&buf,
				tt.environment,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderEnvironment() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderEnvironment() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestRenderEnvironments(t *testing.T) {
	environments := []coremodel.Environment{
		*newTestEnvironment(
			"environment:1234567890abcdef",
			"payment-demo",
			"1.0.0",
			time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
			time.Date(2026, 9, 28, 10, 20, 45, 0, time.UTC),
		),
		*newTestEnvironment(
			"environment:abcdef1234567890",
			"other-demo",
			"2.0.0",
			time.Date(2026, 9, 28, 11, 0, 0, 0, time.UTC),
			time.Date(2026, 9, 28, 11, 5, 0, 0, time.UTC),
		),
	}

	tests := []struct {
		name         string
		environments []coremodel.Environment
		settings     output.Settings
		want         string
	}{
		{
			name:         "empty",
			environments: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "ID  NAME  VERSION  CREATED  UPDATED\n--  ----  -------  -------  -------\n\n",
		},
		{
			name:         "table",
			environments: environments,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            NAME          VERSION  CREATED               UPDATED             \n" +
				"------------  ------------  -------  --------------------  --------------------\n" +
				"1234567890ab  payment-demo  1.0.0    2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n" +
				"abcdef123456  other-demo    2.0.0    2026-09-28T11:00:00Z  2026-09-28T11:05:00Z\n\n",
		},
		{
			name:         "table without headers",
			environments: environments,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  payment-demo  1.0.0  2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n" +
				"abcdef123456  other-demo  2.0.0  2026-09-28T11:00:00Z  2026-09-28T11:05:00Z\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderEnvironments(
				&buf,
				tt.environments,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderEnvironments() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderEnvironments() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func newTestEnvironment(
	id string,
	name string,
	version string,
	createdAt time.Time,
	updatedAt time.Time,
) *coremodel.Environment {
	return (*coremodel.Environment)(&coremodel.Resource[coremodel.EnvironmentSpec]{
		ID:        id,
		Name:      name,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
}
