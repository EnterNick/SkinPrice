//go:build linux

package prompt

import (
	"fmt"
	"os/exec"
)

func (Prompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	cmd := exec.Command(
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
	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return false, nil
	}
	return false, err
}
