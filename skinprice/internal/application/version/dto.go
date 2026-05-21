package version

import (
	"context"
	"io"
	"log/slog"
	"os"
	"time"
)

const (
	CurrentStateFile   = "current.json"
	UpdateManifestFile = "update-manifest.json"
	StableChannel      = "stable"
)

type CurrentVersionDTO struct {
	Version    string    `json:"version"`
	Entrypoint string    `json:"entrypoint"`
	Previous   string    `json:"previous"`
	UpdatedAt  time.Time `json:"updated_at"`
}

type UpdateManifest struct {
	Version             string          `json:"version"`
	Channel             string          `json:"channel"`
	MinSupportedVersion string          `json:"min_supported_version"`
	ReleaseNotes        string          `json:"release_notes"`
	PublishedAt         time.Time       `json:"published_at"`
	Assets              []ManifestAsset `json:"assets"`
}

type ManifestAsset struct {
	OS         string `json:"os"`
	Arch       string `json:"arch"`
	Filename   string `json:"filename"`
	Sha256     string `json:"sha256"`
	Size       int64  `json:"size"`
	Entrypoint string `json:"entrypoint"`
}

type ReleaseMeta struct {
	TagName string
	Assets  []ReleaseAsset
}

type ReleaseAsset struct {
	Name        string
	DownloadURL string
	Size        int64
}

type ReleaseProvider interface {
	GetLatestRelease(ctx context.Context) (ReleaseMeta, error)
	FindAsset(release ReleaseMeta, name string) (ReleaseAsset, error)
}

type AssetDownloader interface {
	Download(ctx context.Context, asset ReleaseAsset) (io.ReadCloser, error)
}

type FileStorage interface {
	Read(path string) ([]byte, error)
	Write(path string, content []byte, mode os.FileMode) error
	WriteAtomic(path string, content []byte) error
	MkdirAll(path string) error
	Exists(path string) (bool, error)
	Rename(oldPath, newPath string) error
	RemoveAll(path string) error
	ExtractZip(src, dst string) error
	ExtractTarGz(src, dst string) error
	Chmod(path string, mode os.FileMode) error
	TempDir(dir, pattern string) (string, error)
}

type AppRunner interface {
	Start(entrypoint string) error
}

type UpdatePrompter interface {
	ShowCheckingForUpdates() (io.Closer, error)
	ConfirmUpdate(currentVersion, newVersion string) (bool, error)
	NotifyUpdateFailed(currentVersion string, updateErr error) error
	NotifyUpdateSuccess(previousVersion, newVersion string) error
}

type Service struct {
	InstallRoot     string
	PlatformOS      string
	PlatformArch    string
	Logger          *slog.Logger
	ReleaseProvider ReleaseProvider
	Downloader      AssetDownloader
	FileStorage     FileStorage
	AppRunner       AppRunner
	Prompter        UpdatePrompter
	Now             func() time.Time
}
