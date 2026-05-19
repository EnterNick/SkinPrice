package main

import (
	"embed"
	"log/slog"
	"os"

	"SkinPrice/skinprice/internal/config"
	"SkinPrice/skinprice/internal/shared/logx"
	"SkinPrice/skinprice/internal/shared/utils"

	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	"github.com/wailsapp/wails/v2/pkg/options/linux"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	utils.LoadDotEnv()
	cfg, err := config.Load()
	if err != nil {
		_, _ = os.Stderr.WriteString("invalid configuration: " + err.Error() + "\n")
		os.Exit(1)
	}
	logger, closer, err := logx.New(logx.Config{
		Level:       cfg.LogLevel,
		Format:      cfg.LogFormat,
		ToFile:      cfg.LogToFile,
		FilePath:    cfg.LogFilePath,
		MaxSizeMB:   cfg.LogMaxSizeMB,
		MaxBackups:  cfg.LogMaxBackups,
		MaxAgeDays:  cfg.LogMaxAgeDays,
		Compress:    cfg.LogCompress,
		AppName:     "SkinPrice",
		Environment: cfg.AppEnv,
	})
	if err != nil {
		_, _ = os.Stderr.WriteString("failed to initialize logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer func() {
		if closeErr := closer.Close(); closeErr != nil {
			slog.Error("failed to close logger", logx.ErrAttrs(closeErr)...)
		}
	}()
	slog.SetDefault(logger)
	logger.Info("logger initialized",
		slog.String("log_level", cfg.LogLevel),
		slog.String("log_format", cfg.LogFormat),
		slog.Bool("log_to_file", cfg.LogToFile),
		slog.String("app_env", cfg.AppEnv),
	)

	app, err := NewApp(logger)
	if err != nil {
		logger.Error("application bootstrap failed", logx.ErrAttrs(err)...)
		os.Exit(1)
	}

	err = wails.Run(&options.App{
		Title:  "SkinPrice",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 27, G: 38, B: 54, A: 1},
		OnStartup:        app.startup,
		Bind: []interface{}{
			app,
		},
		Linux: &linux.Options{
			WebviewGpuPolicy: linux.WebviewGpuPolicyNever,
		},
	})
	if err != nil {
		logger.Error("wails run failed", logx.ErrAttrs(err)...)
		os.Exit(1)
	}
	err = app.Shutdown()

	if err != nil {
		logger.Error("shutdown failed", logx.ErrAttrs(err)...)
		os.Exit(1)
	}
}
