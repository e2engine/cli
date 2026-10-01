package main

import (
	"context"

	"github.com/spf13/cobra"

	environmentcommand "github.com/e2engine/cli/internal/command/environment"
	testcommand "github.com/e2engine/cli/internal/command/test"
	testexecutioncommand "github.com/e2engine/cli/internal/command/testexecution"
	testsuitecommand "github.com/e2engine/cli/internal/command/testsuite"
	testsuiteexecutioncommand "github.com/e2engine/cli/internal/command/testsuiteexecution"
	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/service"
)

func newGetCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "get",
		Short: "Display one or more resources",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newGetEnvironmentCommand(provider, outputSettings),
		newGetTestCommand(provider, outputSettings),
		newGetTestSuiteCommand(provider, outputSettings),
		newGetTestExecutionCommand(provider, outputSettings),
		newGetTestSuiteExecutionCommand(provider, outputSettings),

		newGetEnvironmentsCommand(provider, outputSettings),
		newGetTestsCommand(provider, outputSettings),
		newGetTestSuitesCommand(provider, outputSettings),
		newGetTestExecutionsCommand(provider, outputSettings),
		newGetTestSuiteExecutionsCommand(provider, outputSettings),
	)

	return cmd
}

func newGetEnvironmentCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment <ref>",
		Aliases: []string{"env"},
		Short:   "Display an environment by a prefix of its id or name.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := environmentcommand.NewGetHandler(service)

					return handler.Run(
						ctx,
						environmentcommand.GetParams{
							Reference: args[0],
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

func newGetTestCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test <ref>",
		Aliases: []string{"t"},
		Short:   "Display a test by a prefix of its id or name.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testcommand.NewGetHandler(service)

					return handler.Run(
						ctx,
						testcommand.GetParams{
							Reference: args[0],
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

func newGetTestSuiteCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuite <ref>",
		Aliases: []string{"ts"},
		Short:   "Display a testsuite by a prefix of its id or name.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuitecommand.NewGetHandler(service)

					return handler.Run(
						ctx,
						testsuitecommand.GetParams{
							Reference: args[0],
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

func newGetTestExecutionCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testexecution <ref>",
		Aliases: []string{"te"},
		Short:   "Display a test execution by a prefix of its id.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testexecutioncommand.NewGetHandler(service)

					return handler.Run(
						ctx,
						testexecutioncommand.GetParams{
							Reference: args[0],
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

func newGetTestSuiteExecutionCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuiteexecution <ref>",
		Aliases: []string{"tse"},
		Short:   "Display a testsuite execution by a prefix of its id.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuiteexecutioncommand.NewGetHandler(service)

					return handler.Run(
						ctx,
						testsuiteexecutioncommand.GetParams{
							Reference: args[0],
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
