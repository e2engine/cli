package main

import (
	"context"
	"time"

	"github.com/spf13/cobra"

	testexecutioncommand "github.com/e2engine/cli/internal/command/testexecution"
	testsuiteexecutioncommand "github.com/e2engine/cli/internal/command/testsuiteexecution"
	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/service"
)

var executionCheckTimeout time.Duration

func newCheckCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Wait for an execution to complete and check its result.",
		Args:  cobra.NoArgs,
	}

	cmd.PersistentFlags().DurationVarP(
		&executionCheckTimeout,
		"timeout",
		"t",
		0,
		"Execution check timeout (default: from global config)",
	)

	cmd.AddCommand(
		newCheckTestExecutionCommand(provider, outputSettings),
		newCheckTestSuiteExecutionCommand(provider, outputSettings),
	)

	return cmd
}

func newCheckTestExecutionCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testexecution <ref>",
		Aliases: []string{"te"},
		Short:   "Wait for a test execution to complete and check its result.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					timeout := service.GetConfig().Runtime.ExecutionCheckTimeout

					if executionCheckTimeout != 0 {
						timeout = executionCheckTimeout
					}

					handler := testexecutioncommand.NewCheckHandler(service)

					return handler.Run(
						ctx,
						testexecutioncommand.CheckParams{
							Reference: args[0],
							Timeout:   timeout,
							Output:    settings,
							Writer:    cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}

func newCheckTestSuiteExecutionCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuiteexecution <ref>",
		Aliases: []string{"tse"},
		Short:   "Wait for a testsuite execution to complete and check its result.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					timeout := service.GetConfig().Runtime.ExecutionCheckTimeout

					if executionCheckTimeout != 0 {
						timeout = executionCheckTimeout
					}

					handler := testsuiteexecutioncommand.NewCheckHandler(service)

					return handler.Run(
						ctx,
						testsuiteexecutioncommand.CheckParams{
							Reference: args[0],
							Timeout:   timeout,
							Output:    settings,
							Writer:    cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}
