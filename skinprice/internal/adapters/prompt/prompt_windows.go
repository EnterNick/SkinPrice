//go:build windows

package prompt

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

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
