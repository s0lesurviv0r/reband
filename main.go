package main

import (
	"fmt"
	"os"

	"github.com/s0lesurviv0r/reband/cmd"
)

func main() {
	if err := cmd.RootCmd().Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
