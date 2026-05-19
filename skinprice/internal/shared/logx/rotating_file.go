package logx

import (
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"time"
)

type fileOptions struct {
	MaxSizeBytes int64
	MaxBackups   int
	MaxAgeDays   int
	Compress     bool
}

type rotatingFileWriter struct {
	path string
	opts fileOptions

	mu   sync.Mutex
	file *os.File
	size int64
}

func newRotatingFileWriter(path string, opts fileOptions) (*rotatingFileWriter, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return nil, err
	}

	w := &rotatingFileWriter{path: path, opts: opts}
	if err := w.openExisting(); err != nil {
		return nil, err
	}
	return w, nil
}

func (w *rotatingFileWriter) Write(p []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()

	if err := w.rotateIfNeeded(int64(len(p))); err != nil {
		return 0, err
	}

	n, err := w.file.Write(p)
	w.size += int64(n)
	return n, err
}

func (w *rotatingFileWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.file == nil {
		return nil
	}
	return w.file.Close()
}

func (w *rotatingFileWriter) openExisting() error {
	file, err := os.OpenFile(w.path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}

	info, err := file.Stat()
	if err != nil {
		_ = file.Close()
		return err
	}

	w.file = file
	w.size = info.Size()
	return nil
}

func (w *rotatingFileWriter) rotateIfNeeded(incoming int64) error {
	if w.file == nil {
		if err := w.openExisting(); err != nil {
			return err
		}
	}

	if w.opts.MaxSizeBytes <= 0 || w.size+incoming <= w.opts.MaxSizeBytes {
		return nil
	}

	if err := w.file.Close(); err != nil {
		return err
	}

	rotatedName := fmt.Sprintf("%s.%s", w.path, time.Now().UTC().Format("20060102-150405"))
	if err := os.Rename(w.path, rotatedName); err != nil && !os.IsNotExist(err) {
		return err
	}

	if err := w.openExisting(); err != nil {
		return err
	}

	if w.opts.Compress {
		if err := compressFile(rotatedName); err != nil {
			return err
		}
	}

	return w.cleanup()
}

func (w *rotatingFileWriter) cleanup() error {
	entries, err := listRotatedFiles(w.path)
	if err != nil {
		return err
	}

	if w.opts.MaxAgeDays > 0 {
		cutoff := time.Now().AddDate(0, 0, -w.opts.MaxAgeDays)
		filtered := entries[:0]
		for _, entry := range entries {
			if entry.modTime.Before(cutoff) {
				_ = os.Remove(entry.path)
				continue
			}
			filtered = append(filtered, entry)
		}
		entries = filtered
	}

	if w.opts.MaxBackups > 0 && len(entries) > w.opts.MaxBackups {
		for _, entry := range entries[:len(entries)-w.opts.MaxBackups] {
			_ = os.Remove(entry.path)
		}
	}

	return nil
}

type rotatedFile struct {
	path    string
	modTime time.Time
}

func listRotatedFiles(path string) ([]rotatedFile, error) {
	dir := filepath.Dir(path)
	base := filepath.Base(path) + "."

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	files := make([]rotatedFile, 0, len(entries))
	for _, entry := range entries {
		name := entry.Name()
		if entry.IsDir() || !strings.HasPrefix(name, base) {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			return nil, err
		}

		files = append(files, rotatedFile{
			path:    filepath.Join(dir, name),
			modTime: info.ModTime(),
		})
	}

	sort.Slice(files, func(i, j int) bool {
		return files[i].modTime.Before(files[j].modTime)
	})

	return files, nil
}

func compressFile(path string) error {
	src, err := os.Open(path)
	if err != nil {
		return err
	}
	defer func() {
		_ = src.Close()
	}()

	dst, err := os.Create(path + ".gz")
	if err != nil {
		return err
	}

	gz := gzip.NewWriter(dst)
	if _, err := io.Copy(gz, src); err != nil {
		_ = gz.Close()
		_ = dst.Close()
		return err
	}
	if err := gz.Close(); err != nil {
		_ = dst.Close()
		return err
	}
	if err := dst.Close(); err != nil {
		return err
	}

	return os.Remove(path)
}
