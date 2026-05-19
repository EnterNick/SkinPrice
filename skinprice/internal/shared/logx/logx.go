package logx

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"

	"SkinPrice/skinprice/internal/shared/errx"
)

type Config struct {
	Level       string
	Format      string
	ToFile      bool
	FilePath    string
	MaxSizeMB   int
	MaxBackups  int
	MaxAgeDays  int
	Compress    bool
	AppName     string
	Environment string
}

func New(cfg Config) (*slog.Logger, io.Closer, error) {
	if cfg.AppName == "" {
		cfg.AppName = "skinprice"
	}

	level := parseLevel(cfg.Level)
	consoleHandler := newHandler(os.Stdout, cfg.Format, level, true)

	handlers := []slog.Handler{consoleHandler}
	closers := make([]io.Closer, 0, 1)

	if cfg.ToFile {
		path, err := resolveLogPath(cfg)
		if err != nil {
			return nil, nil, fmt.Errorf("resolve log path: %w", err)
		}
		writer, err := newRotatingFileWriter(path, fileOptions{
			MaxSizeBytes: int64(defaultInt(cfg.MaxSizeMB, 20)) * 1024 * 1024,
			MaxBackups:   defaultInt(cfg.MaxBackups, 5),
			MaxAgeDays:   defaultInt(cfg.MaxAgeDays, 14),
			Compress:     cfg.Compress,
		})
		if err != nil {
			return nil, nil, fmt.Errorf("init file logger: %w", err)
		}
		closers = append(closers, writer)
		handlers = append(handlers, newHandler(writer, "json", level, false))
		cfg.FilePath = path
	}

	logger := slog.New(&fanoutHandler{handlers: handlers}).With(
		slog.String("app", cfg.AppName),
		slog.String("env", strings.TrimSpace(cfg.Environment)),
	)

	return logger, multiCloser(closers), nil
}

func WithComponent(logger *slog.Logger, component string) *slog.Logger {
	return Safe(logger).With(slog.String("component", component))
}

func Safe(logger *slog.Logger) *slog.Logger {
	if logger != nil {
		return logger
	}
	if slog.Default() != nil {
		return slog.Default()
	}
	return slog.New(slog.NewTextHandler(os.Stdout, nil))
}

func ErrAttrs(err error) []any {
	if err == nil {
		return nil
	}

	attrs := []any{slog.String("error", err.Error())}

	var ex *errx.Error
	if errors.As(err, &ex) {
		if ex.Op != "" {
			attrs = append(attrs, slog.String("op", ex.Op))
		}
		if ex.Code != "" {
			attrs = append(attrs, slog.String("code", string(ex.Code)))
		}
		for key, value := range SanitizeAttrs(ex.Fields) {
			attrs = append(attrs, slog.Any(key, value))
		}
	}

	return attrs
}

func resolveLogPath(cfg Config) (string, error) {
	if strings.TrimSpace(cfg.FilePath) != "" {
		return filepath.Abs(cfg.FilePath)
	}

	baseDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}

	appDir := cfg.AppName
	if appDir == "" {
		appDir = "skinprice"
	}

	return filepath.Join(baseDir, appDir, "logs", "skinprice.log"), nil
}

func parseLevel(value string) slog.Level {
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "debug":
		return slog.LevelDebug
	case "warn", "warning":
		return slog.LevelWarn
	case "error":
		return slog.LevelError
	default:
		return slog.LevelInfo
	}
}

func newHandler(w io.Writer, format string, level slog.Level, addSource bool) slog.Handler {
	opts := &slog.HandlerOptions{Level: level, AddSource: addSource}
	if strings.EqualFold(format, "json") {
		return slog.NewJSONHandler(w, opts)
	}
	return slog.NewTextHandler(w, opts)
}

func defaultInt(value, fallback int) int {
	if value > 0 {
		return value
	}
	return fallback
}

type fanoutHandler struct {
	handlers []slog.Handler
}

func (h *fanoutHandler) Enabled(ctx context.Context, level slog.Level) bool {
	for _, handler := range h.handlers {
		if handler.Enabled(ctx, level) {
			return true
		}
	}
	return false
}

func (h *fanoutHandler) Handle(ctx context.Context, record slog.Record) error {
	var result error
	for _, handler := range h.handlers {
		if !handler.Enabled(ctx, record.Level) {
			continue
		}
		if err := handler.Handle(ctx, record.Clone()); err != nil {
			result = errors.Join(result, err)
		}
	}
	return result
}

func (h *fanoutHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithAttrs(attrs))
	}
	return &fanoutHandler{handlers: handlers}
}

func (h *fanoutHandler) WithGroup(name string) slog.Handler {
	handlers := make([]slog.Handler, 0, len(h.handlers))
	for _, handler := range h.handlers {
		handlers = append(handlers, handler.WithGroup(name))
	}
	return &fanoutHandler{handlers: handlers}
}

type multiCloser []io.Closer

func (m multiCloser) Close() error {
	var result error
	for _, closer := range m {
		if closer == nil {
			continue
		}
		if err := closer.Close(); err != nil {
			result = errors.Join(result, err)
		}
	}
	return result
}
