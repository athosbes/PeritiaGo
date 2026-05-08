package guard

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/athosbes/PeritiaGo/internal/config"
)

func TestAllowWriteInsideOutput(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	path := filepath.Join(outputDir, "report.json")
	_, err := AllowWrite(path)
	if err != nil {
		t.Fatalf("expected write allowed inside output dir, got: %v", err)
	}
}

func TestAllowWriteOutsideOutput(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	outsidePath := filepath.Join(tmp, "stolen_data.txt")
	_, err := AllowWrite(outsidePath)
	if err == nil {
		t.Fatalf("expected block for path outside output directory")
	}
}

func TestAllowWriteTraversal(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	// Attempt traversal: output/../evil.txt -> should be tmp/evil.txt
	traversalPath := filepath.Join(outputDir, "..", "evil.txt")
	_, err := AllowWrite(traversalPath)
	if err == nil {
		t.Fatalf("expected block for path traversal attempt")
	}
}

func TestForensicModeDisabled(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: false,
		OutputDir:    outputDir,
	}

	outsidePath := filepath.Join(tmp, "anywhere.txt")
	_, err := AllowWrite(outsidePath)
	if err != nil {
		t.Fatalf("expected write allowed when forensic mode is disabled, got: %v", err)
	}
}

func TestSymlinkEscape(t *testing.T) {
	tmp := t.TempDir()

	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	outsideDir := filepath.Join(tmp, "outside")
	os.MkdirAll(outsideDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	link := filepath.Join(outputDir, "link_to_outside")
	// On Windows, symlink creation might require Admin or Developer Mode.
	// We'll skip this test if we can't create a real symlink that resolves.
	_ = os.Symlink(outsideDir, link)

	// Verify it's a real resolving symlink
	if res, err := filepath.EvalSymlinks(link); err != nil || res == link {
		t.Skipf("skipping symlink test: could not create a resolving symlink (likely permission issue)")
	}

	evilPath := filepath.Join(link, "evil.txt")
	_, err := AllowWrite(evilPath)
	if err == nil {
		t.Fatalf("expected block for symlink escape pointing outside output directory")
	}
}

func TestAllowWriteUNC(t *testing.T) {
	tmp := t.TempDir()
	outputDir := filepath.Join(tmp, "output")
	os.MkdirAll(outputDir, 0755)

	config.CurrentConfig = &config.AppConfig{
		ForensicMode: true,
		OutputDir:    outputDir,
	}

	// Attempting a UNC-like path that points elsewhere
	uncPath := `\\localhost\C$\Windows\temp.txt`
	if _, err := AllowWrite(uncPath); err == nil {
		t.Errorf("expected block for UNC path: %s", uncPath)
	}
}
