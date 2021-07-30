package cmds

import (
	"github.com/spf13/cobra"
)

//TODO: Create delete commands
func deleteCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "delete",
		Short: "Delete homepi resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(&cobra.Command{
		Use:   "user",
		Short: "Delete a user",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "accessory",
		Short: "Delete an accessory",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
	cmd.AddCommand(&cobra.Command{
		Use:   "role",
		Short: "Delete a role",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	})
	return cmd
}
