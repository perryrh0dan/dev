package image

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newPullCmd(cfg *config.Config, engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "pull",
		Short: "Pull the latest dev image and all tagged variants in use",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Image == "" {
				return fmt.Errorf("image not configured (set image in ~/.config/dev/config.json or DOCKER_DEV_ENV)")
			}

			fmt.Printf("Pulling %s...\n", cfg.Image)
			if err := engine.PullImage(cfg.Image); err != nil {
				return fmt.Errorf("pull %s: %w", cfg.Image, err)
			}

			envs, err := config.LoadEnvironments()
			if err != nil {
				return fmt.Errorf("load environments: %w", err)
			}

			seen := map[string]bool{}
			for _, env := range envs {
				if env.Tag == "" || seen[env.Tag] {
					continue
				}
				seen[env.Tag] = true
				tagged := cfg.Image + ":" + env.Tag
				fmt.Printf("Pulling %s...\n", tagged)
				if err := engine.PullImage(tagged); err != nil {
					return fmt.Errorf("pull %s: %w", tagged, err)
				}
			}
			return nil
		},
	}
}
