package main

import (
	"fmt"
	"os"

	"github.com/homepi/homepi/src/cmds"
)

func main() {
	if err := cmds.Run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
