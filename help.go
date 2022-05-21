package main

import (
	"fmt"
	"github.com/fatih/color"
)

func PrintHelp(appName string) {
	color.Blue("Running ED PVP Log Prepare v%s\n\n", VERSION)
	color.Yellow("Commands")
	color.Set(color.FgYellow)
	fmt.Print("    " + appName + " ")
	color.Set(color.FgWhite)
	fmt.Print(": The default command. Run this to bundle the recent files\n\n")
	color.Set(color.FgYellow)
	fmt.Print("    " + appName + " ")
	color.Set(color.FgBlue)
	fmt.Print("config ")
	color.Set(color.FgWhite)
	fmt.Print(": The config Command. Run this to change the Config\n\n")
}
