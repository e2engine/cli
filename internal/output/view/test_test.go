package view

import (
	"bytes"
	"testing"
	"time"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func TestRenderTest(t *testing.T) {
	test := newTestTest(
		"test:1234567890abcdef",
		"successful-payment",
		"1.0.0",
		time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
		time.Date(2026, 9, 28, 10, 20, 45, 0, time.UTC),
	)

	tests := []struct {
		name     string
		test     *coremodel.Test
		settings output.Settings
		want     string
	}{
		{
			name: "nil test",
			test: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "",
		},
		{
			name: "table",
			test: test,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            NAME                VERSION  CREATED               UPDATED             \n" +
				"------------  ------------------  -------  --------------------  --------------------\n" +
				"1234567890ab  successful-payment  1.0.0    2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n\n",
		},
		{
			name: "table without headers",
			test: test,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  successful-payment  1.0.0  2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n\n",
		},
		{
			name: "json",
			test: test,
			settings: output.Settings{
				Format: render.FormatJSON,
			},
			want: "{\n" +
				"\t\"kind\": \"\",\n" +
				"\t\"id\": \"test:1234567890abcdef\",\n" +
				"\t\"name\": \"successful-payment\",\n" +
				"\t\"version\": \"1.0.0\",\n" +
				"\t\"spec\": {\n" +
				"\t\t\"request\": {},\n" +
				"\t\t\"expect\": {}\n" +
				"\t},\n" +
				"\t\"created_at\": \"2026-09-28T10:15:30Z\",\n" +
				"\t\"updated_at\": \"2026-09-28T10:20:45Z\"\n" +
				"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTest(
				&buf,
				tt.test,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTest() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTest() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestRenderTests(t *testing.T) {
	tests := []coremodel.Test{
		*newTestTest(
			"test:1234567890abcdef",
			"successful-payment",
			"1.0.0",
			time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
			time.Date(2026, 9, 28, 10, 20, 45, 0, time.UTC),
		),
		*newTestTest(
			"test:abcdef1234567890",
			"fraud-rejection",
			"2.0.0",
			time.Date(2026, 9, 28, 11, 0, 0, 0, time.UTC),
			time.Date(2026, 9, 28, 11, 5, 0, 0, time.UTC),
		),
	}

	testCases := []struct {
		name     string
		tests    []coremodel.Test
		settings output.Settings
		want     string
	}{
		{
			name:  "empty",
			tests: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "ID  NAME  VERSION  CREATED  UPDATED\n--  ----  -------  -------  -------\n\n",
		},
		{
			name:  "table",
			tests: tests,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            NAME                VERSION  CREATED               UPDATED             \n" +
				"------------  ------------------  -------  --------------------  --------------------\n" +
				"1234567890ab  successful-payment  1.0.0    2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n" +
				"abcdef123456  fraud-rejection     2.0.0    2026-09-28T11:00:00Z  2026-09-28T11:05:00Z\n\n",
		},
		{
			name:  "table without headers",
			tests: tests,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  successful-payment  1.0.0  2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n" +
				"abcdef123456  fraud-rejection  2.0.0  2026-09-28T11:00:00Z  2026-09-28T11:05:00Z\n\n",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTests(
				&buf,
				tt.tests,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTests() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTests() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func newTestTest(
	id string,
	name string,
	version string,
	createdAt time.Time,
	updatedAt time.Time,
) *coremodel.Test {
	return (*coremodel.Test)(&coremodel.Resource[coremodel.TestSpec]{
		ID:        id,
		Name:      name,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
}
