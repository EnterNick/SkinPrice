//go:build windows

package prompt

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"
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
	script := `Add-Type -AssemblyName System.Windows.Forms; ` +
		`$form = New-Object System.Windows.Forms.Form; ` +
		`$form.Text = 'SkinPrice'; ` +
		`$form.StartPosition = 'CenterScreen'; ` +
		`$form.TopMost = $true; ` +
		`$form.Width = 340; ` +
		`$form.Height = 120; ` +
		`$form.FormBorderStyle = 'FixedDialog'; ` +
		`$form.ControlBox = $false; ` +
		`$label = New-Object System.Windows.Forms.Label; ` +
		`$label.AutoSize = $false; ` +
		`$label.TextAlign = 'MiddleCenter'; ` +
		`$label.Dock = 'Fill'; ` +
		`$label.Text = 'Checking for updates...'; ` +
		`$form.Controls.Add($label); ` +
		`[System.Windows.Forms.Application]::Run($form)`
	cmd := exec.CommandContext(context.Background(), "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return processCloser{cmd: cmd}, nil
}

func (Prompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	script := fmt.Sprintf(`Add-Type -AssemblyName PresentationFramework; $result = [System.Windows.MessageBox]::Show('Update SkinPrice from %s to %s?','SkinPrice Update','YesNo','Question'); if ($result -eq 'Yes') { exit 0 } else { exit 1 }`, currentVersion, newVersion)
	cmd := exec.CommandContext(context.Background(), "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
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
	message := fmt.Sprintf("Failed to update SkinPrice. Running current version %s.`n%s", currentVersion, updateErr.Error())
	return runNotification("SkinPrice Update Failed", message)
}

func (Prompter) NotifyUpdateSuccess(previousVersion, newVersion string) error {
	message := fmt.Sprintf("SkinPrice updated from %s to %s.", previousVersion, newVersion)
	return runNotification("SkinPrice Updated", message)
}

func runNotification(title, message string) error {
	script := fmt.Sprintf(`Add-Type -AssemblyName System.Windows.Forms; Add-Type -AssemblyName System.Drawing; $notify = New-Object System.Windows.Forms.NotifyIcon; $notify.Icon = [System.Drawing.SystemIcons]::Information; $notify.BalloonTipTitle = '%s'; $notify.BalloonTipText = '%s'; $notify.Visible = $true; $notify.ShowBalloonTip(5000); Start-Sleep -Seconds 6; $notify.Dispose()`, escapePowerShellSingleQuoted(title), escapePowerShellSingleQuoted(message))
	cmd := exec.CommandContext(context.Background(), "powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	return cmd.Run()
}

func escapePowerShellSingleQuoted(value string) string {
	return strings.ReplaceAll(value, `'`, `''`)
}
