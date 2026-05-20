package logx

import (
	"path/filepath"
	"testing"
)

func TestResolveLogPathFromDirectory(t *testing.T) {
	dir := t.TempDir()

	got, err := resolveLogPath(Config{
		FilePath: dir + "/",
		AppName:  "SkinPrice",
	})
	if err != nil {
		t.Fatalf("resolveLogPath() error = %v", err)
	}

	want := filepath.Join(dir, "skinprice.log")
	if got != want {
		t.Fatalf("resolveLogPath() = %s, want %s", got, want)
	}
}

func TestResolveLogPathFromExistingDirectory(t *testing.T) {
	dir := t.TempDir()

	got, err := resolveLogPath(Config{
		FilePath: dir,
		AppName:  "SkinPrice",
	})
	if err != nil {
		t.Fatalf("resolveLogPath() error = %v", err)
	}

	want := filepath.Join(dir, "skinprice.log")
	if got != want {
		t.Fatalf("resolveLogPath() = %s, want %s", got, want)
	}
}
