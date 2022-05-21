package cmd

import (
	"bufio"
	"edpvp-log-prepare/config"
	"edpvp-log-prepare/util"
	"edpvp-log-prepare/zipper"
	"fmt"
	"github.com/fatih/color"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)
import "github.com/skratchdot/open-golang/open"

func aggregateLogFiles(config config.AppConfig) []string {
	// Get all files in Log Directory
	files, err := ioutil.ReadDir(config.EliteLogPath)
	if err != nil {
		panic(err)
	}

	filterCmdrs := len(config.CmdrInclude) > 0
	relevantFiles := make([]string, 0)

	for _, file := range files {
		if file.IsDir() {
			continue
		}
		fileName := file.Name()

		if !strings.HasSuffix(fileName, ".log") {
			continue
		}
		fileMeta := util.HandleFile(fileName)

		// This is a log file. Check if file matches either Horizons or Odyssey format
		if fileMeta.Type == util.UNKNOWN {
			continue
		}
		if fileMeta.TimeStamp < config.LastExecutionTimeStamp {
			continue
		}
		if filterCmdrs {
			// check the log header. check if the CMDR is on the whitelist
			cmdr, success := util.GetCommanderInLog(filepath.Join(config.EliteLogPath, fileName))
			if !success {
				continue
			}
			isCmdrInWhitelist := false
			cmdrUpper := strings.ToUpper(cmdr)
			for _, c := range config.CmdrInclude {
				if cmdrUpper == c {
					isCmdrInWhitelist = true
					break
				}
			}
			if !isCmdrInWhitelist {
				continue
			}
		}
		relevantFiles = append(relevantFiles, filepath.Join(config.EliteLogPath, fileName))
	}
	return relevantFiles
}

func MainCommand() {
	hasConfig, cfg := config.GetConfig()
	if !hasConfig {
		color.Yellow("This appears to be your first time running this program. \n" +
			"This will run you through a Configuration of the App")

		cfg, _ = BuildNewConfigFromScratch()
		color.White("Saving config...")
		err := config.SetConfig(cfg)
		if err != nil {
			panic(err)
		}
	}

	relevantLogFiles := aggregateLogFiles(cfg)

	if len(relevantLogFiles) == 0 {
		color.Red("No relevant Log files were found. No ZIP has been generated...")
		color.Green("Press any key to continue...")
		reader := bufio.NewReader(os.Stdin)
		reader.Discard(reader.Buffered())
		_, _ = reader.ReadString('\n')
		return
	}

	// Now that the relevant files are aggregated... create a Temporary Directory and copy over all relevant files
	// This is done because maybe some file is holding the original file and renaming it would cause issues, etc
	tempDir, err := ioutil.TempDir("", "edpvp-log-prepare-*")
	if err != nil {
		panic(err)
	}
	defer os.Remove(tempDir)
	color.White(tempDir)
	moveRelevantFilesToTempDir(relevantLogFiles, tempDir)

	directoryPath := zipper.ZipUpFiles(tempDir, cfg)
	cfg.LastExecutionTimeStamp = time.Now().Unix()
	err = config.SetConfig(cfg)
	if err != nil {
		panic(err)
	}
	color.Blue("\n\nDone building ZIP files. Do you want to open the directory?")
	color.Blue(directoryPath)
	color.Green("Y/N")

	if util.GetYesNo(bufio.NewReader(os.Stdin)) {
		err := open.Run(directoryPath)
		if err != nil {
			panic(err)
		}
	}

}

func moveRelevantFilesToTempDir(logFiles []string, tempDir string) {
	for _, file := range logFiles {
		_, justFileName := filepath.Split(file)
		newFileName := util.ConvertOddyToHorizonsName(justFileName)
		fileContent, err := ioutil.ReadFile(file)
		if err != nil {
			fmt.Println(err)
			continue
		}

		err = ioutil.WriteFile(filepath.Join(tempDir, newFileName), fileContent, 0644)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
