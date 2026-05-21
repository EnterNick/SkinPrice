//go:build linux

package prompt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
)

type processCloser struct {
	cmd *exec.Cmd
}

func (p processCloser) Close() error {
	if p.cmd == nil || p.cmd.Process == nil {
		return nil
	}
	if err := p.cmd.Process.Kill(); err != nil && !errors.Is(err, os.ErrProcessDone) {
		return err
	}
	_, _ = p.cmd.Process.Wait()
	return nil
}

func (Prompter) ShowCheckingForUpdates() (io.Closer, error) {
	if _, err := exec.LookPath("zenity"); err != nil {
		return nil, fmt.Errorf("zenity is not installed")
	}
	cmd := exec.CommandContext(
		context.Background(),
		"zenity",
		"--info",
		"--title=SkinPrice",
		"--text=Checking for updates...",
		"--no-wrap",
		"--width=320",
	)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return processCloser{cmd: cmd}, nil
}

func (Prompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	cmd := exec.CommandContext(
		context.Background(),
		"zenity",
		"--question",
		"--title=SkinPrice Update",
		"--text",
		fmt.Sprintf("Update SkinPrice from %s to %s?", currentVersion, newVersion),
	)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}

func (Prompter) NotifyUpdateFailed(currentVersion string, updateErr error) error {
	message := fmt.Sprintf("Failed to update SkinPrice. Running current version %s.\n%s", currentVersion, updateErr.Error())
	return runNotification("SkinPrice Update Failed", message)
}

func (Prompter) NotifyUpdateSuccess(previousVersion, newVersion string) error {
	message := fmt.Sprintf("SkinPrice updated from %s to %s.", previousVersion, newVersion)
	return runNotification("SkinPrice Updated", message)
}

func runNotification(title, message string) error {
	if _, err := exec.LookPath("notify-send"); err == nil {
		cmd := exec.CommandContext(context.Background(), "notify-send", title, message)
		return cmd.Run()
	}
	if _, err := exec.LookPath("zenity"); err == nil {
		cmd := exec.CommandContext(context.Background(), "zenity", "--notification", "--title", title, "--text", message)
		return cmd.Run()
	}
	return fmt.Errorf("no supported notification command found")
}
