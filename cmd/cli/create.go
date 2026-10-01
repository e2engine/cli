package main

import (
	"context"

	"github.com/spf13/cobra"

	environmentcommand "github.com/e2engine/cli/internal/command/environment"
	testcommand "github.com/e2engine/cli/internal/command/test"
	testsuitecommand "github.com/e2engine/cli/internal/command/testsuite"
	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/service"
)

func newCreateCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new resource",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newCreateEnvironmentCommand(provider, outputSettings),
		newCreateTestCommand(provider, outputSettings),
		newCreateTestSuiteCommand(provider, outputSettings),
	)

	return cmd
}

func newCreateEnvironmentCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment <spec-path>",
		Aliases: []string{"env"},
		Short:   "Create a new environment from a given specification.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := environmentcommand.NewCreateHandler(service)

					return handler.Run(
						ctx,
						environmentcommand.CreateParams{
							SpecPath: args[0],
							Output:   settings,
							Writer:   cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}

func newCreateTestCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test <spec-path>",
		Aliases: []string{"t"},
		Short:   "Create a new test from a given specification.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testcommand.NewCreateHandler(service)

					return handler.Run(
						ctx,
						testcommand.CreateParams{
							SpecPath: args[0],
							Output:   settings,
							Writer:   cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}

func newCreateTestSuiteCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuite <spec-path>",
		Aliases: []string{"ts"},
		Short:   "Create a new test suite from a given specification.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuitecommand.NewCreateHandler(service)

					return handler.Run(
						ctx,
						testsuitecommand.CreateParams{
							SpecPath: args[0],
							Output:   settings,
							Writer:   cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	return cmd
}
