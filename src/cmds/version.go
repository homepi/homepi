package cmds

import (
	"fmt"

	"github.com/spf13/cobra"
)

func versionCommand() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print current version of homepi",
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Println(fmt.Sprintf("%s version %s", cmd.Parent().Name(), cmd.Root().Version))
		},
	}
}
