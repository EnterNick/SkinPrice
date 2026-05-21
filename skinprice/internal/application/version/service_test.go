package version

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"SkinPrice/skinprice/internal/adapters/osfile"
)

func TestRunLaunchesCurrentWhenGitHubUnavailable(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "versions/0.1.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.1.0", "skinprice"))

	runner := &fakeRunner{}
	service := newService(root, fakeReleaseProvider{err: errors.New("offline")}, fakeDownloader{}, runner, &fakePrompter{})

	result, err := service.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Version != "0.1.0" {
		t.Fatalf("Run() version = %s", result.Version)
	}
	if len(runner.started) != 1 {
		t.Fatalf("runner started = %d, want 1", len(runner.started))
	}
}

func TestRunDoesNotOfferUpdateWhenManifestVersionMatches(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.2.0",
		Entrypoint: "versions/0.2.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.2.0", "skinprice"))

	release, downloader := manifestOnlyRelease(t, UpdateManifest{
		Version:             "0.2.0",
		Channel:             StableChannel,
		MinSupportedVersion: "0.1.0",
		PublishedAt:         time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
		Assets: []ManifestAsset{{
			OS: "linux", Arch: "amd64", Filename: "skinprice-linux-amd64.tar.gz", Entrypoint: "skinprice",
		}},
	})

	prompter := &fakePrompter{confirmResult: true}
	runner := &fakeRunner{}
	service := newService(root, release, downloader, runner, prompter)

	if _, err := service.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if prompter.calls != 0 {
		t.Fatalf("prompt calls = %d, want 0", prompter.calls)
	}
}

func TestRunFallsBackWhenUserDeclinesUpdate(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "versions/0.1.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.1.0", "skinprice"))

	updateBytes, manifest := tarballManifest(t, "0.2.0", "skinprice", []byte("linux-binary"))
	release, downloader := releaseWithAssets(t, manifest, map[string][]byte{
		UpdateManifestFile:             marshalManifest(t, manifest),
		"skinprice-linux-amd64.tar.gz": updateBytes,
	})

	runner := &fakeRunner{}
	prompter := &fakePrompter{confirmResult: false}
	service := newService(root, release, downloader, runner, prompter)

	if _, err := service.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	current := readCurrentState(t, root)
	if current.Version != "0.1.0" {
		t.Fatalf("current version = %s, want 0.1.0", current.Version)
	}
	if len(runner.started) != 1 {
		t.Fatalf("runner started = %d, want 1", len(runner.started))
	}
}

func TestRunInstallsUpdateAndLaunchesNewVersion(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "versions/0.1.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.1.0", "skinprice"))

	updateBytes, manifest := tarballManifest(t, "0.2.0", "skinprice", []byte("linux-binary"))
	release, downloader := releaseWithAssets(t, manifest, map[string][]byte{
		UpdateManifestFile:             marshalManifest(t, manifest),
		"skinprice-linux-amd64.tar.gz": updateBytes,
	})

	runner := &fakeRunner{}
	prompter := &fakePrompter{confirmResult: true}
	service := newService(root, release, downloader, runner, prompter)

	result, err := service.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Version != "0.2.0" {
		t.Fatalf("Run() version = %s, want 0.2.0", result.Version)
	}

	current := readCurrentState(t, root)
	if current.Version != "0.2.0" || current.Previous != "0.1.0" {
		t.Fatalf("current state = %+v", current)
	}
	if !strings.HasSuffix(current.Entrypoint, "versions/0.2.0/skinprice") {
		t.Fatalf("current entrypoint = %s", current.Entrypoint)
	}

	info, err := os.Stat(filepath.Join(root, "versions", "0.2.0", "skinprice"))
	if err != nil {
		t.Fatalf("Stat() error = %v", err)
	}
	if info.Mode().Perm()&0o111 == 0 {
		t.Fatalf("binary mode = %v, want executable", info.Mode().Perm())
	}
	if len(prompter.succeeded) != 1 || prompter.succeeded[0] != "0.1.0|0.2.0" {
		t.Fatalf("success notifications = %#v, want [0.1.0|0.2.0]", prompter.succeeded)
	}
	if prompter.checkingShown != 1 || prompter.checkingClosed != 1 {
		t.Fatalf("checking window lifecycle = shown:%d closed:%d, want 1/1", prompter.checkingShown, prompter.checkingClosed)
	}
}

