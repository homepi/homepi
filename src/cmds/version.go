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
			fmt.Fprint(
				cmd.OutOrStdout(),
				fmt.Sprintf(
					"%s\n%s\n%s\n%s\n%s",
					fmt.Sprintf("HomePi version: %s", versionInfo.Version),
					fmt.Sprintf("Build type: %s", versionInfo.BuildType),
					fmt.Sprintf("Build time: %s", versionInfo.BuildTime),
					fmt.Sprintf("Golang: %s", versionInfo.GoVersion),
					fmt.Sprintf("Compiled by: %s", versionInfo.CompiledBy),
				),
			)
		},
	}
}
