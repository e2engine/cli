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

func newDeleteCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete a resource",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newDeleteEnvironmentCommand(provider, outputSettings),
		newDeleteTestCommand(provider, outputSettings),
		newDeleteTestSuiteCommand(provider, outputSettings),
		newDeleteTestExecutionCommand(provider, outputSettings),
		newDeleteTestSuiteExecutionCommand(provider, outputSettings),
	)

	return cmd
}

func newDeleteEnvironmentCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment <ref>",
		Aliases: []string{"env"},
		Short:   "Delete an environment by a prefix of its id or name.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := environmentcommand.NewDeleteHandler(service)

					return handler.Run(
						ctx,
						environmentcommand.DeleteParams{
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

func newDeleteTestCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test <ref>",
		Aliases: []string{"t"},
		Short:   "Delete a test by a prefix of its id or name.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testcommand.NewDeleteHandler(service)

					return handler.Run(
						ctx,
						testcommand.DeleteParams{
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

func newDeleteTestSuiteCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuite <ref>",
		Aliases: []string{"ts"},
		Short:   "Delete a test suite by a prefix of its id or name.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuitecommand.NewDeleteHandler(service)

					return handler.Run(
						ctx,
						testsuitecommand.DeleteParams{
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

func newDeleteTestExecutionCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testexecution <id>",
		Aliases: []string{"te"},
		Short:   "Delete a test execution by a prefix of its id.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testexecutioncommand.NewDeleteHandler(service)

					return handler.Run(
						ctx,
						testexecutioncommand.DeleteParams{
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

func newDeleteTestSuiteExecutionCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuiteexecution <id>",
		Aliases: []string{"tse"},
		Short:   "Delete a test suite execution by a prefix of its id.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuiteexecutioncommand.NewDeleteHandler(service)

					return handler.Run(
						ctx,
						testsuiteexecutioncommand.DeleteParams{
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
