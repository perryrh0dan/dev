package identity

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/config"
	"github.com/spf13/cobra"
)

func newActivateCmd() *cobra.Command {
	var email string

	cmd := &cobra.Command{
		Use:   "activate",
		Short: "Activate an existing git identity by email",
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := config.ActivateIdentity(email); err != nil {
				return fmt.Errorf("activate identity: %w", err)
			}
			fmt.Printf("Identity %q is now active\n", email)
			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Email of the identity to activate (required)")
	_ = cmd.MarkFlagRequired("email")

	return cmd
}
