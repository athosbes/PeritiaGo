package executor

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// SystemBinaries defines absolute paths for critical system tools to prevent hijacking.
var SystemBinaries = map[string]string{
	"powershell": filepath.Join(os.Getenv("SystemRoot"), "System32", "WindowsPowerShell", "v1.0", "powershell.exe"),
	"wmic":       filepath.Join(os.Getenv("SystemRoot"), "System32", "wbem", "wmic.exe"),
	"wevtutil":   filepath.Join(os.Getenv("SystemRoot"), "System32", "wevtutil.exe"),
	"control":    filepath.Join(os.Getenv("SystemRoot"), "System32", "control.exe"),
}

// Execute runs a command with the given name and arguments.
// If the name is a key in SystemBinaries, the absolute path is used.
func Execute(name string, args ...string) ([]byte, error) {
	return ExecuteContext(context.Background(), name, args...)
}

// ExecuteContext runs a command with context, absolute path enforcement, and secure argument handling.
func ExecuteContext(ctx context.Context, name string, args ...string) ([]byte, error) {
	binPath, ok := SystemBinaries[strings.ToLower(name)]
	if !ok {
		// If not a known system binary, we check if it's an absolute path already
		if !filepath.IsAbs(name) {
			return nil, fmt.Errorf("insecure command execution: %s is not an absolute path and not in allowed system binaries", name)
		}
		binPath = name
	}

	// Verify binary exists
	if _, err := os.Stat(binPath); err != nil {
		return nil, fmt.Errorf("binary not found: %s", binPath)
	}

	cmd := exec.CommandContext(ctx, binPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("command execution failed (%s): %w\nOutput: %s", binPath, err, string(output))
	}

	return output, nil
}

// Start runs a command asynchronously (useful for GUI tools like control panel).
func Start(name string, args ...string) (*exec.Cmd, error) {
	binPath, ok := SystemBinaries[strings.ToLower(name)]
	if !ok {
		if !filepath.IsAbs(name) {
			return nil, fmt.Errorf("insecure command start: %s is not an absolute path", name)
		}
		binPath = name
	}

	cmd := exec.Command(binPath, args...)
	if err := cmd.Start(); err != nil {
		return nil, fmt.Errorf("failed to start command (%s): %w", binPath, err)
	}
	return cmd, nil
}