func TestRunFallsBackWhenManifestMissing(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "versions/0.1.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.1.0", "skinprice"))

	release := fakeReleaseProvider{release: ReleaseMeta{
		TagName: "v0.2.0",
		Assets: []ReleaseAsset{{
			Name: "skinprice-linux-amd64.tar.gz", DownloadURL: "skinprice-linux-amd64.tar.gz",
		}},
	}}

	runner := &fakeRunner{}
	service := newService(root, release, fakeDownloader{}, runner, &fakePrompter{})
	if _, err := service.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if len(runner.started) != 1 {
		t.Fatalf("runner started = %d, want 1", len(runner.started))
	}
}

func TestRunBlocksUpdateBelowMinSupportedVersion(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "versions/0.1.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.1.0", "skinprice"))

	release, downloader := manifestOnlyRelease(t, UpdateManifest{
		Version:             "0.3.0",
		Channel:             StableChannel,
		MinSupportedVersion: "0.2.0",
		PublishedAt:         time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
		Assets: []ManifestAsset{{
			OS: "linux", Arch: "amd64", Filename: "skinprice-linux-amd64.tar.gz", Entrypoint: "skinprice",
		}},
	})

	service := newService(root, release, downloader, &fakeRunner{}, &fakePrompter{})
	_, err := service.Run(context.Background())
	if !errors.Is(err, ErrBootstrapRequired) {
		t.Fatalf("Run() error = %v, want ErrBootstrapRequired", err)
	}
}

func TestRunRecoversWhenCurrentStateMissing(t *testing.T) {
	root := t.TempDir()
	updateBytes, manifest := zipManifest(t, "0.2.0", "SkinPrice.exe", []byte("windows-binary"))
	release, downloader := releaseWithAssets(t, manifest, map[string][]byte{
		UpdateManifestFile:            marshalManifest(t, manifest),
		"skinprice-windows-amd64.zip": updateBytes,
	})

	runner := &fakeRunner{}
	service := newService(root, release, downloader, runner, &fakePrompter{})
	service.PlatformOS = "windows"

	result, err := service.Run(context.Background())
	if err != nil {
		t.Fatalf("Run() error = %v", err)
	}
	if result.Version != "0.2.0" {
		t.Fatalf("Run() version = %s", result.Version)
	}
	current := readCurrentState(t, root)
	if current.Version != "0.2.0" {
		t.Fatalf("current version = %s", current.Version)
	}
}

func TestRunRejectsCurrentEntrypointTraversal(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "../outside",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})

	service := newService(root, fakeReleaseProvider{err: errors.New("offline")}, fakeDownloader{}, &fakeRunner{}, &fakePrompter{})
	_, err := service.Run(context.Background())
	if !errors.Is(err, ErrNoRunnableVersion) {
		t.Fatalf("Run() error = %v, want ErrNoRunnableVersion", err)
	}
}

