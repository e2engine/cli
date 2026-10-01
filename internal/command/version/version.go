package version

import (
	"context"
	"fmt"
	"io"
	"runtime"
	"time"

	"github.com/e2engine/core/model"

	"github.com/e2engine/cli/internal/output"
	"github.com/e2engine/cli/internal/output/view"
)

var (
	// Injected by the build system.
	version      string
	gitCommit    string
	gitTreeState string
	buildTime    = "1970-01-01T00:00:00Z"
)

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

type Params struct {
	Output output.Settings
	Writer io.Writer
}

func (h *Handler) Run(_ context.Context, params Params) error {
	return view.RenderObject(params.Writer, buildVersion(), params.Output)
}

func buildVersion() model.Version {
	return adjustVersion(model.Version{
		Version:      version,
		BuildTime:    parseBuildTime(),
		GitCommit:    gitCommit,
		GitTreeDirty: gitTreeState != "clean",
		GoVersion:    runtime.Version(),
		Compiler:     runtime.Compiler,
		Platform:     fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH),
	})
}

func parseBuildTime() time.Time {
	t, err := time.Parse(time.RFC3339Nano, buildTime)
	if err != nil {
		return time.Time{}
	}

	return t
}

func adjustVersion(v model.Version) model.Version {
	if v.Version != "" {
		return v
	}

	v.Version = "devel"

	switch {
	case len(v.GitCommit) >= 7:
		v.Version += "+" + v.GitCommit[:7]

	case v.GitCommit != "":
		v.Version += "+" + v.GitCommit

	default:
		v.Version += "+unknown"
	}

	if v.GitTreeDirty {
		v.Version += ".dirty"
	}

	return v
}
