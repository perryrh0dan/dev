package image

import (
	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func NewImageCmd(cfg *config.Config, engine container.ContainerEngine) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "image",
		Short: "Manage dev container images",
	}
	cmd.AddCommand(newPullCmd(cfg, engine))
	cmd.AddCommand(newTagsCmd(cfg, engine))
	return cmd
}
