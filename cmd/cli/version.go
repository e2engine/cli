package main

import (
	"github.com/spf13/cobra"

	versioncommand "github.com/e2engine/cli/internal/command/version"
	"github.com/e2engine/cli/internal/output"
)

func newVersionCommand(outputSettings *output.Settings) *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Show e2engine tool version",
		RunE: func(cmd *cobra.Command, _ []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return versioncommand.NewHandler().Run(cmd.Context(), versioncommand.Params{
				Writer: cmd.OutOrStdout(),
				Output: settings,
			})
		},
	}
}
