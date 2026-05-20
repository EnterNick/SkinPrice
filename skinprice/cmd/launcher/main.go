package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"

	filedownloader "SkinPrice/skinprice/internal/adapters/filedownloader"
	githubrelease "SkinPrice/skinprice/internal/adapters/githubrelease"
	"SkinPrice/skinprice/internal/adapters/osfile"
	processlauncher "SkinPrice/skinprice/internal/adapters/processlauncher"
	"SkinPrice/skinprice/internal/adapters/prompt"
	appversion "SkinPrice/skinprice/internal/application/version"
	"SkinPrice/skinprice/internal/shared/logx"
	"SkinPrice/skinprice/internal/shared/utils"
)

const defaultRepo = "EnterNick/SkinPrice"

func main() {
	utils.LoadDotEnv()

	executablePath, err := os.Executable()
	if err != nil {
		_, _ = os.Stderr.WriteString("failed to resolve executable path: " + err.Error() + "\n")
		os.Exit(1)
	}
	installRoot := filepath.Dir(executablePath)

	logger, closer, err := logx.New(logx.Config{
		Level:       "info",
		Format:      "text",
		ToFile:      true,
		FilePath:    filepath.Join(installRoot, "logs", "launcher.log"),
		MaxSizeMB:   10,
		MaxBackups:  5,
		MaxAgeDays:  14,
		Compress:    true,
		AppName:     "skinprice-launcher",
		Environment: "production",
	})
	if err != nil {
		_, _ = os.Stderr.WriteString("failed to initialize launcher logger: " + err.Error() + "\n")
		os.Exit(1)
	}
	defer func() {
		_ = closer.Close()
	}()
	slog.SetDefault(logger)

	httpClient := &http.Client{Timeout: 30 * time.Second}
	service := appversion.Service{
		InstallRoot: installRoot,
		Logger:      logger,
		ReleaseProvider: githubrelease.Client{
			BaseURL: "https://api.github.com",
			Repo:    repoName(),
			HTTP:    httpClient,
		},
		Downloader:  filedownloader.Downloader{HTTP: httpClient},
		FileStorage: osfile.Storage{},
		AppRunner:   processlauncher.Runner{},
		Prompter:    prompt.Prompter{},
	}

	if _, err := service.Run(context.Background()); err != nil {
		logger.Error("launcher failed", slog.String("error", err.Error()))
		_, _ = os.Stderr.WriteString("launcher failed: " + err.Error() + "\n")
		os.Exit(1)
	}
}

func repoName() string {
	if value := os.Getenv("SKINPRICE_GITHUB_REPO"); value != "" {
		return value
	}
	return defaultRepo
}