func TestRunKeepsCurrentStateWhenChecksumMismatch(t *testing.T) {
	root := t.TempDir()
	writeCurrentState(t, root, CurrentVersionDTO{
		Version:    "0.1.0",
		Entrypoint: "versions/0.1.0/skinprice",
		UpdatedAt:  time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
	})
	writeExecutable(t, filepath.Join(root, "versions", "0.1.0", "skinprice"))

	updateBytes, manifest := tarballManifest(t, "0.2.0", "skinprice", []byte("linux-binary"))
	manifest.Assets[0].Sha256 = strings.Repeat("0", 64)
	release, downloader := releaseWithAssets(t, manifest, map[string][]byte{
		UpdateManifestFile:             marshalManifest(t, manifest),
		"skinprice-linux-amd64.tar.gz": updateBytes,
	})

	runner := &fakeRunner{}
	prompter := &fakePrompter{confirmResult: true}
	service := newService(root, release, downloader, runner, prompter)
	if _, err := service.Run(context.Background()); err != nil {
		t.Fatalf("Run() error = %v", err)
	}

	current := readCurrentState(t, root)
	if current.Version != "0.1.0" {
		t.Fatalf("current version = %s, want 0.1.0", current.Version)
	}
	if len(prompter.failed) != 1 || !strings.HasPrefix(prompter.failed[0], "0.1.0|") {
		t.Fatalf("failure notifications = %#v, want one for current version 0.1.0", prompter.failed)
	}
	if prompter.checkingShown != 1 || prompter.checkingClosed != 1 {
		t.Fatalf("checking window lifecycle = shown:%d closed:%d, want 1/1", prompter.checkingShown, prompter.checkingClosed)
	}
}

func TestCompareVersionsNormalizesTagPrefix(t *testing.T) {
	cmp, err := compareVersions("v0.2.0", "0.2.0")
	if err != nil {
		t.Fatalf("compareVersions() error = %v", err)
	}
	if cmp != 0 {
		t.Fatalf("compareVersions() = %d, want 0", cmp)
	}
}

func TestResolveEntrypointRejectsAbsoluteAndTraversalPaths(t *testing.T) {
	root := t.TempDir()
	if _, err := resolveEntrypoint(root, "/tmp/skinprice"); !errors.Is(err, ErrEntrypointInvalid) {
		t.Fatalf("absolute path error = %v", err)
	}
	if _, err := resolveEntrypoint(root, "../skinprice"); !errors.Is(err, ErrEntrypointInvalid) {
		t.Fatalf("traversal path error = %v", err)
	}
}

type fakeReleaseProvider struct {
	release ReleaseMeta
	err     error
}

func (f fakeReleaseProvider) GetLatestRelease(context.Context) (ReleaseMeta, error) {
	return f.release, f.err
}

func (f fakeReleaseProvider) FindAsset(release ReleaseMeta, name string) (ReleaseAsset, error) {
	for _, asset := range release.Assets {
		if asset.Name == name {
			return asset, nil
		}
	}
	return ReleaseAsset{}, errors.New("not found")
}

type fakeDownloader map[string][]byte

func (f fakeDownloader) Download(_ context.Context, asset ReleaseAsset) (io.ReadCloser, error) {
	content, ok := f[asset.DownloadURL]
	if !ok {
		return nil, errors.New("missing asset")
	}
	return io.NopCloser(bytes.NewReader(content)), nil
}

type fakeRunner struct {
	started []string
}

func (f *fakeRunner) Start(entrypoint string) error {
	f.started = append(f.started, entrypoint)
	return nil
}

type fakePrompter struct {
	confirmResult  bool
	calls          int
	failed         []string
	succeeded      []string
	checkingShown  int
	checkingClosed int
}

type fakeCloser struct {
	onClose func() error
}

func (f fakeCloser) Close() error {
	if f.onClose != nil {
		return f.onClose()
	}
	return nil
}

func (f *fakePrompter) ShowCheckingForUpdates() (io.Closer, error) {
	f.checkingShown++
	return fakeCloser{onClose: func() error {
		f.checkingClosed++
		return nil
	}}, nil
}

func (f *fakePrompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	_ = currentVersion
	_ = newVersion
	f.calls++
	return f.confirmResult, nil
}

func (f *fakePrompter) NotifyUpdateFailed(currentVersion string, updateErr error) error {
	f.failed = append(f.failed, currentVersion+"|"+updateErr.Error())
	return nil
}

func (f *fakePrompter) NotifyUpdateSuccess(previousVersion, newVersion string) error {
	f.succeeded = append(f.succeeded, previousVersion+"|"+newVersion)
	return nil
}

