package main

import (
	"edpvp-log-prepare/cmd"
	"github.com/fatih/color"
	"github.com/teamseodo/cli"
	"log"
	"os"
)

func main() {
	mainCommand := cli.NewCommand("app", "app", "The Default Command – "+
		"Use it to build the ZIP files")
	mainCommand.Do(cmd.MainCommand)

	help, err := mainCommand.FindHelp(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}
	help.ShowHelp()

	command, err := mainCommand.Parse(os.Args[1:])
	if err != nil {
		log.Fatal(err)
	}

	color.Blue("Running edpvp-log-prepare v%s\n", VERSION)
	err = command.Run()
	if err != nil {
		command.Help().ShowHelp()
	}
	CheckUpdate()
}
