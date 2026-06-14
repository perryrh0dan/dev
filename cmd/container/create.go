package container

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newCreateCmd(engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "create <name>",
		Short: "Create a named dev volume",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]
			if err := engine.CreateVolume(name); err != nil {
				return fmt.Errorf("create volume %q: %w", name, err)
			}
			fmt.Printf("Created volume %q\n", name)
			return nil
		},
	}
}
