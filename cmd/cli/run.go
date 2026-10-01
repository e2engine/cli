package main

import (
	"context"

	"github.com/spf13/cobra"

	testcommand "github.com/e2engine/cli/internal/command/test"
	testsuitecommand "github.com/e2engine/cli/internal/command/testsuite"
	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/service"
)

func newRunCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "run",
		Short: "Run a test or test suite against an environment.",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newRunTestCommand(provider, outputSettings),
		newRunTestSuiteCommand(provider, outputSettings),
	)

	return cmd
}

func newRunTestCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test <test-ref> <env-ref>",
		Aliases: []string{"t"},
		Short:   "Run a test.",
		Args:    cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testcommand.NewRunHandler(service)

					return handler.Run(
						ctx,
						testcommand.RunParams{
							EnvironmentReference: args[1],
							TestReference:        args[0],
							Output:               settings,
							Writer:               cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}

func newRunTestSuiteCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuite <testsuite-ref> <env-ref>",
		Aliases: []string{"ts"},
		Short:   "Run a test suite.",
		Args:    cobra.ExactArgs(2),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuitecommand.NewRunHandler(service)

					return handler.Run(
						ctx,
						testsuitecommand.RunParams{
							EnvironmentReference: args[1],
							TestSuiteReference:   args[0],
							Output:               settings,
							Writer:               cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}
