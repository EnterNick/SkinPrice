//go:build windows

package prompt

import (
	"fmt"
	"os/exec"
)

func (Prompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	script := fmt.Sprintf(`Add-Type -AssemblyName PresentationFramework; $result = [System.Windows.MessageBox]::Show('Update SkinPrice from %s to %s?','SkinPrice Update','YesNo','Question'); if ($result -eq 'Yes') { exit 0 } else { exit 1 }`, currentVersion, newVersion)
	cmd := exec.Command("powershell", "-NoProfile", "-NonInteractive", "-Command", script)
	err := cmd.Run()
	if err == nil {
		return true, nil
	}
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}
