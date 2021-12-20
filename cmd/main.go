/*
Copyright © 2021 Max Brauer <mamachanko>

*/
package main

import (
	"github.com/mamachanko/need/pkg/cmd"
	"os"
)

func main() {
	command := cmd.NewNeedCmd()

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}
