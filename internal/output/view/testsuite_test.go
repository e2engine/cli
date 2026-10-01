package view

import (
	"bytes"
	"testing"
	"time"

	coremodel "github.com/e2engine/core/model"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func TestRenderTestSuite(t *testing.T) {
	ts := newTestTestSuite(
		"test-suite:1234567890abcdef",
		"smoke",
		"1.0.0",
		time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
		time.Date(2026, 9, 28, 10, 20, 45, 0, time.UTC),
	)

	tests := []struct {
		name     string
		ts       *coremodel.TestSuite
		settings output.Settings
		want     string
	}{
		{
			name: "nil test suite",
			ts:   nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "",
		},
		{
			name: "table",
			ts:   ts,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            NAME   VERSION  CREATED               UPDATED             \n" +
				"------------  -----  -------  --------------------  --------------------\n" +
				"1234567890ab  smoke  1.0.0    2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n\n",
		},
		{
			name: "table without headers",
			ts:   ts,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  smoke  1.0.0  2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n\n",
		},
		{
			name: "json",
			ts:   ts,
			settings: output.Settings{
				Format: render.FormatJSON,
			},
			want: "{\n" +
				"\t\"kind\": \"\",\n" +
				"\t\"id\": \"test-suite:1234567890abcdef\",\n" +
				"\t\"name\": \"smoke\",\n" +
				"\t\"version\": \"1.0.0\",\n" +
				"\t\"spec\": {\n" +
				"\t\t\"selectors\": {}\n" +
				"\t},\n" +
				"\t\"created_at\": \"2026-09-28T10:15:30Z\",\n" +
				"\t\"updated_at\": \"2026-09-28T10:20:45Z\"\n" +
				"}\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTestSuite(
				&buf,
				tt.ts,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTestSuite() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTestSuite() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func TestRenderTestSuites(t *testing.T) {
	testSuites := []coremodel.TestSuite{
		*newTestTestSuite(
			"test-suite:1234567890abcdef",
			"smoke",
			"1.0.0",
			time.Date(2026, 9, 28, 10, 15, 30, 0, time.UTC),
			time.Date(2026, 9, 28, 10, 20, 45, 0, time.UTC),
		),
		*newTestTestSuite(
			"test-suite:abcdef1234567890",
			"regression",
			"2.0.0",
			time.Date(2026, 9, 28, 11, 0, 0, 0, time.UTC),
			time.Date(2026, 9, 28, 11, 5, 0, 0, time.UTC),
		),
	}

	tests := []struct {
		name       string
		testSuites []coremodel.TestSuite
		settings   output.Settings
		want       string
	}{
		{
			name:       "empty",
			testSuites: nil,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "ID  NAME  VERSION  CREATED  UPDATED\n--  ----  -------  -------  -------\n\n",
		},
		{
			name:       "table",
			testSuites: testSuites,
			settings: output.Settings{
				Format: render.FormatTable,
			},
			want: "" +
				"ID            NAME        VERSION  CREATED               UPDATED             \n" +
				"------------  ----------  -------  --------------------  --------------------\n" +
				"1234567890ab  smoke       1.0.0    2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n" +
				"abcdef123456  regression  2.0.0    2026-09-28T11:00:00Z  2026-09-28T11:05:00Z\n\n",
		},
		{
			name:       "table without headers",
			testSuites: testSuites,
			settings: output.Settings{
				Format:      render.FormatTable,
				IsNoHeaders: true,
			},
			want: "" +
				"1234567890ab  smoke  1.0.0  2026-09-28T10:15:30Z  2026-09-28T10:20:45Z\n" +
				"abcdef123456  regression  2.0.0  2026-09-28T11:00:00Z  2026-09-28T11:05:00Z\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var buf bytes.Buffer

			err := RenderTestSuites(
				&buf,
				tt.testSuites,
				tt.settings,
			)
			if err != nil {
				t.Fatalf("RenderTestSuites() error = %v", err)
			}

			if got := buf.String(); got != tt.want {
				t.Errorf(
					"RenderTestSuites() output:\n%q\nwant:\n%q",
					got,
					tt.want,
				)
			}
		})
	}
}

func newTestTestSuite(
	id string,
	name string,
	version string,
	createdAt time.Time,
	updatedAt time.Time,
) *coremodel.TestSuite {
	return (*coremodel.TestSuite)(&coremodel.Resource[coremodel.TestSuiteSpec]{
		ID:        id,
		Name:      name,
		Version:   version,
		CreatedAt: createdAt,
		UpdatedAt: updatedAt,
	})
}
