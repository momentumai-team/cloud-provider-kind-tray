package actions

import (
	"io"
	"os/exec"
)

func Start() (*exec.Cmd, io.ReadCloser, io.ReadCloser, error) {
	// Define the command to run
	cmd := exec.Command("sudo", "cloud-provider-kind", "-enable-load-balancer-status")

	// Get a pipe for the command's stdout
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, nil, nil, err
	}

	// Get a pipe for the command's stderr
	stderr, err := cmd.StderrPipe()
	if err != nil {
		return nil, nil, nil, err
	}
	// Start the command
	if err := cmd.Start(); err != nil {
		return nil, nil, nil, err
	}
	return cmd, stdout, stderr, nil
}
