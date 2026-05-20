package version

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

var (
	ErrBootstrapRequired = errors.New("installed version is too old; download a fresh bootstrap package")
	ErrEntrypointInvalid = errors.New("entrypoint must be a safe relative path")
	ErrNoRunnableVersion = errors.New("no runnable version available")
)

type LaunchResult struct {
	Entrypoint string
	Version    string
}

func (s Service) Run(ctx context.Context) (LaunchResult, error) {
	logger := s.logger()
	currentPath := filepath.Join(s.InstallRoot, CurrentStateFile)

	current, currentEntrypoint, currentValid, currentErr := s.loadRunnableCurrent(currentPath)
	if currentErr != nil {
		logger.Warn("failed to load current state", slog.String("error", currentErr.Error()))
	}

	release, err := s.ReleaseProvider.GetLatestRelease(ctx)
	if err != nil {
		if currentValid {
			logger.Warn("failed to fetch latest release, launching current version", slog.String("error", err.Error()))
			return s.launchCurrent(current, currentEntrypoint)
		}
		return LaunchResult{}, fmt.Errorf("%w: %w", ErrNoRunnableVersion, err)
	}

	manifestAsset, err := s.ReleaseProvider.FindAsset(release, UpdateManifestFile)
	if err != nil {
		if currentValid {
			logger.Warn("release manifest asset not found, launching current version", slog.String("error", err.Error()))
			return s.launchCurrent(current, currentEntrypoint)
		}
		return LaunchResult{}, fmt.Errorf("%w: %w", ErrNoRunnableVersion, err)
	}

	manifest, err := s.downloadManifest(ctx, manifestAsset)
	if err != nil {
		if currentValid {
			logger.Warn("failed to read release manifest, launching current version", slog.String("error", err.Error()))
			return s.launchCurrent(current, currentEntrypoint)
		}
		return LaunchResult{}, fmt.Errorf("%w: %w", ErrNoRunnableVersion, err)
	}

	if currentValid && manifest.MinSupportedVersion != "" {
		compareMin, cmpErr := compareVersions(current.Version, manifest.MinSupportedVersion)
		if cmpErr != nil {
			if currentValid {
				return s.launchCurrent(current, currentEntrypoint)
			}
			return LaunchResult{}, cmpErr
		}
		if compareMin < 0 {
			return LaunchResult{}, ErrBootstrapRequired
		}
	}

	asset, err := selectManifestAsset(manifest, s.platformOS(), s.platformArch())
	if err != nil {
		if currentValid {
			logger.Warn("no manifest asset for current platform, launching current version", slog.String("error", err.Error()))
			return s.launchCurrent(current, currentEntrypoint)
		}
		return LaunchResult{}, fmt.Errorf("%w: %w", ErrNoRunnableVersion, err)
	}

	if currentValid {
		cmp, cmpErr := compareVersions(current.Version, manifest.Version)
		if cmpErr != nil {
			return s.launchCurrent(current, currentEntrypoint)
		}
		if cmp >= 0 {
			return s.launchCurrent(current, currentEntrypoint)
		}

		confirmed, confirmErr := s.Prompter.ConfirmUpdate(current.Version, manifest.Version)
		if confirmErr != nil {
			logger.Warn("update prompt failed, launching current version", slog.String("error", confirmErr.Error()))
			return s.launchCurrent(current, currentEntrypoint)
		}
		if !confirmed {
			return s.launchCurrent(current, currentEntrypoint)
		}
	}

	installedCurrent, err := s.installRelease(ctx, release, manifest, asset, current)
	if err != nil {
		if currentValid {
			logger.Warn("update install failed, launching current version", slog.String("error", err.Error()))
			return s.launchCurrent(current, currentEntrypoint)
		}
		return LaunchResult{}, err
	}

	installedEntrypoint, err := resolveEntrypoint(s.InstallRoot, installedCurrent.Entrypoint)
	if err != nil {
		return LaunchResult{}, err
	}
	return s.launchCurrent(installedCurrent, installedEntrypoint)
}

func (s Service) loadRunnableCurrent(currentPath string) (CurrentVersionDTO, string, bool, error) {
	content, err := s.FileStorage.Read(currentPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return CurrentVersionDTO{}, "", false, nil
		}
		return CurrentVersionDTO{}, "", false, err
	}

	var current CurrentVersionDTO
	if err := json.Unmarshal(content, &current); err != nil {
		return CurrentVersionDTO{}, "", false, err
	}

	current.Version = normalizeVersion(current.Version)
	current.Previous = normalizeVersion(current.Previous)

	entrypoint, err := resolveEntrypoint(s.InstallRoot, current.Entrypoint)
	if err != nil {
		return current, "", false, err
	}

	exists, err := s.FileStorage.Exists(entrypoint)
	if err != nil {
		return current, "", false, err
	}
	if !exists {
		return current, "", false, fmt.Errorf("entrypoint does not exist: %s", entrypoint)
	}

	return current, entrypoint, true, nil
}

