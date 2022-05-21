package config

import (
	"fmt"
	"strings"
)

type AppConfig struct {
	EliteLogPath           string   `json:"eliteLogPath,omitempty"`
	CmdrInclude            []string `json:"cmdrInclude,omitempty"`
	OutputDir              string   `json:"outputDir,omitempty"`
	LastExecutionTimeStamp int64    `json:"lastExecutionTimeStamp,omitempty"`
}

func ConfigToPrintString(config AppConfig) string {

	var whiteListString string
	if len(config.CmdrInclude) == 0 {
		whiteListString = "Not Used"
	} else {
		whiteListString = strings.Join(config.CmdrInclude, ", ")
	}
	return fmt.Sprintf(":::Config::: \n"+
		"Elite Logs Path: \n"+
		"\t%s\n"+
		"Whitelist:\n"+
		"\t%s\n"+
		"ZIP Output Directory:\n"+
		"\t%s\n"+
		"Last Execution TimeStamp:\n"+
		"\t%d",
		config.EliteLogPath,
		whiteListString,
		config.OutputDir,
		config.LastExecutionTimeStamp)
}
