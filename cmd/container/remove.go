package container

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newRemoveCmd(engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "remove <name>",
		Short: "Remove a dev environment and its companion volumes",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			name := args[0]

			// Primary volume is required — fail hard if missing
			if err := engine.RemoveVolumes(name); err != nil {
				return fmt.Errorf("remove volume %q: %w", name, err)
			}

			// Companion volumes may not exist — warn but don't fail
			for _, companion := range []string{name + "-history", name + "-resurrect"} {
				if err := engine.RemoveVolumes(companion); err != nil {
					fmt.Fprintf(cmd.ErrOrStderr(), "warning: could not remove %q: %v\n", companion, err)
				}
			}

			fmt.Fprintf(cmd.OutOrStdout(), "Removed dev environment %q\n", name)
			return nil
		},
	}
}
