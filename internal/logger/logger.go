package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

var Log *slog.Logger

// Init initializes the structured logger to write to both stdout and a JSONL file in outputDir.
func Init(outputDir string) error {
	logFilePath := filepath.Join(outputDir, "execution_log.jsonl")
	file, err := os.Create(logFilePath)
	if err != nil {
		return err
	}

	// Use MultiWriter to log to both terminal and file
	mw := io.MultiWriter(os.Stdout, file)

	handler := slog.NewJSONHandler(mw, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	Log = slog.New(handler)
	return nil
}
