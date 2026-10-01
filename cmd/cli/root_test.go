package main

import (
	"bytes"
	"testing"

	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
)

func newTestStreams() basicStreams {
	return basicStreams{
		in:     bytes.NewReader(nil),
		out:    new(bytes.Buffer),
		errOut: new(bytes.Buffer),
	}
}

func TestNewRootCommand(t *testing.T) {
	streams := newTestStreams()

	cmd := newRootCommand(
		streams,
		providers{},
	)

	if cmd.Use != appName {
		t.Errorf("Use = %q, want %q", cmd.Use, appName)
	}

	if !cmd.SilenceErrors {
		t.Error("SilenceErrors = false, want true")
	}

	if !cmd.SilenceUsage {
		t.Error("SilenceUsage = false, want true")
	}

	for _, name := range []string{
		"output",
		"quiet",
		"no-headers",
		"verbose",
	} {
		if cmd.PersistentFlags().Lookup(name) == nil {
			t.Errorf("persistent flag %q is not registered", name)
		}
	}

	for _, name := range []string{
		"validate",
		"create",
		"get",
		"delete",
		"version",
		"run",
		"internal-runner",
	} {
		if _, _, err := cmd.Find([]string{name}); err != nil {
			t.Errorf("command %q is not registered: %v", name, err)
		}
	}
}

func TestRootCommandOutputSettings(t *testing.T) {
	tests := []struct {
		name      string
		args      []string
		format    render.Format
		quiet     bool
		verbose   bool
		noHeaders bool
	}{
		{
			name: "defaults",
			args: []string{
				"version",
			},
		},
		{
			name: "all flags",
			args: []string{
				"--output", "json",
				"--quiet",
				"--verbose",
				"--no-headers",
				"version",
			},
			format:    render.FormatJSON,
			quiet:     true,
			verbose:   true,
			noHeaders: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRootSettings()

			streams := newTestStreams()
			cmd := newRootCommand(
				streams,
				providers{},
			)

			cmd.SetArgs(tt.args)

			// "get" itself accepts no arguments and requires no service,
			// so this exercises PersistentPreRun without building one.
			if err := cmd.Execute(); err != nil {
				t.Fatalf("Execute() error = %v", err)
			}

			if outputSettings.Format != tt.format {
				t.Errorf(
					"Format = %q, want %q",
					outputSettings.Format,
					tt.format,
				)
			}

			if outputSettings.IsQuiet != tt.quiet {
				t.Errorf(
					"Quiet = %v, want %v",
					outputSettings.IsQuiet,
					tt.quiet,
				)
			}

			if outputSettings.IsVerbose != tt.verbose {
				t.Errorf(
					"Verbose = %v, want %v",
					outputSettings.IsVerbose,
					tt.verbose,
				)
			}

			if outputSettings.IsNoHeaders != tt.noHeaders {
				t.Errorf(
					"NoHeaders = %v, want %v",
					outputSettings.IsNoHeaders,
					tt.noHeaders,
				)
			}
		})
	}
}

func TestRootCommandArgumentValidation(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{
			name: "validate environment missing spec",
			args: []string{
				"validate",
				"environment",
			},
		},
		{
			name: "create test missing spec",
			args: []string{
				"create",
				"test",
			},
		},
		{
			name: "get environment missing ref",
			args: []string{
				"get",
				"environment",
			},
		},
		{
			name: "delete testsuite missing ref",
			args: []string{
				"delete",
				"testsuite",
			},
		},
		{
			name: "run test missing environment ref",
			args: []string{
				"run",
				"test",
				"test-ref",
			},
		},
		{
			name: "internal runner rejects arguments",
			args: []string{
				"internal-runner",
				"unexpected",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetRootSettings()

			streams := newTestStreams()
			cmd := newRootCommand(
				streams,
				providers{},
			)

			cmd.SetArgs(tt.args)

			if err := cmd.Execute(); err == nil {
				t.Fatal("Execute() error = nil, want non-nil")
			}
		})
	}
}

func TestWithDefaultOutputFormat(t *testing.T) {
	tests := []struct {
		name     string
		settings output.Settings
		format   render.Format
		want     render.Format
	}{
		{
			name:   "empty",
			format: render.FormatYAML,
			want:   render.FormatYAML,
		},
		{
			name: "already set",
			settings: output.Settings{
				Format: render.FormatJSON,
			},
			format: render.FormatYAML,
			want:   render.FormatJSON,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := withDefaultOutputFormat(
				&tt.settings,
				tt.format,
			)

			if got.Format != tt.want {
				t.Errorf(
					"Format = %q, want %q",
					got.Format,
					tt.want,
				)
			}

			// The helper must not mutate the original settings.
			if tt.settings.Format != "" &&
				tt.settings.Format != got.Format {
				t.Errorf(
					"original Format = %q after call",
					tt.settings.Format,
				)
			}
		})
	}
}

func resetRootSettings() {
	outputFormat = ""
	isQuiet = false
	isVerbose = false
	isNoHeaders = false
	*outputSettings = output.Settings{}
}
