package processlauncher

import (
	"context"
	"os/exec"
	"path/filepath"
)

type Runner struct{}

func (Runner) Start(entrypoint string) error {
	cmd := exec.CommandContext(context.Background(), entrypoint)
	cmd.Dir = filepath.Dir(entrypoint)
	return cmd.Start()
}
