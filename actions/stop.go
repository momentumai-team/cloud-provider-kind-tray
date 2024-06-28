package actions

import (
	"os/exec"
	"syscall"
)

func Stop(cmd *exec.Cmd) error {
	// Stop the process
	if err := cmd.Process.Signal(syscall.SIGTERM); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
