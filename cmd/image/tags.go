package image

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newTagsCmd(cfg *config.Config, engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "tags",
		Short: "List locally available tags for the dev image",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Image == "" {
				return fmt.Errorf("image not configured (set image in ~/.config/dev/config.json or DOCKER_DEV_ENV)")
			}
			tags, err := engine.ListImageTags(cfg.Image)
			if err != nil {
				return fmt.Errorf("list tags for %s: %w", cfg.Image, err)
			}
			for _, t := range tags {
				fmt.Println(t)
			}
			return nil
		},
	}
}
