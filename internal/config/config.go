package config

import (
	"context"
	"strings"

	runtimeconfig "github.com/e2engine/core/execute/runtime/config"
	coremodel "github.com/e2engine/core/model"
	"github.com/e2engine/core/pkg/keys"
	"github.com/e2engine/core/pkg/log"
	"github.com/e2engine/core/scheduler"
	"github.com/e2engine/core/transport"
	dbconfig "github.com/e2engine/repository/sqlite/config"
	runnerconfig "github.com/e2engine/runner-local/pkg/config"
	configlib "github.com/ygrebnov/config"
	"github.com/ygrebnov/errorc"
	loglib "github.com/ygrebnov/log"
	logconfig "github.com/ygrebnov/log/pkg/config"

	"github.com/e2engine/cli/pkg/errors"
)

const appName = "e2engine"

type Config struct {
	Output *OutputConfig `json:"output" yaml:"output" default:"dive" validateElem:"omitempty"`
	// Logger is optional. If not provided, a default logger config will be used.
	// Default logger config is provided not via tags, but via GetDefaultLoggerConfig() function.
	Logger    *loglib.Config        `json:"logger,omitempty" yaml:"logger,omitempty" validateElem:"omitempty"`
	DB        *dbconfig.Config      `json:"db" yaml:"db" default:"dive" validateElem:"omitempty"`
	Runtime   *runtimeconfig.Config `json:"runtime" yaml:"runtime" default:"dive" validateElem:"omitempty"`
	Scheduler *scheduler.Config     `json:"scheduler" yaml:"scheduler" default:"dive" validateElem:"omitempty"`
	Runner    *runnerconfig.Config  `json:"runner" yaml:"runner" default:"dive" validateElem:"omitempty"`
	Transport *transport.Config     `json:"transport,omitempty" yaml:"transport,omitempty" validateElem:"omitempty"`
}

func (c *Config) Validate(ctx context.Context) error {
	if c.Logger != nil {
		if err := c.Logger.Validate(ctx); err != nil {
			return err
		}
	}

	if c.Transport != nil {
		if err := c.Transport.Validate(); err != nil {
			return err
		}
	}

	return nil
}

type OutputConfig struct {
	ListMaxSize int `json:"list_max_size" yaml:"list_max_size" env:"LIST_MAX_SIZE" default:"50" validate:"min(1),max(50)"`
}

func GetDefaultLoggerConfig() (*loglib.Config, error) {
	cfg := &loglib.Config{
		AppName: appName,
		Sinks: []loglib.SinkConfig{
			{Kind: loglib.KindStdErr},
			{Kind: loglib.KindFile, Level: loglib.LevelDebug},
		},
	}

	if err := cfg.ApplyDefaults(); err != nil {
		return nil, errorc.With(
			errors.ErrCannotConfigureLogger,
			errorc.Error(keys.Cause, err),
		)
	}

	return cfg, nil
}

func Load(ctx context.Context, logger log.Logger) (*Config, error) {
	var cfg Config
	if err := configlib.Load(
		ctx,
		&cfg,
		configlib.WithEnvPrefix(strings.ToUpper(appName)),
		configlib.WithAppName(appName),
		configlib.WithValidationRules(
			append(
				logconfig.GetValidationRules(),
				coremodel.GetValidationRules()...,
			)...,
		),
		configlib.WithLogger(logger),
	); err != nil {
		return nil, errorc.With(
			errors.ErrCannotLoadConfig,
			errorc.Error(keys.Cause, err),
		)
	}

	if err := cfg.Validate(ctx); err != nil {
		return nil, err
	}

	return &cfg, nil
}
