package main

import (
	"github.com/spf13/cobra"

	"github.com/e2engine/cli/internal/command/internalrunner"
)

func newInternalRunnerCommand() *cobra.Command {
	return &cobra.Command{
		Use:    "internal-runner",
		Hidden: true,
		Args:   cobra.NoArgs,

		RunE: func(cmd *cobra.Command, _ []string) error {
			return internalrunner.Run(cmd.Context())
		},
	}
}
