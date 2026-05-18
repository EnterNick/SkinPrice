package logx

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRotatingFileWriterRotatesAndCompresses(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "skinprice.log")

	writer, err := newRotatingFileWriter(path, fileOptions{
		MaxSizeBytes: 32,
		MaxBackups:   2,
		MaxAgeDays:   7,
		Compress:     true,
	})
	if err != nil {
		t.Fatalf("newRotatingFileWriter() error = %v", err)
	}
	defer func() {
		_ = writer.Close()
	}()

	for range 4 {
		if _, err := writer.Write([]byte("0123456789abcdef\n")); err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}

	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("ReadDir() error = %v", err)
	}

	var currentFound bool
	var rotatedFound bool
	for _, entry := range entries {
		switch {
		case entry.Name() == "skinprice.log":
			currentFound = true
		case strings.HasPrefix(entry.Name(), "skinprice.log.") && strings.HasSuffix(entry.Name(), ".gz"):
			rotatedFound = true
		}
	}

	if !currentFound {
		t.Fatal("expected current log file to exist")
	}
	if !rotatedFound {
		t.Fatal("expected at least one rotated compressed log file")
	}
}