func (s Service) downloadManifest(ctx context.Context, asset ReleaseAsset) (UpdateManifest, error) {
	content, err := s.downloadBytes(ctx, asset)
	if err != nil {
		return UpdateManifest{}, err
	}

	var manifest UpdateManifest
	if err := json.Unmarshal(content, &manifest); err != nil {
		return UpdateManifest{}, err
	}

	manifest.Version = normalizeVersion(manifest.Version)
	manifest.MinSupportedVersion = normalizeVersion(manifest.MinSupportedVersion)
	if manifest.Channel == "" {
		manifest.Channel = StableChannel
	}
	for i := range manifest.Assets {
		manifest.Assets[i].Entrypoint = filepath.ToSlash(manifest.Assets[i].Entrypoint)
	}

	return manifest, nil
}

func (s Service) installRelease(
	ctx context.Context,
	release ReleaseMeta,
	manifest UpdateManifest,
	manifestAsset ManifestAsset,
	current CurrentVersionDTO,
) (CurrentVersionDTO, error) {
	releaseAsset, err := s.ReleaseProvider.FindAsset(release, manifestAsset.Filename)
	if err != nil {
		return CurrentVersionDTO{}, err
	}

	archiveBytes, err := s.downloadBytes(ctx, releaseAsset)
	if err != nil {
		return CurrentVersionDTO{}, err
	}
	if err := verifyArchive(archiveBytes, manifestAsset); err != nil {
		return CurrentVersionDTO{}, err
	}

	versionsRoot := filepath.Join(s.InstallRoot, "versions")
	tmpRoot := filepath.Join(versionsRoot, ".tmp")
	if err := s.FileStorage.MkdirAll(tmpRoot); err != nil {
		return CurrentVersionDTO{}, err
	}

	stageDir, err := s.FileStorage.TempDir(tmpRoot, manifest.Version+"-")
	if err != nil {
		return CurrentVersionDTO{}, err
	}

	cleanupStage := true
	defer func() {
		if cleanupStage {
			_ = s.FileStorage.RemoveAll(stageDir)
		}
	}()

	archiveName := manifestAsset.Filename
	archivePath := filepath.Join(stageDir, archiveName)
	if err := s.FileStorage.Write(archivePath, archiveBytes, 0o644); err != nil {
		return CurrentVersionDTO{}, err
	}

	extractDir := filepath.Join(stageDir, "content")
	if err := s.FileStorage.MkdirAll(extractDir); err != nil {
		return CurrentVersionDTO{}, err
	}

	if err := extractArchive(s.FileStorage, archivePath, extractDir, archiveName); err != nil {
		return CurrentVersionDTO{}, err
	}

	binaryPath, err := resolveEntrypoint(extractDir, manifestAsset.Entrypoint)
	if err != nil {
		return CurrentVersionDTO{}, err
	}

	exists, err := s.FileStorage.Exists(binaryPath)
	if err != nil {
		return CurrentVersionDTO{}, err
	}
	if !exists {
		return CurrentVersionDTO{}, fmt.Errorf("update entrypoint missing: %s", manifestAsset.Entrypoint)
	}

	if s.platformOS() != "windows" {
		if err := s.FileStorage.Chmod(binaryPath, 0o755); err != nil {
			return CurrentVersionDTO{}, err
		}
	}

	finalDir := filepath.Join(versionsRoot, manifest.Version)
	finalExists, err := s.FileStorage.Exists(finalDir)
	if err != nil {
		return CurrentVersionDTO{}, err
	}

	if !finalExists {
		if err := s.FileStorage.Rename(extractDir, finalDir); err != nil {
			return CurrentVersionDTO{}, err
		}
		cleanupStage = false
	} else {
		cleanupStage = true
	}

	entrypoint := filepath.ToSlash(filepath.Join("versions", manifest.Version, manifestAsset.Entrypoint))
	next := CurrentVersionDTO{
		Version:    manifest.Version,
		Entrypoint: entrypoint,
		Previous:   current.Version,
		UpdatedAt:  s.now().UTC(),
	}

	payload, err := json.MarshalIndent(next, "", "  ")
	if err != nil {
		return CurrentVersionDTO{}, err
	}
	payload = append(payload, '\n')

	if err := s.FileStorage.WriteAtomic(filepath.Join(s.InstallRoot, CurrentStateFile), payload); err != nil {
		return CurrentVersionDTO{}, err
	}

	return next, nil
}

func (s Service) launchCurrent(current CurrentVersionDTO, entrypoint string) (LaunchResult, error) {
	if err := s.AppRunner.Start(entrypoint); err != nil {
		return LaunchResult{}, err
	}
	return LaunchResult{Entrypoint: entrypoint, Version: current.Version}, nil
}

