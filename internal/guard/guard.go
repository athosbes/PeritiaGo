package guard

import (
	"errors"
	"os"
	"path/filepath"
	"strings"

	"github.com/athosbes/PeritiaGo/internal/config"
	"github.com/athosbes/PeritiaGo/internal/logger"
)

// AllowWrite checks if a write operation to the given path is permitted.
// In forensic mode, writes are only allowed within the designated output directory.
// It returns the resolved, cleaned, and safe absolute path to be used for the operation.
func AllowWrite(path string) (string, error) {
	if config.CurrentConfig == nil {
		return path, nil // Config not initialized yet, assume safe
	}

	// Normalize paths to prevent traversal bypasses
	cleanPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", err
	}

	if !config.CurrentConfig.ForensicMode {
		return cleanPath, nil
	}

	cleanOutput, err := filepath.Abs(filepath.Clean(config.CurrentConfig.OutputDir))
	if err != nil {
		return "", err
	}

	// Resolve symlinks and junctions to prevent escapes
	// First resolve the output directory itself
	resolvedOutput, err := filepath.EvalSymlinks(cleanOutput)
	if err == nil {
		cleanOutput = resolvedOutput
	}

	// For the target path, we resolve it as much as possible.
	// We MUST resolve symlinks for all components of the path to prevent escapes.
	cleanPath = resolvePathRobustly(cleanPath)

	// Resolve the output directory itself again after normalization to handle potential UNC mappings
	if res, err := filepath.EvalSymlinks(cleanOutput); err == nil {
		cleanOutput = res
	}

	// Final normalization again after resolution
	cleanOutput = filepath.Clean(cleanOutput)
	cleanPath = filepath.Clean(cleanPath)

	// In Windows, drive letters and UNC paths need careful handling.
	o := strings.ToLower(normalizeDrivePath(cleanOutput))
	p := strings.ToLower(normalizeDrivePath(cleanPath))

	// Prefix check after absolute resolution is robust against symlink escapes.
	// We use strings.HasPrefix for robust prefix checking.
	sep := string(filepath.Separator)
	prefix := o
	if !strings.HasSuffix(prefix, sep) {
		prefix += sep
	}

	// SPECIAL CASE FOR WINDOWS: EvalSymlinks can resolve to a different drive letter mapping
	// (e.g., a junction to another drive). We must ensure the final path is truly inside.
	if !strings.HasPrefix(p, prefix) && p != o {
		if logger.Log != nil {
			logger.Log.Warn("write blocked in forensic mode (prefix check)",
				"path", cleanPath,
				"normalized_path", p,
				"output_dir", o,
			)
		}
		return "", errors.New("write blocked: attempt to write outside output directory in forensic mode")
	}

	return cleanPath, nil
}

func normalizeDrivePath(path string) string {
	if strings.HasPrefix(path, `\\?\UNC\`) {
		path = `\\` + path[8:]
	} else if strings.HasPrefix(path, `\??\UNC\`) {
		path = `\\` + path[8:]
	} else {
		path = strings.TrimPrefix(path, `\\?\`)
		path = strings.TrimPrefix(path, `\??\`)
	}

	if len(path) >= 2 && path[1] == ':' {
		return strings.ToUpper(string(path[0])) + path[1:]
	}
	return path
}

func resolvePathRobustly(path string) string {
	// First try full resolution
	if resolved, err := filepath.EvalSymlinks(path); err == nil {
		abs, _ := filepath.Abs(resolved)
		return filepath.Clean(abs)
	}

	// If it fails, resolve component by component
	path = filepath.Clean(path)
	vol := filepath.VolumeName(path)
	rest := strings.TrimPrefix(path, vol)
	parts := strings.Split(rest, string(filepath.Separator))

	current := vol
	// Ensure absolute start for Windows drive letters and Unix roots
	if strings.HasPrefix(rest, string(filepath.Separator)) {
		if vol != "" {
			if !strings.HasSuffix(current, string(filepath.Separator)) {
				current += string(filepath.Separator)
			}
		} else {
			current = string(filepath.Separator)
		}
	}

	for i := 0; i < len(parts); i++ {
		if parts[i] == "" {
			continue
		}
		// On Windows, if current is "C:", filepath.Join(current, "Users") results in "C:Users" (relative)
		// instead of "C:\Users" (absolute). We must ensure the separator is present.
		if vol != "" && current == vol && !strings.HasSuffix(current, string(filepath.Separator)) {
			current += string(filepath.Separator)
		}

		next := filepath.Join(current, parts[i])
		
		// TO PREVENT TOCTOU/SYMLINK ESCAPES:
		// We use Lstat to check if the component exists and is a symlink.
		// If it is, we MUST resolve it before continuing.
		if info, err := os.Lstat(next); err == nil {
			if info.Mode()&os.ModeSymlink != 0 {
				if resolved, err := filepath.EvalSymlinks(next); err == nil {
					current = resolved
					continue
				}
			}
			current = next
		} else {
			// If we can't lstat the next part (likely because it doesn't exist yet),
			// we just join the rest.
			for j := i; j < len(parts); j++ {
				current = filepath.Join(current, parts[j])
			}
			break
		}
	}

	abs, err := filepath.Abs(current)
	if err != nil {
		return current
	}
	return filepath.Clean(abs)
}
