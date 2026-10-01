package service

import (
	"context"
	"path/filepath"
	"testing"

	"github.com/e2engine/core/execute"
	dbconfig "github.com/e2engine/repository/sqlite/config"

	"github.com/e2engine/cli/internal/config"
	"github.com/e2engine/cli/internal/util/paths"
)

func TestWithRunner(t *testing.T) {
	s := &settings{}

	withRunner()(s)

	if !s.buildRunner {
		t.Fatal("withRunner() did not enable buildRunner")
	}
}

func TestGetDataDir_FromEnvironment(t *testing.T) {
	want := filepath.Join(
		t.TempDir(),
		"e2engine-data",
	)

	t.Setenv("E2ENGINE_DATA_DIR", want)

	got, err := getDataDir()
	if err != nil {
		t.Fatalf("getDataDir() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"getDataDir() = %q, want %q",
			got,
			want,
		)
	}
}

func TestGetDataDir_TrimsEnvironment(t *testing.T) {
	want := filepath.Join(
		t.TempDir(),
		"e2engine-data",
	)

	t.Setenv(
		"E2ENGINE_DATA_DIR",
		"  "+want+"  ",
	)

	got, err := getDataDir()
	if err != nil {
		t.Fatalf("getDataDir() error = %v", err)
	}

	if got != want {
		t.Fatalf(
			"getDataDir() = %q, want %q",
			got,
			want,
		)
	}
}

func TestGetDataDir_Default(t *testing.T) {
	tests := []struct {
		name  string
		value string
	}{
		{
			name:  "empty",
			value: "",
		},
		{
			name:  "whitespace",
			value: "   ",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Setenv(
				"E2ENGINE_DATA_DIR",
				test.value,
			)

			want, err := paths.GetDataDir(appName)
			if err != nil {
				t.Fatalf(
					"paths.GetDataDir() error = %v",
					err,
				)
			}

			got, err := getDataDir()
			if err != nil {
				t.Fatalf(
					"getDataDir() error = %v",
					err,
				)
			}

			if got != want {
				t.Fatalf(
					"getDataDir() = %q, want %q",
					got,
					want,
				)
			}
		})
	}
}

func TestBuildSchedulerService(t *testing.T) {
	jobs := make(chan execute.TestJob, 1)

	scheduler, err := buildSchedulerService(jobs)
	if err != nil {
		t.Fatalf(
			"buildSchedulerService() error = %v",
			err,
		)
	}

	if scheduler == nil {
		t.Fatal("buildSchedulerService() = nil")
	}
}

func TestBuildLogger_DefaultConfig(t *testing.T) {
	logger, err := BuildLogger(nil)
	if err != nil {
		t.Fatalf("BuildLogger() error = %v", err)
	}
	if logger == nil {
		t.Fatal("BuildLogger() = nil")
	}

	t.Cleanup(func() {
		if err := logger.Close(); err != nil {
			t.Errorf("logger.Close() error = %v", err)
		}
	})
}

func TestBuildRepositories(t *testing.T) {
	dataDir := t.TempDir()
	t.Setenv("E2ENGINE_DATA_DIR", dataDir)

	cfg := &config.Config{
		DB: &dbconfig.Config{
			Name:         "e2engine.sqlite",
			Driver:       "sqlite",
			MaxOpenConns: 1,
			MaxIdleConns: 1,
		},
	}

	repositories, closer, err := BuildRepositories(
		context.Background(),
		cfg,
	)
	if err != nil {
		t.Fatalf(
			"BuildRepositories() error = %v",
			err,
		)
	}

	if closer == nil {
		t.Fatal("BuildRepositories() closer = nil")
	}

	t.Cleanup(func() {
		if err := closer.Close(); err != nil {
			t.Errorf("closer.Close() error = %v", err)
		}
	})

	// Check whichever repository fields are guaranteed by coreapi.Repositories.
	if repositories.Environment == nil {
		t.Error("repositories.Environment = nil")
	}
	if repositories.Test == nil {
		t.Error("repositories.Test = nil")
	}
	if repositories.TestExecution == nil {
		t.Error("repositories.TestExecution = nil")
	}
}
