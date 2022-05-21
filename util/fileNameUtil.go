package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type EliteLogFileName struct {
	OriginalName string
	Type         EliteLogFileType
	TimeStamp    int64
}

type EliteLogFileType byte

const (
	ERRORED  EliteLogFileType = 0xFF
	UNKNOWN  EliteLogFileType = 0x00
	HORIZONS EliteLogFileType = 0x01
	ODYSSEY  EliteLogFileType = 0x02
)

func HandleFile(fileName string) EliteLogFileName {
	returnVal := EliteLogFileName{
		OriginalName: fileName,
	}

	returnVal.Type = getLogFileType(fileName)
	if returnVal.Type == UNKNOWN {
		return returnVal
	}

	// try to extract Timestamp
	returnVal.TimeStamp = getTimeStamp(fileName, returnVal.Type)
	return returnVal
}

func GetCommanderInLog(pathName string) (string, bool) {
	file, err := os.Open(pathName)
	if err != nil {
		fmt.Println(err)
		return "", false
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()
		lineAsStr := string(line)
		if strings.Contains(lineAsStr, "\"event\":\"Commander\"") {
			type RelevantData struct {
				Event string `json:"event,omitempty"`
				Name  string `json:"Name,omitempty"`
			}
			var data RelevantData
			err = json.Unmarshal(line, &data)
			if err != nil || data.Name == "" {
				fmt.Println(err)
				return "", false
			}

			return data.Name, data.Event == "Commander"
		}
	}
	return "", false
}

func ConvertOddyToHorizonsName(name string) string {
	// if already horizons, return as is
	if isHorizonsFile(name) {
		return name
	}

	// 2022-05-20T170535
	// YYYY-MM-DDThhmmss
	splitString := strings.Split(name, ".")
	timeStampOddy := splitString[1]

	// 2022-05-20  170535
	dateAndTimeSplit := strings.Split(timeStampOddy, "T")
	//2022  05  20
	dates := strings.Split(dateAndTimeSplit[0], "-")
	// on the first one, remove the first two chars. turn 2022 into 22
	dates[0] = dates[0][2:]

	asHorizonsTimeStamp := strings.Join(dates, "") + dateAndTimeSplit[1]
	splitString[1] = asHorizonsTimeStamp
	return strings.Join(splitString, ".")
}

func getTimeStamp(name string, fileType EliteLogFileType) int64 {

	if fileType == ODYSSEY {
		name = ConvertOddyToHorizonsName(name)
	}

	timeCodePart := strings.Split(name, ".")[1]

	year, _ := strconv.Atoi("20" + timeCodePart[0:2])
	month, _ := strconv.Atoi(timeCodePart[2:4])
	day, _ := strconv.Atoi(timeCodePart[4:6])
	hr, _ := strconv.Atoi(timeCodePart[6:8])
	min, _ := strconv.Atoi(timeCodePart[8:10])
	sec, _ := strconv.Atoi(timeCodePart[10:12])

	date := time.Date(year, time.Month(month), day, hr, min, sec, 0, time.UTC)

	return date.Unix()

}

func getLogFileType(fileName string) EliteLogFileType {
	if isHorizonsFile(fileName) {
		return HORIZONS
	}
	if isOdysseyFile(fileName) {
		return ODYSSEY
	}
	return UNKNOWN
}

func isHorizonsFile(fileName string) bool {
	// Horizons Example
	// Journal.220212014914.01.log
	match, err := regexp.MatchString("\\bJournal\\.\\d*\\.\\d*.log\\b", fileName)
	if err != nil {
		panic(err)
	}
	return match
}

func isOdysseyFile(fileName string) bool {
	// Odyssey Example
	// Journal.2022-05-20T170535.01.log
	match, err := regexp.MatchString("\\bJournal\\.\\d{4}-\\d{2}-\\d{2}T\\d*\\.\\d*\\.log\\b", fileName)
	if err != nil {
		panic(err)
	}
	return match
}
