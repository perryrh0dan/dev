package container

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/spf13/cobra"
)

func newCopyCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "copy <container-name> <path>",
		Short: "Copy stdin into a file inside a running container",
		Long:  "Reads from stdin and writes it to <path> inside the named container. Example: pbpaste | dev container copy myenv /root/file.txt",
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			containerName := args[0]
			targetPath := args[1]

			c := exec.Command("docker", "exec", "-i", containerName,
				"sh", "-c", `cat > "$1"`, "--", targetPath)
			c.Stdin = os.Stdin
			c.Stdout = os.Stdout
			c.Stderr = os.Stderr
			if err := c.Run(); err != nil {
				return fmt.Errorf("copy to container %q at %q: %w", containerName, targetPath, err)
			}
			return nil
		},
	}
}
