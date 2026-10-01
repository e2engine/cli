package version

import (
	"testing"
	"time"

	"github.com/e2engine/core/model"
)

func TestAdjustVersion(t *testing.T) {
	tests := []struct {
		name string
		in   model.Version
		want string
	}{
		{
			name: "release version unchanged",
			in: model.Version{
				Version:      "1.2.3",
				GitCommit:    "abcdef123456",
				GitTreeDirty: true,
			},
			want: "1.2.3",
		},
		{
			name: "development with long commit",
			in: model.Version{
				GitCommit: "abcdef123456",
			},
			want: "devel+abcdef1",
		},
		{
			name: "development with short commit",
			in: model.Version{
				GitCommit: "abc123",
			},
			want: "devel+abc123",
		},
		{
			name: "development without commit",
			in:   model.Version{},
			want: "devel+unknown",
		},
		{
			name: "dirty development build",
			in: model.Version{
				GitCommit:    "abcdef123456",
				GitTreeDirty: true,
			},
			want: "devel+abcdef1.dirty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := adjustVersion(tt.in)

			if got.Version != tt.want {
				t.Errorf(
					"adjustVersion().Version = %q, want %q",
					got.Version,
					tt.want,
				)
			}
		})
	}
}

func TestParseBuildTime(t *testing.T) {
	original := buildTime
	t.Cleanup(func() {
		buildTime = original
	})

	buildTime = "2026-09-22T08:15:30.123456789Z"

	got := parseBuildTime()

	want := time.Date(
		2026, time.September, 22,
		8, 15, 30, 123456789,
		time.UTC,
	)

	if !got.Equal(want) {
		t.Errorf("parseBuildTime() = %v, want %v", got, want)
	}
}

func TestParseBuildTime_Invalid(t *testing.T) {
	original := buildTime
	t.Cleanup(func() {
		buildTime = original
	})

	buildTime = "invalid"

	got := parseBuildTime()

	if !got.IsZero() {
		t.Errorf("parseBuildTime() = %v, want zero time", got)
	}
}
