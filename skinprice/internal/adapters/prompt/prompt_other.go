//go:build !linux && !windows

package prompt

import "fmt"

func (Prompter) ConfirmUpdate(currentVersion, newVersion string) (bool, error) {
	return false, fmt.Errorf("update prompt is not supported on this platform")
}
