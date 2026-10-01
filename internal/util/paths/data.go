package paths

import (
	"os"
	"path/filepath"
	"runtime"

	"github.com/e2engine/core/pkg/keys"
	"github.com/ygrebnov/errorc"

	"github.com/e2engine/cli/pkg/errors"
)

// GetDataDir returns the OS-appropriate per-user application data directory for the given appName.
// Linux/Unix:  $XDG_DATA_HOME or ~/.local/share/appName
// macOS:       ~/Library/Application Support/appName
// Windows:     %APPDATA%\appName or ~/AppData/Roaming/appName
func GetDataDir(appName string) (string, error) {
	switch runtime.GOOS {
	case "darwin":
		home, err := getUserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "Library", "Application Support", appName), nil
	case "windows":
		if appdata := os.Getenv("APPDATA"); appdata != "" {
			return filepath.Join(appdata, appName), nil
		}
		home, err := getUserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, "AppData", "Roaming", appName), nil
	default:
		// XDG on Unix
		if v := os.Getenv("XDG_DATA_HOME"); v != "" && filepath.IsAbs(v) {
			return filepath.Join(v, appName), nil
		}
		home, err := getUserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(home, ".local", "share", appName), nil
	}
}

func getUserHomeDir() (string, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		return "", errorc.With(
			errors.ErrCannotResolveUserHomeDir,
			errorc.Error(keys.Cause, err),
		)
	}
	return home, nil
}
