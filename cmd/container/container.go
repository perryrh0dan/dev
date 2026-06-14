package container

import (
	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func NewContainerCmd(cfg *config.Config, engine container.ContainerEngine) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "container",
		Short: "Manage dev containers and volumes",
	}
	cmd.AddCommand(newCreateCmd(engine))
	cmd.AddCommand(newListCmd(engine))
	cmd.AddCommand(newRemoveCmd(engine))
	cmd.AddCommand(newStopCmd(cfg, engine))
	cmd.AddCommand(newAttachCmd(engine))
	cmd.AddCommand(newBackupCmd(engine))
	cmd.AddCommand(newRestoreCmd(engine))
	cmd.AddCommand(newCopyCmd())
	cmd.AddCommand(newShowCmd())
	return cmd
}
