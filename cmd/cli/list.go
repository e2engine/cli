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

func newGetEnvironmentsCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	var (
		environmentsListLimit                           int
		environmentsOrderBy, environmentsOrderDirection string
	)

	cmd := &cobra.Command{
		Use:     "environments",
		Aliases: []string{"envs"},
		Short:   "Display environments list",
		Args:    cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := environmentcommand.NewListHandler(service)

					return handler.Run(
						ctx,
						&environmentcommand.ListParams{
							Limit:          environmentsListLimit,
							OrderBy:        environmentsOrderBy,
							OrderDirection: environmentsOrderDirection,
							Output:         settings,
							Writer:         cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	cmd.Flags().IntVar(
		&environmentsListLimit,
		"limit",
		0,
		"Maximum number of environments to display; values above the configured maximum are capped",
	)

	cmd.Flags().StringVar(
		&environmentsOrderBy,
		"order-by",
		"",
		"Order by field. Accepted values are: id, name, version, createdAt, updatedAt",
	)

	cmd.Flags().StringVar(
		&environmentsOrderDirection,
		"order-direction",
		"",
		"Order direction. Accepted values are: asc, desc",
	)

	return cmd
}

func newGetTestsCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	var (
		testsListLimit                    int
		testsOrderBy, testsOrderDirection string
	)

	cmd := &cobra.Command{
		Use:   "tests",
		Short: "Display tests list",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testcommand.NewListHandler(service)

					return handler.Run(
						ctx,
						&testcommand.ListParams{
							Limit:          testsListLimit,
							OrderBy:        testsOrderBy,
							OrderDirection: testsOrderDirection,
							Output:         settings,
							Writer:         cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	cmd.Flags().IntVar(
		&testsListLimit,
		"limit",
		0,
		"Maximum number of tests to display; values above the configured maximum are capped",
	)

	cmd.Flags().StringVar(
		&testsOrderBy,
		"order-by",
		"",
		"Order by field. Accepted values are: id, name, version, createdAt, updatedAt",
	)

	cmd.Flags().StringVar(
		&testsOrderDirection,
		"order-direction",
		"",
		"Order direction. Accepted values are: asc, desc",
	)

	return cmd
}

func newGetTestSuitesCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	var (
		testSuitesListLimit                         int
		testSuitesOrderBy, testSuitesOrderDirection string
	)

	cmd := &cobra.Command{
		Use:   "testsuites",
		Short: "Display test suites list",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuitecommand.NewListHandler(service)

					return handler.Run(
						ctx,
						&testsuitecommand.ListParams{
							Limit:          testSuitesListLimit,
							OrderBy:        testSuitesOrderBy,
							OrderDirection: testSuitesOrderDirection,
							Output:         settings,
							Writer:         cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	cmd.Flags().IntVar(
		&testSuitesListLimit,
		"limit",
		0,
		"Maximum number of test suites to display; values above the configured maximum are capped",
	)

	cmd.Flags().StringVar(
		&testSuitesOrderBy,
		"order-by",
		"",
		"Order by field. Accepted values are: id, name, version, createdAt, updatedAt",
	)

	cmd.Flags().StringVar(
		&testSuitesOrderDirection,
		"order-direction",
		"",
		"Order direction. Accepted values are: asc, desc",
	)

	return cmd
}

func newGetTestExecutionsCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	var (
		testExecutionsListLimit                             int
		testExecutionsOrderBy, testExecutionsOrderDirection string
	)

	cmd := &cobra.Command{
		Use:   "testexecutions",
		Short: "Display test executions list",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testexecutioncommand.NewListHandler(service)

					return handler.Run(
						ctx,
						&testexecutioncommand.ListParams{
							Limit:          testExecutionsListLimit,
							OrderBy:        testExecutionsOrderBy,
							OrderDirection: testExecutionsOrderDirection,
							Output:         settings,
							Writer:         cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	cmd.Flags().IntVar(
		&testExecutionsListLimit,
		"limit",
		0,
		"Maximum number of test executions to display; values above the configured maximum are capped",
	)

	cmd.Flags().StringVar(
		&testExecutionsOrderBy,
		"order-by",
		"",
		"Order by field. Accepted values are: id, status, startedAt, finishedAt",
	)

	cmd.Flags().StringVar(
		&testExecutionsOrderDirection,
		"order-direction",
		"",
		"Order direction. Accepted values are: asc, desc",
	)

	return cmd
}

func newGetTestSuiteExecutionsCommand(provider serviceProvider, outputSettings *output.Settings) *cobra.Command {
	var (
		testSuiteExecutionsListLimit                                  int
		testSuiteExecutionsOrderBy, testSuiteExecutionsOrderDirection string
	)

	cmd := &cobra.Command{
		Use:   "testsuiteexecutions",
		Short: "Display test suite executions list",
		Args:  cobra.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {
			settings := withDefaultOutputFormat(outputSettings, "table")

			return runWithService(
				cmd,
				args,
				provider,
				func(ctx context.Context, service *service.Service) error {
					handler := testsuiteexecutioncommand.NewListHandler(service)

					return handler.Run(
						ctx,
						&testsuiteexecutioncommand.ListParams{
							Limit:          testSuiteExecutionsListLimit,
							OrderBy:        testSuiteExecutionsOrderBy,
							OrderDirection: testSuiteExecutionsOrderDirection,
							Output:         settings,
							Writer:         cmd.OutOrStdout(),
						},
					)
				},
			)
		},
	}

	cmd.Flags().IntVar(
		&testSuiteExecutionsListLimit,
		"limit",
		0,
		"Maximum number of test suite executions to display; values above the configured maximum are capped",
	)

	cmd.Flags().StringVar(
		&testSuiteExecutionsOrderBy,
		"order-by",
		"",
		"Order by field. Accepted values are: id, status, startedAt, finishedAt",
	)

	cmd.Flags().StringVar(
		&testSuiteExecutionsOrderDirection,
		"order-direction",
		"",
		"Order direction. Accepted values are: asc, desc",
	)

	return cmd
}
