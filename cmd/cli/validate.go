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

func newValidateCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "validate",
		Short: "Validate a resource specification",
		Args:  cobra.NoArgs,
	}

	cmd.AddCommand(
		newValidateEnvironmentCommand(provider, outputSettings),
		newValidateTestCommand(provider, outputSettings),
		newValidateTestSuiteCommand(provider, outputSettings),
	)

	return cmd
}

func newValidateEnvironmentCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "environment",
		Aliases: []string{"env"},
		Short:   "Validate an environment specification.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := environmentcommand.NewValidateHandler(service)

					return handler.Run(
						ctx,
						environmentcommand.ValidateParams{
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

func newValidateTestCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "test",
		Aliases: []string{"t"},
		Short:   "Validate a test specification.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testcommand.NewValidateHandler(service)

					return handler.Run(
						ctx,
						testcommand.ValidateParams{
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

func newValidateTestSuiteCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "testsuite",
		Aliases: []string{"ts"},
		Short:   "Validate a test suite specification.",
		Args:    cobra.ExactArgs(1),

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuitecommand.NewValidateHandler(service)

					return handler.Run(
						ctx,
						testsuitecommand.ValidateParams{
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
