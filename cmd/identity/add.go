package identity

import (
	"fmt"

	"github.com/perryrh0dan/dev/internal/config"
	"github.com/spf13/cobra"
)

func newAddCmd() *cobra.Command {
	var email, name, keyid string

	cmd := &cobra.Command{
		Use:   "add",
		Short: "Add or update a git identity and make it active",
		RunE: func(cmd *cobra.Command, args []string) error {
			id := config.Identity{Email: email, Name: name, KeyID: keyid}
			if err := config.SaveIdentity(id); err != nil {
				return fmt.Errorf("save identity: %w", err)
			}
			fmt.Printf("Identity %q (%s) saved and activated\n", name, email)
			return nil
		},
	}

	cmd.Flags().StringVar(&email, "email", "", "Git email address (required)")
	cmd.Flags().StringVar(&name, "name", "", "Git display name (required)")
	cmd.Flags().StringVar(&keyid, "keyid", "", "GPG key ID (required)")
	_ = cmd.MarkFlagRequired("email")
	_ = cmd.MarkFlagRequired("name")
	_ = cmd.MarkFlagRequired("keyid")

	return cmd
}
