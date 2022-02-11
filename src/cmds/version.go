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
			cmd.Println(fmt.Sprintf("HomePi version: %s", versionInfo.Version))
			cmd.Println(fmt.Sprintf("Build type: %s", versionInfo.BuildType))
			cmd.Println(fmt.Sprintf("Build time: %s", versionInfo.BuildTime))
			cmd.Println(fmt.Sprintf("Golang: %s", versionInfo.GoVersion))
			cmd.Println(fmt.Sprintf("Compiled by: %s", versionInfo.CompiledBy))
		},
	}
}
