package fs

import (
	"encoding/json"
	"path/filepath"
	"strings"

	"github.com/e2engine/core/model"
	"github.com/e2engine/core/pkg/errors"
	"github.com/e2engine/core/pkg/keys"
	"github.com/e2engine/core/util"
	"github.com/ygrebnov/errorc"
	"gopkg.in/yaml.v3"
)

func Load[E model.Entity](path string) (*E, error) {
	data, err := util.ReadFile(path)
	if err != nil {
		return nil, err
	}

	e := new(E)
	ext := strings.ToLower(filepath.Ext(path))

	switch ext {
	case ".yaml", ".yml":
		if err := yaml.Unmarshal(data, e); err != nil {
			return nil, err
		}

	case ".json":
		if err := json.Unmarshal(data, e); err != nil {
			return nil, err
		}

	default:
		return nil, errorc.With(
			errors.ErrUnsupportedFileType,
			errorc.String(keys.FilePath, path),
			errorc.String(keys.FileFormat, ext),
		)
	}

	return e, nil
}
