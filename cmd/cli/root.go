package main

import (
	"context"
	nativeerrors "errors"
	"fmt"
	"strings"
	"time"

	"github.com/e2engine/core/pkg/keys"
	"github.com/e2engine/core/pkg/log"
	"github.com/spf13/cobra"
	"github.com/ygrebnov/errorc"
	"github.com/ygrebnov/model/validation"
	"github.com/ygrebnov/render"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/service"
	"github.com/e2engine/cli/pkg/errors"
)

const appName = "e2engine"

var (
	outputFormat                    string
	isQuiet, isVerbose, isNoHeaders bool
	outputSettings                  = new(output.Settings)
)

type providers struct {
	base       serviceProvider
	withRunner serviceProvider
}

func newRootCommand(streams ioStreams, providers providers) *cobra.Command {
	cmd := &cobra.Command{
		Use:           appName,
		Short:         "Manage and run end-to-end tests",
		SilenceErrors: true,
		SilenceUsage:  true,
	}

	cmd.PersistentFlags().StringVarP(
		&outputFormat,
		"output",
		"o",
		"",
		"Output format (available: table, json, yaml)",
	)

	cmd.PersistentFlags().BoolVarP(
		&isQuiet,
		"quiet",
		"q",
		false,
		"Suppress non-error output",
	)

	cmd.PersistentFlags().BoolVar(
		&isNoHeaders,
		"no-headers",
		false,
		"Omit headers in table output",
	)

	cmd.PersistentFlags().BoolVarP(
		&isVerbose,
		"verbose",
		"v",
		false,
		"Enable verbose output",
	)

	cmd.PersistentPreRun = func(_ *cobra.Command, _ []string) {
		*outputSettings = output.NewSettings(outputFormat, isQuiet, isVerbose, isNoHeaders)
	}

	cmd.SetIn(streams.In())
	cmd.SetOut(streams.Out())
	cmd.SetErr(streams.ErrOut())

	cmd.AddCommand(
		newValidateCommand(providers.base, outputSettings),
		newCreateCommand(providers.base, outputSettings),
		newGetCommand(providers.base, outputSettings),
		newDeleteCommand(providers.base, outputSettings),
		newCheckCommand(providers.base, outputSettings),
		newVersionCommand(outputSettings),
		newRunCommand(providers.withRunner, outputSettings),

		newInternalRunnerCommand(),
	)

	return cmd
}

func executeRootCommand(
	ctx context.Context,
	streams ioStreams,
	providers providers,
) int {
	cmd := newRootCommand(streams, providers)

	err := cmd.ExecuteContext(ctx)
	if err == nil {
		return 0
	}
	if isQuiet {
		return 1
	}
	if isVerbose {
		if vErr, ok := nativeerrors.AsType[*validation.Error](err); ok {
			settings := withDefaultOutputFormat(outputSettings, "yaml")

			renderErr := output.RenderValidationError(
				streams.ErrOut(),
				vErr,
				settings,
			)

			if renderErr != nil {
				_, _ = fmt.Fprintln(streams.ErrOut(), renderErr)
			}
			return 1
		}
	}
	_, _ = fmt.Fprintln(streams.ErrOut(), err)
	return 1
}

func runWithService(
	cmd *cobra.Command,
	args []string,
	provider serviceProvider,
	run func(context.Context, *service.Service) error,
) error {
	startTime := time.Now()
	ctx := cmd.Context()

	s, err := provider.Get(ctx)
	if err != nil {
		return errorc.With(errors.ErrCannotInitializeService, errorc.Error(keys.Cause, err))
	}

	logger := s.GetLogger()

	startRecord := log.NewRecord(
		startTime,
		log.LevelDebug,
		"start executing command",
		log.String("command", cmd.CommandPath()),
		log.String("args", strings.Join(args, " ")),
	)
	logger.LogRecord(startRecord)

	err = run(ctx, s)

	endTime := time.Now()
	endRecord := log.NewRecord(
		endTime,
		log.LevelDebug,
		"end executing command",
		log.String("command", cmd.CommandPath()),
		log.String("args", strings.Join(args, " ")),
		log.String("duration", endTime.Sub(startTime).String()),
		log.Err("err", err),
	)
	logger.LogRecord(endRecord)

	return err
}

func withDefaultOutputFormat(
	settings *output.Settings,
	format render.Format,
) output.Settings {
	result := *settings
	if result.Format == "" {
		result.Format = format
	}

	return result
}