func (s Service) logger() *slog.Logger {
	if s.Logger != nil {
		return s.Logger
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

func (s Service) now() time.Time {
	if s.Now != nil {
		return s.Now()
	}
	return time.Now()
}

func (s Service) platformOS() string {
	if s.PlatformOS != "" {
		return s.PlatformOS
	}
	return runtime.GOOS
}

func (s Service) platformArch() string {
	if s.PlatformArch != "" {
		return s.PlatformArch
	}
	return runtime.GOARCH
}

func (s Service) downloadBytes(ctx context.Context, asset ReleaseAsset) ([]byte, error) {
	reader, err := s.Downloader.Download(ctx, asset)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = reader.Close()
	}()
	return io.ReadAll(reader)
}

func selectManifestAsset(manifest UpdateManifest, targetOS, targetArch string) (ManifestAsset, error) {
	for _, asset := range manifest.Assets {
		if asset.OS == targetOS && asset.Arch == targetArch {
			return asset, nil
		}
	}
	return ManifestAsset{}, fmt.Errorf("asset for %s/%s not found", targetOS, targetArch)
}

func resolveEntrypoint(root, entrypoint string) (string, error) {
	entrypoint = strings.TrimSpace(filepath.ToSlash(entrypoint))
	if entrypoint == "" {
		return "", ErrEntrypointInvalid
	}
	if filepath.IsAbs(entrypoint) {
		return "", ErrEntrypointInvalid
	}

	cleanRel := filepath.Clean(filepath.FromSlash(entrypoint))
	if cleanRel == "." || cleanRel == "" || cleanRel == ".." {
		return "", ErrEntrypointInvalid
	}
	if strings.HasPrefix(cleanRel, ".."+string(os.PathSeparator)) {
		return "", ErrEntrypointInvalid
	}

	joined := filepath.Join(root, cleanRel)
	rootAbs, err := filepath.Abs(root)
	if err != nil {
		return "", err
	}
	joinedAbs, err := filepath.Abs(joined)
	if err != nil {
		return "", err
	}
	rel, err := filepath.Rel(rootAbs, joinedAbs)
	if err != nil {
		return "", err
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", ErrEntrypointInvalid
	}
	return joinedAbs, nil
}

func normalizeVersion(value string) string {
	return strings.TrimPrefix(strings.TrimSpace(value), "v")
}

func compareVersions(left, right string) (int, error) {
	leftParts, err := parseSemver(normalizeVersion(left))
	if err != nil {
		return 0, err
	}
	rightParts, err := parseSemver(normalizeVersion(right))
	if err != nil {
		return 0, err
	}

	for i := 0; i < 3; i++ {
		if leftParts[i] < rightParts[i] {
			return -1, nil
		}
		if leftParts[i] > rightParts[i] {
			return 1, nil
		}
	}
	return 0, nil
}

func parseSemver(value string) ([3]int, error) {
	var result [3]int
	parts := strings.Split(value, ".")
	if len(parts) != 3 {
		return result, fmt.Errorf("invalid version: %s", value)
	}
	for i, part := range parts {
		num, err := strconv.Atoi(part)
		if err != nil || num < 0 {
			return result, fmt.Errorf("invalid version: %s", value)
		}
		result[i] = num
	}
	return result, nil
}

func verifyArchive(content []byte, asset ManifestAsset) error {
	if asset.Size > 0 && int64(len(content)) != asset.Size {
		return fmt.Errorf("unexpected asset size: got %d want %d", len(content), asset.Size)
	}
	sum := sha256.Sum256(content)
	actual := hex.EncodeToString(sum[:])
	if asset.Sha256 != "" && !strings.EqualFold(actual, asset.Sha256) {
		return fmt.Errorf("sha256 mismatch: got %s want %s", actual, asset.Sha256)
	}
	return nil
}

func extractArchive(storage FileStorage, archivePath, dst, archiveName string) error {
	switch {
	case strings.HasSuffix(strings.ToLower(archiveName), ".zip"):
		return storage.ExtractZip(archivePath, dst)
	case strings.HasSuffix(strings.ToLower(archiveName), ".tar.gz"):
		return storage.ExtractTarGz(archivePath, dst)
	default:
		return fmt.Errorf("unsupported archive format: %s", archiveName)
	}
}

func BuildZip(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	writer := zip.NewWriter(&buf)
	for name, content := range files {
		item, err := writer.Create(name)
		if err != nil {
			return nil, err
		}
		if _, err := item.Write(content); err != nil {
			return nil, err
		}
	}
	if err := writer.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func BuildTarGz(files map[string][]byte) ([]byte, error) {
	var buf bytes.Buffer
	gz := gzip.NewWriter(&buf)
	tw := tar.NewWriter(gz)
	for name, content := range files {
		header := &tar.Header{
			Name: name,
			Mode: 0o644,
			Size: int64(len(content)),
		}
		if err := tw.WriteHeader(header); err != nil {
			return nil, err
		}
		if _, err := tw.Write(content); err != nil {
			return nil, err
		}
	}
	if err := tw.Close(); err != nil {
		return nil, err
	}
	if err := gz.Close(); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
