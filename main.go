package main

import (
	"edpvp-log-prepare/cmd"
	"os"
	"path/filepath"
)

func main() {
	_, executingFileName := filepath.Split(os.Args[0])

	if len(os.Args) == 1 {
		cmd.MainCommand()
	} else if len(os.Args) > 1 {
		if os.Args[1] == "config" {
			cmd.RenewConfigCommand()
		} else {
			PrintHelp(executingFileName)
		}
	} else {
		PrintHelp(executingFileName)
	}

	CheckUpdate()
}
