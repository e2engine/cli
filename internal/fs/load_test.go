package fs_test

import (
	"errors"
	"os"
	"path/filepath"
	"testing"

	"github.com/e2engine/core/model"
	coreerrors "github.com/e2engine/core/pkg/errors"

	"github.com/e2engine/cli/internal/fs"
)

func TestLoad(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
		want     *model.Environment
	}{
		{
			name:     "yaml",
			filename: "environment.yaml",
			content: `
kind: Environment
name: local
version: 1.0.0
description: Local environment
`,
			want: &model.Environment{
				Kind:        model.ResourceKindEnvironment,
				Name:        "local",
				Version:     "1.0.0",
				Description: "Local environment",
			},
		},
		{
			name:     "yml",
			filename: "environment.yml",
			content: `
kind: Environment
name: local
version: 1.0.0
`,
			want: &model.Environment{
				Kind:    model.ResourceKindEnvironment,
				Name:    "local",
				Version: "1.0.0",
			},
		},
		{
			name:     "json",
			filename: "environment.json",
			content: `{
				"kind": "Environment",
				"name": "local",
				"version": "1.0.0",
				"description": "Local environment"
			}`,
			want: &model.Environment{
				Kind:        model.ResourceKindEnvironment,
				Name:        "local",
				Version:     "1.0.0",
				Description: "Local environment",
			},
		},
		{
			name:     "uppercase extension",
			filename: "environment.YAML",
			content: `
kind: Environment
name: local
version: 1.0.0
`,
			want: &model.Environment{
				Kind:    model.ResourceKindEnvironment,
				Name:    "local",
				Version: "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeFile(t, tt.filename, tt.content)

			got, err := fs.Load[model.Environment](path)
			if err != nil {
				t.Fatalf("Load() error = %v", err)
			}

			if got.Kind != tt.want.Kind {
				t.Errorf("Kind = %q, want %q", got.Kind, tt.want.Kind)
			}
			if got.Name != tt.want.Name {
				t.Errorf("Name = %q, want %q", got.Name, tt.want.Name)
			}
			if got.Version != tt.want.Version {
				t.Errorf("Version = %q, want %q", got.Version, tt.want.Version)
			}
			if got.Description != tt.want.Description {
				t.Errorf(
					"Description = %q, want %q",
					got.Description,
					tt.want.Description,
				)
			}
		})
	}
}

func TestLoad_InvalidContent(t *testing.T) {
	tests := []struct {
		name     string
		filename string
		content  string
	}{
		{
			name:     "invalid yaml",
			filename: "environment.yaml",
			content:  "name: [",
		},
		{
			name:     "invalid json",
			filename: "environment.json",
			content:  `{"name":`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			path := writeFile(t, tt.filename, tt.content)

			got, err := fs.Load[model.Environment](path)
			if err == nil {
				t.Fatal("Load() error = nil, want error")
			}
			if got != nil {
				t.Errorf("Load() = %#v, want nil", got)
			}
		})
	}
}

func TestLoad_UnsupportedFileType(t *testing.T) {
	path := writeFile(
		t,
		"environment.txt",
		"kind: Environment\nname: local\n",
	)

	got, err := fs.Load[model.Environment](path)

	if !errors.Is(err, coreerrors.ErrUnsupportedFileType) {
		t.Fatalf(
			"Load() error = %v, want %v",
			err,
			coreerrors.ErrUnsupportedFileType,
		)
	}
	if got != nil {
		t.Errorf("Load() = %#v, want nil", got)
	}
}

func TestLoad_FileDoesNotExist(t *testing.T) {
	path := filepath.Join(t.TempDir(), "missing.yaml")

	got, err := fs.Load[model.Environment](path)

	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if got != nil {
		t.Errorf("Load() = %#v, want nil", got)
	}
}

func TestLoad_Directory(t *testing.T) {
	path := t.TempDir()

	got, err := fs.Load[model.Environment](path)

	if err == nil {
		t.Fatal("Load() error = nil, want error")
	}
	if got != nil {
		t.Errorf("Load() = %#v, want nil", got)
	}
}

func writeFile(
	t *testing.T,
	name string,
	content string,
) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), name)

	if err := os.WriteFile(path, []byte(content), 0o600); err != nil {
		t.Fatalf("os.WriteFile() error = %v", err)
	}

	return path
}
