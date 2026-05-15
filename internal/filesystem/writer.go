package filesystem

import (
	"os"

	"github.com/athosbes/PeritiaGo/internal/config"
	"github.com/athosbes/PeritiaGo/internal/guard"
	"github.com/athosbes/PeritiaGo/internal/logger"
)

// WriteFile is a wrapper around os.WriteFile that checks permissions via the guard package.
// It ensures that in forensic mode, writes are only allowed within the output directory.
func WriteFile(filename string, data []byte, perm os.FileMode) error {
	safePath, err := guard.AllowWrite(filename)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Warn("write blocked by guard",
				"path", filename,
				"execution_id", config.CurrentConfig.ExecutionID,
			)
		}
		return err
	}

	if logger.Log != nil {
		logger.Log.Info("write allowed",
			"path", safePath,
			"execution_id", config.CurrentConfig.ExecutionID,
		)
	}
	return os.WriteFile(safePath, data, perm)
}

// Create is a wrapper around os.Create that checks permissions via the guard package.
func Create(name string) (*os.File, error) {
	safePath, err := guard.AllowWrite(name)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Warn("create blocked by guard",
				"path", name,
				"execution_id", config.CurrentConfig.ExecutionID,
			)
		}
		return nil, err
	}

	if logger.Log != nil {
		logger.Log.Info("create allowed",
			"path", safePath,
			"execution_id", config.CurrentConfig.ExecutionID,
		)
	}
	return os.Create(safePath)
}

// MkdirAll is a wrapper around os.MkdirAll that checks permissions via the guard package.
func MkdirAll(path string, perm os.FileMode) error {
	safePath, err := guard.AllowWrite(path)
	if err != nil {
		if logger.Log != nil {
			logger.Log.Warn("mkdir blocked by guard",
				"path", path,
				"execution_id", config.CurrentConfig.ExecutionID,
			)
		}
		return err
	}

	if logger.Log != nil {
		logger.Log.Info("mkdir allowed",
			"path", safePath,
			"execution_id", config.CurrentConfig.ExecutionID,
		)
	}
	return os.MkdirAll(safePath, perm)
}
