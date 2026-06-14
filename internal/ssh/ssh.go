package ssh

import (
	"fmt"
	"os"
	"os/exec"
)

// ConnectAll iterates over hosts and SSH-connects to each sequentially.
func ConnectAll(hosts []string) error {
	if len(hosts) == 0 {
		return fmt.Errorf("no remote hosts configured (set remote_dev_env in ~/.config/dev/config.json)")
	}
	var lastErr error
	for _, host := range hosts {
		if err := connect(host); err != nil {
			fmt.Fprintf(os.Stderr, "ssh to %s exited: %v\n", host, err)
			lastErr = err
		}
	}
	return lastErr
}

func connect(host string) error {
	args := buildSSHArgs(host)
	cmd := exec.Command("ssh", args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func buildSSHArgs(host string) []string {
	return []string{host}
}
