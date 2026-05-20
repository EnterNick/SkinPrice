//go:build linux

package prompt

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
)

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
