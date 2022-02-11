package main

import (
	"fmt"
	"os"
	"runtime"

	"github.com/homepi/homepi/src/cmds"
	"github.com/homepi/homepi/src/core"
	"github.com/spf13/cobra"
)

var (
	BranchName string
	Version    string
	CompiledBy string
	BuildTime  string
)

func main() {

	rootCmd := &cobra.Command{
		Use:     "homepi",
		Version: Version,
		Long: `
    __  __                     ____  _ 
   / / / /___  ________  ___  / __ \(_)
  / /_/ / __ \/ __  __ \/ _ \/ /_/ / / 
 / __  / /_/ / / / / / /  __/ ____/ /  
/_/ /_/\____/_/ /_/ /_/\___/_/   /_/`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	rootCmd.SetArgs(os.Args[1:])

	vi := &core.VersionInfo{
		Version:    Version,
		BranchName: BranchName,
		CompiledBy: CompiledBy,
		GoVersion:  runtime.Version(),
		BuildTime:  BuildTime,
	}
	if err := cmds.RegisterAndRun(vi, rootCmd); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
