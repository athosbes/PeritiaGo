package capture

import (
	"fmt"
	"image/png"
	"os"
	"path/filepath"

	"github.com/kbinani/screenshot"
)

// CaptureScreen takes a screenshot of the main display and saves it to the specified path.
func CaptureScreen(outputPath string) error {
	// Assumes primary display
	bounds := screenshot.GetDisplayBounds(0)

	img, err := screenshot.CaptureRect(bounds)
	if err != nil {
		return fmt.Errorf("failed to capture screen: %w", err)
	}

	// Ensure directory exists
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory for screenshot: %w", err)
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create screenshot file: %w", err)
	}
	defer file.Close()

	if err := png.Encode(file, img); err != nil {
		return fmt.Errorf("failed to encode screenshot as PNG: %w", err)
	}

	return nil
}
