package cmd

import (
	"bufio"
	"edpvp-log-prepare/config"
	"edpvp-log-prepare/util"
	"fmt"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func doesDefaultEliteLogLocationExist() (bool, string) {
	homeDir, _ := os.UserHomeDir()
	path := filepath.Join(
		homeDir,
		"Saved Games",
		"Frontier Developments",
		"Elite Dangerous",
	)
	return checkIfDirectoryForEliteLogsExists(path), path

}

func checkIfDirectoryForEliteLogsExists(path string) bool {
	data, err := os.Stat(path)
	if err != nil {
		fmt.Printf("Tried to find the Elite Dangerous logs at ")
		color.Set(color.FgYellow)
		fmt.Print(path)
		color.Unset()
		fmt.Print(", but that Folder does not exist...\n")
		return false
	}
	if !data.IsDir() {
		color.White("Tried to find the Elite Dangerous logs at '%s', but that directory is not a folder", path)
		return false
	}
	return true
}

const ConfigLen = 2

func getNewLogFileLocation(reader *bufio.Reader) string {
	exists, folderPath := doesDefaultEliteLogLocationExist()
	if exists {
		fmt.Print("Found a Folder at ")
		color.Set(color.FgYellow)
		fmt.Print(folderPath)
		color.Unset()
		fmt.Println(" . Should this be where I look for Elite logs?")
		color.Green("Y/N")
		if util.GetYesNo(reader) {
			return folderPath
		}
		exists = false
	}

	for !exists {
		color.Green("Please Enter the Path to the Elite Log Directory...")
		reader.Discard(reader.Buffered())
		userInput, err := reader.ReadString('\n')
		userInput = strings.Trim(userInput, "\r\n")
		if err != nil {
			if err == io.EOF {
				continue // Wait again
			}
			if err.Error() == "unexpected newline" {
				continue // Try again
			}
			panic(err)
		}
		existsUserInput := checkIfDirectoryForEliteLogsExists(userInput)
		if existsUserInput {
			// Check if any log files are present. If not, ask the user if they really want to use this folder
			files, err := ioutil.ReadDir(userInput)
			if err != nil {
				panic(err)
			}
			doesContainLogFiles := false
			for _, file := range files {
				if file.IsDir() {
					continue
				}
				if strings.HasSuffix(file.Name(), ".log") {
					doesContainLogFiles = true
					break
				}
			}
			if !doesContainLogFiles {
				color.Red("The Folder you provided does no contain any .log files. Are you sure this is the" +
					"correct folder?")
				color.Green("Y/N")
				if !util.GetYesNo(reader) {
					continue
				}
			}
			exists = true
			folderPath = userInput
		}
	}
	return folderPath
}

func BuildNewConfigFromScratch() (config.AppConfig, error) {
	reader := bufio.NewReader(os.Stdin)
	cfg := config.AppConfig{}
	// Logs Location
	color.Blue("\n\n\n[1/%d] === Elite Dangerous Logs Folder", ConfigLen)
	cfg.EliteLogPath = getNewLogFileLocation(reader)

	color.Blue("\n\n\n[2/%d] === Output Directory", ConfigLen)
	cfg.OutputDir = getNewOutputDir(cfg.EliteLogPath, reader)

	color.Blue("\n\n\n[3/%d] === Whitelist CMDRs", ConfigLen)
	color.White("This tool can make sure that only Log Files from specific CMDRs are " +
		"Zipped to be sent to the bot. Do you want to filter specific CMDRs? ")
	color.Green("Y/N")
	if util.GetYesNo(reader) {
		cfg.CmdrInclude = getCommanderWhitelist(reader)
	} else {
		cfg.CmdrInclude = make([]string, 0)
	}

	color.Blue("\n\n\n[4/%d] === Set Last Logfile's Timestamp", ConfigLen)
	color.White("The tool will collect all logs from the last time you started the tool.\n" +
		"As this is the first time you run this tool you can set a Date. Logs that have an older Timestamp in their Name" +
		"as the given Date will be considered.\nPlease provide the Timestamp you wish to set. If you dont know what a " +
		"Timestamp is, go to https://www.unixtimestamp.com/, enter the desired date, and paste the number it throws back at you.\n" +
		"If you want to ZIP all previous logs, you can just enter '0' ")
	cfg.LastExecutionTimeStamp = getLastExecTimeStamp(reader)

	color.White("\n\n%s\n\nIf you wish to make changes, call this program with 'config' as the first argument",
		config.ConfigToPrintString(cfg))

	return cfg, nil
}

func getLastExecTimeStamp(reader *bufio.Reader) int64 {
	for true {
		color.Green("Please enter a Unix Timestamp...")
		var timeStamp int64
		reader.Discard(reader.Buffered())
		var timeStampString, err = reader.ReadString('\n')
		timeStampString = strings.Trim(timeStampString, "\r\n")

		if err != nil {
			if err.Error() == "unexpected newline" {
				// try again
				continue
			}
			panic(err)
		}
		timeStamp, err = strconv.ParseInt(timeStampString, 10, 64)
		if err != nil {
			color.Red("Failed to parse input as Int")
			color.Red(err.Error())
			continue
		}

		if timeStamp < 0 {
			color.Red("Negative Timestamps are not supported. Retry")
			continue
		}

		return timeStamp
	}
	panic("Should never happen")
}

func getCommanderWhitelist(reader *bufio.Reader) []string {
	for true {
		color.Green("Enter a list of CMDRs to include in the ZIP, seperated by commas")
		color.White("Example:\nWDX, SomeOtherCmdr, SomeoneElse")
		reader.Discard(reader.Buffered())
		var userInput, err = reader.ReadString('\n')
		userInput = strings.Trim(userInput, "\r\n")

		if err != nil {
			if err.Error() == "unexpected newline" {
				userInput = ""
			} else {
				panic(err)
			}
		}
		cmdrs := strings.Split(strings.ToUpper(userInput), ",")

		if len(cmdrs) == 0 {
			color.Yellow("No CMDRs entered. This will disable filtering. Is this okay?")
			color.Green("Y/N")
			if util.GetYesNo(reader) {
				return make([]string, 0)
			} else {
				continue
			}
		}

		color.White("Only logs of the following CMDRs will be considered:")
		for i, cmdr := range cmdrs {
			cmdrs[i] = strings.Trim(cmdr, " ")
			color.White("\tCMDR %s", cmdr)
		}

		color.Green("Is this correct? Y/N")
		if util.GetYesNo(reader) {
			return cmdrs
		}
	}

	panic("Should never happen")
}

func getNewOutputDir(logPath string, reader *bufio.Reader) string {
	returnFilePath := filepath.Join(logPath, "#edpvp")
	color.White("This defines where the output ZIP files are located. This defaults to logLocation/#edpvp, or in" +
		" your instance: \n")
	color.Yellow(returnFilePath)
	color.Green("Is the default path okay? Y/N")
	isInputOkay := util.GetYesNo(reader)
	for !isInputOkay {
		color.Green("Please enter where you want the output ZIPs to be placed")
		reader.Discard(reader.Buffered())
		userInput, err := reader.ReadString('\n')
		userInput = strings.Trim(userInput, "\r\n")

		if err != nil {
			if err.Error() == "unexpected newline" {
				continue // Try again
			}
			panic(err)
		}
		userInput = strings.Trim(userInput, " ")
		userInputData, err := os.Stat(userInput)
		if err != nil {
			if os.IsNotExist(err) {
				color.Yellow("%s does not exists. The program will try to create that directory "+
					"when saving files. Is that okay?", userInput)
				color.Green("Y/N")
				if util.GetYesNo(reader) {
					returnFilePath = userInput
					isInputOkay = true
				}
			} else {
				color.Red("Failed to get Path. Try again. Error: %s", err.Error())
			}
		} else if !userInputData.IsDir() {
			color.Red("The provided Path is not a directory. Try again")
		} else {
			color.White("Found Directory %s", userInputData)
			color.Green("Set this Directory as Output Directory? Y/N")
			if util.GetYesNo(reader) {
				returnFilePath = userInput
				isInputOkay = true
			}
		}

	}

	return returnFilePath
}
