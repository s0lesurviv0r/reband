package main

import (
	"github.com/s0lesurviv0r/channel-conv/cmd"
)

func main() {
	err := cmd.RootCmd().Execute()
	if err != nil {
		panic(err)
	}
}
