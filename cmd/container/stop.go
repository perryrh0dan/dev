package container

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/config"
	"github.com/perryrh0dan/dev/internal/container"
	"github.com/spf13/cobra"
)

func newStopCmd(cfg *config.Config, engine container.ContainerEngine) *cobra.Command {
	return &cobra.Command{
		Use:   "stop",
		Short: "Stop all running dev containers",
		RunE: func(cmd *cobra.Command, args []string) error {
			if cfg.Image == "" {
				return fmt.Errorf("image not configured (set image in ~/.config/dev/config.json or DOCKER_DEV_ENV)")
			}
			if err := engine.StopContainers(cfg.Image); err != nil {
				return fmt.Errorf("stop containers: %w", err)
			}
			fmt.Println("Stopped all dev containers")
			return nil
		},
	}
}
