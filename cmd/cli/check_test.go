package main

import (
	"testing"
	"time"

	"github.com/e2engine/cli/internal/output"
)

func TestNewCheckCommand(t *testing.T) {
	executionCheckTimeout = 0
	t.Cleanup(func() {
		executionCheckTimeout = 0
	})

	settings := &output.Settings{}

	cmd := newCheckCommand(nil, settings)

	if cmd.Use != "check" {
		t.Errorf("Use = %q, want %q", cmd.Use, "check")
	}

	if cmd.Short != "Wait for an execution to complete and check its result." {
		t.Errorf(
			"Short = %q, want %q",
			cmd.Short,
			"Wait for an execution to complete and check its result.",
		)
	}

	if err := cmd.Args(cmd, nil); err != nil {
		t.Errorf("Args() error = %v", err)
	}

	flag := cmd.PersistentFlags().Lookup("timeout")
	if flag == nil {
		t.Fatal("timeout flag not found")
	}

	if flag.Shorthand != "t" {
		t.Errorf("timeout shorthand = %q, want %q", flag.Shorthand, "t")
	}

	if flag.DefValue != "0s" {
		t.Errorf("timeout default = %q, want %q", flag.DefValue, "0s")
	}
}

func TestNewCheckCommand_RejectsArguments(t *testing.T) {
	executionCheckTimeout = 0
	t.Cleanup(func() {
		executionCheckTimeout = 0
	})

	cmd := newCheckCommand(nil, &output.Settings{})

	err := cmd.Args(cmd, []string{"unexpected"})
	if err == nil {
		t.Fatal("Args() error = nil, want error")
	}
}

func TestNewCheckCommand_TimeoutFlag(t *testing.T) {
	executionCheckTimeout = 0
	t.Cleanup(func() {
		executionCheckTimeout = 0
	})

	cmd := newCheckCommand(nil, &output.Settings{})

	err := cmd.PersistentFlags().Set("timeout", "15s")
	if err != nil {
		t.Fatalf("Set(timeout) error = %v", err)
	}

	if executionCheckTimeout != 15*time.Second {
		t.Errorf(
			"executionCheckTimeout = %v, want %v",
			executionCheckTimeout,
			15*time.Second,
		)
	}
}

func TestNewCheckCommand_TimeoutFlag_Shorthand(t *testing.T) {
	executionCheckTimeout = 0
	t.Cleanup(func() {
		executionCheckTimeout = 0
	})

	cmd := newCheckCommand(nil, &output.Settings{})

	err := cmd.PersistentFlags().Set("timeout", "20s")
	if err != nil {
		t.Fatalf("Set(t) error = %v", err)
	}

	if executionCheckTimeout != 20*time.Second {
		t.Errorf(
			"executionCheckTimeout = %v, want %v",
			executionCheckTimeout,
			20*time.Second,
		)
	}
}

func TestNewCheckTestExecutionCommand(t *testing.T) {
	cmd := newCheckTestExecutionCommand(
		nil,
		&output.Settings{},
	)

	if cmd.Use != "testexecution <ref>" {
		t.Errorf(
			"Use = %q, want %q",
			cmd.Use,
			"testexecution <ref>",
		)
	}

	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "te" {
		t.Errorf(
			"Aliases = %v, want [te]",
			cmd.Aliases,
		)
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			wantErr: true,
		},
		{
			name: "reference",
			args: []string{
				"execution-ref",
			},
		},
		{
			name: "too many arguments",
			args: []string{
				"execution-ref",
				"unexpected",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Args(cmd, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Args() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("Args() error = %v", err)
			}
		})
	}
}

func TestNewCheckTestSuiteExecutionCommand(t *testing.T) {
	cmd := newCheckTestSuiteExecutionCommand(
		nil,
		&output.Settings{},
	)

	if cmd.Use != "testsuiteexecution <ref>" {
		t.Errorf(
			"Use = %q, want %q",
			cmd.Use,
			"testsuiteexecution <ref>",
		)
	}

	if len(cmd.Aliases) != 1 || cmd.Aliases[0] != "tse" {
		t.Errorf(
			"Aliases = %v, want [tse]",
			cmd.Aliases,
		)
	}

	tests := []struct {
		name    string
		args    []string
		wantErr bool
	}{
		{
			name:    "no arguments",
			wantErr: true,
		},
		{
			name: "reference",
			args: []string{
				"execution-ref",
			},
		},
		{
			name: "too many arguments",
			args: []string{
				"execution-ref",
				"unexpected",
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := cmd.Args(cmd, tt.args)

			if tt.wantErr {
				if err == nil {
					t.Fatal("Args() error = nil, want error")
				}

				return
			}

			if err != nil {
				t.Fatalf("Args() error = %v", err)
			}
		})
	}
}
