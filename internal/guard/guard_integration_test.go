package guard_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/athosbes/PeritiaGo/internal/config"
	"github.com/athosbes/PeritiaGo/internal/filesystem"
)

func TestRealWriteBlocked(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	// Try to write outside using the wrapper
	outsidePath := filepath.Join(tmp, "forbidden.txt")
	err := filesystem.WriteFile(outsidePath, []byte("secret"), 0644)

	if err == nil {
		t.Fatalf("expected error when writing outside output directory via wrapper")
	}

	// Verify file does not exist
	if _, err := os.Stat(outsidePath); !os.IsNotExist(err) {
		t.Fatalf("file should not have been created: %s", outsidePath)
	}
}

func TestRealWriteAllowed(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	// Attempt write through filesystem wrapper
	path := filepath.Join(outputDir, "test.txt")
	err := filesystem.WriteFile(path, []byte("data"), 0644)
	if err != nil {
		t.Fatalf("expected write allowed inside output dir via wrapper, got: %v", err)
	}

	// Attempt write outside
	outsidePath := filepath.Join(tmp, "evil.txt")
	err = filesystem.WriteFile(outsidePath, []byte("data"), 0644)
	if err == nil {
		t.Fatalf("expected write blocked outside output dir via wrapper")
	}
}