func newService(root string, provider fakeReleaseProvider, downloader fakeDownloader, runner *fakeRunner, prompter *fakePrompter) Service {
	return Service{
		InstallRoot:     root,
		PlatformOS:      "linux",
		PlatformArch:    "amd64",
		Logger:          slog.New(slog.NewTextHandler(io.Discard, nil)),
		ReleaseProvider: provider,
		Downloader:      downloader,
		FileStorage:     osfile.Storage{},
		AppRunner:       runner,
		Prompter:        prompter,
		Now: func() time.Time {
			return time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC)
		},
	}
}

func manifestOnlyRelease(t *testing.T, manifest UpdateManifest) (fakeReleaseProvider, fakeDownloader) {
	t.Helper()
	return releaseWithAssets(t, manifest, map[string][]byte{
		UpdateManifestFile: marshalManifest(t, manifest),
	})
}

func releaseWithAssets(t *testing.T, manifest UpdateManifest, assets map[string][]byte) (fakeReleaseProvider, fakeDownloader) {
	t.Helper()
	release := ReleaseMeta{TagName: "v" + manifest.Version}
	for name := range assets {
		release.Assets = append(release.Assets, ReleaseAsset{Name: name, DownloadURL: name})
	}
	return fakeReleaseProvider{release: release}, fakeDownloader(assets)
}

func tarballManifest(t *testing.T, versionName, entrypoint string, content []byte) ([]byte, UpdateManifest) {
	t.Helper()
	archive, err := BuildTarGz(map[string][]byte{entrypoint: content})
	if err != nil {
		t.Fatalf("BuildTarGz() error = %v", err)
	}
	sum := sha256.Sum256(archive)
	return archive, UpdateManifest{
		Version:             versionName,
		Channel:             StableChannel,
		MinSupportedVersion: "0.1.0",
		ReleaseNotes:        "Release " + versionName,
		PublishedAt:         time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
		Assets: []ManifestAsset{{
			OS:         "linux",
			Arch:       "amd64",
			Filename:   "skinprice-linux-amd64.tar.gz",
			Sha256:     hex.EncodeToString(sum[:]),
			Size:       int64(len(archive)),
			Entrypoint: entrypoint,
		}},
	}
}

func zipManifest(t *testing.T, versionName, entrypoint string, content []byte) ([]byte, UpdateManifest) {
	t.Helper()
	archive, err := BuildZip(map[string][]byte{entrypoint: content})
	if err != nil {
		t.Fatalf("BuildZip() error = %v", err)
	}
	sum := sha256.Sum256(archive)
	return archive, UpdateManifest{
		Version:             versionName,
		Channel:             StableChannel,
		MinSupportedVersion: "0.1.0",
		ReleaseNotes:        "Release " + versionName,
		PublishedAt:         time.Date(2026, 5, 20, 12, 0, 0, 0, time.UTC),
		Assets: []ManifestAsset{{
			OS:         "windows",
			Arch:       "amd64",
			Filename:   "skinprice-windows-amd64.zip",
			Sha256:     hex.EncodeToString(sum[:]),
			Size:       int64(len(archive)),
			Entrypoint: entrypoint,
		}},
	}
}

func marshalManifest(t *testing.T, manifest UpdateManifest) []byte {
	t.Helper()
	payload, err := json.Marshal(manifest)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	return payload
}

func writeCurrentState(t *testing.T, root string, dto CurrentVersionDTO) {
	t.Helper()
	payload, err := json.Marshal(dto)
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(filepath.Join(root, CurrentStateFile), payload, 0o644); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}

func readCurrentState(t *testing.T, root string) CurrentVersionDTO {
	t.Helper()
	payload, err := os.ReadFile(filepath.Join(root, CurrentStateFile))
	if err != nil {
		t.Fatalf("ReadFile() error = %v", err)
	}
	var dto CurrentVersionDTO
	if err := json.Unmarshal(payload, &dto); err != nil {
		t.Fatalf("json.Unmarshal() error = %v", err)
	}
	return dto
}

func writeExecutable(t *testing.T, path string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("MkdirAll() error = %v", err)
	}
	if err := os.WriteFile(path, []byte("binary"), 0o755); err != nil {
		t.Fatalf("WriteFile() error = %v", err)
	}
}
