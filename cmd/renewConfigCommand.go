package cmd

import (
	"edpvp-log-prepare/config"
	"fmt"
	"github.com/fatih/color"
)

func RenewConfigCommand() {
	fmt.Println("\n\n\n-----------------\nThis is a WIP. If you wish to reset the Config, go to ")
	cfgPath := config.BuildConfigDirIfNotExistsAndReturnDir()
	color.Yellow(cfgPath)
	fmt.Println("and delete the config.json")
}
