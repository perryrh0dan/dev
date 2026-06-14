package container

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newAttachCmd(engine container.ContainerEngine) *cobra.Command {
	var shell string

	cmd := &cobra.Command{
		Use:   "attach <container-id>",
		Short: "Attach to a running dev container",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			id := args[0]
			if err := engine.Exec(id, []string{shell}); err != nil {
				return fmt.Errorf("attach to %q: %w", id, err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&shell, "shell", "/bin/zsh", "Shell to launch in the container")
	return cmd
}
