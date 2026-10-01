package config_test

import (
	"context"
	"testing"

	"github.com/ygrebnov/log"

	"github.com/e2engine/cli/internal/config"
)

func TestGetDefaultLoggerConfig(t *testing.T) {
	cfg, err := config.GetDefaultLoggerConfig()
	if err != nil {
		t.Fatalf("GetDefaultLoggerConfig() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("GetDefaultLoggerConfig() = nil")
	}

	if cfg.AppName != "e2engine" {
		t.Errorf(
			"AppName = %q, want %q",
			cfg.AppName,
			"e2engine",
		)
	}

	if len(cfg.Sinks) != 2 {
		t.Fatalf(
			"len(Sinks) = %d, want 2",
			len(cfg.Sinks),
		)
	}

	if cfg.Sinks[0].Kind != log.KindStdErr {
		t.Errorf(
			"Sinks[0].Kind = %q, want %q",
			cfg.Sinks[0].Kind,
			log.KindStdErr,
		)
	}

	if cfg.Sinks[1].Kind != log.KindFile {
		t.Errorf(
			"Sinks[1].Kind = %q, want %q",
			cfg.Sinks[1].Kind,
			log.KindFile,
		)
	}

	if cfg.Sinks[1].Level != log.LevelDebug {
		t.Errorf(
			"Sinks[1].Level = %q, want %q",
			cfg.Sinks[1].Level,
			log.LevelDebug,
		)
	}
}

func TestConfig_Validate_NilOptionalConfigs(t *testing.T) {
	cfg := &config.Config{}

	err := cfg.Validate(context.Background())
	if err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestConfig_Validate_DefaultLoggerConfig(t *testing.T) {
	loggerCfg, err := config.GetDefaultLoggerConfig()
	if err != nil {
		t.Fatalf("GetDefaultLoggerConfig() error = %v", err)
	}

	cfg := &config.Config{
		Logger: loggerCfg,
	}

	if err := cfg.Validate(context.Background()); err != nil {
		t.Fatalf("Validate() error = %v", err)
	}
}

func TestLoad_Defaults(t *testing.T) {
	cfg, err := config.Load(
		context.Background(),
		nil,
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg == nil {
		t.Fatal("Load() = nil")
	}

	if cfg.Output == nil {
		t.Fatal("Output = nil")
	}

	if cfg.Output.ListMaxSize != 50 {
		t.Errorf(
			"Output.ListMaxSize = %d, want 50",
			cfg.Output.ListMaxSize,
		)
	}

	if cfg.DB == nil {
		t.Error("DB = nil, want default config")
	}

	if cfg.Runtime == nil {
		t.Error("Runtime = nil, want default config")
	}

	if cfg.Runner == nil {
		t.Error("Runner = nil, want default config")
	}

	if cfg.Scheduler == nil {
		t.Error("Scheduler = nil, want default config")
	}
}

func TestLoad_OutputListMaxSizeFromEnvironment(t *testing.T) {
	t.Setenv(
		"E2ENGINE_OUTPUT_LIST_MAX_SIZE",
		"25",
	)

	cfg, err := config.Load(
		context.Background(),
		nil,
	)
	if err != nil {
		t.Fatalf("Load() error = %v", err)
	}

	if cfg.Output == nil {
		t.Fatal("Output = nil")
	}

	if cfg.Output.ListMaxSize != 25 {
		t.Errorf(
			"Output.ListMaxSize = %d, want 25",
			cfg.Output.ListMaxSize,
		)
	}
}

func TestLoad_InvalidOutputListMaxSize(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "below minimum",
			value: "0",
		},
		{
			name:  "above maximum",
			value: "51",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv(
				"E2ENGINE_OUTPUT_LIST_MAX_SIZE",
				tt.value,
			)

			cfg, err := config.Load(
				context.Background(),
				nil,
			)

			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
			if cfg != nil {
				t.Errorf(
					"Load() = %#v, want nil",
					cfg,
				)
			}
		})
	}
}
