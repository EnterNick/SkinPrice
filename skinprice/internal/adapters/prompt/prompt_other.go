//go:build !linux && !windows

package prompt

import (
	"fmt"
	"io"
)

func (Prompter) ShowCheckingForUpdates() (io.Closer, error) {
	return nil, nil
}

func (Prompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	return false, fmt.Errorf("update prompt is not supported on this platform")
}

func (Prompter) NotifyUpdateFailed(currentVersion string, updateErr error) error {
	return nil
}

func (Prompter) NotifyUpdateSuccess(previousVersion, newVersion string) error {
	return nil
}
