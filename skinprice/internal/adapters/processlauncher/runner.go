package processlauncher

import (
	"os/exec"
	"path/filepath"
)

type Runner struct{}

func (Runner) Start(entrypoint string) error {
	cmd := exec.Command(entrypoint)
	cmd.Dir = filepath.Dir(entrypoint)
	return cmd.Start()
}
