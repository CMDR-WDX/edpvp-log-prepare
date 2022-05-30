/**
This takes an Entire List of Logs. It will write them into a ZIP and check if the ZIP is > 8 MB. If yes, the contents
are split and two ZIPs are created instead.

*/

package zipper

import (
	"archive/zip"
	"github.com/fatih/color"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

const ZipFileMaxSize = 1024 * 1024 * 8 // 8 MB

var filenameCounter = 1

type ZipResult struct {
	logsInZip  []Log
	outputFile string
}

type Log struct {
	Size     int64
	Filepath string
}

func FilePathsAsLogs(logs []string) ([]Log, error) {
	returnArray := make([]Log, 0, len(logs))
	for _, entry := range logs {
		data, err := os.Stat(entry)
		if err != nil {
			return returnArray, err
		}
		returnArray = append(returnArray, Log{
			Size:     data.Size(),
			Filepath: entry,
		})
	}
	return returnArray, nil
}

func CreateZipFilesAndMoveToOutputDirectory(logsInTempDir []Log, baseOutDir string) (string, error) {
	color.Blue("Starting Zipping...")
	timeBefore := time.Now()
	// Create a directory under BaseOutDir with the current Timestamp
	directoryName := filepath.Join(baseOutDir, generateDirectoryPrefix())
	err := os.MkdirAll(directoryName, 0633)
	if err != nil {
		return "", err
	}

	_, err = zipUpLogs(logsInTempDir, directoryName)
	if err != nil {
		return "", err
	}

	timeDiff := time.Now().Sub(timeBefore)

	color.Blue("Zipping Done. Took a total of %f Seconds", timeDiff.Seconds())

	return directoryName, nil
}

func generateDirectoryPrefix() string {
	t := time.Now()
	return t.Format("2006-01-02-15-04-05")
}

func zipUpLogs(logs []Log, outDir string) ([]string, error) {
	result, err := zipLogs(logs, outDir)
	if err != nil {
		return make([]string, 0), err
	}
	// Get file stats for Output.
	info, err := os.Stat(result.outputFile)
	if err != nil {
		return make([]string, 0), err
	}

	returnArr := make([]string, 0)

	if info.Size() > ZipFileMaxSize {
		// The Zip is greater than 8 MB. Split up the Zip File in half. Can be solved recursively
		color.Yellow("Chunk %s is too big (%f MB). Splitting up", info.Name(), float64(info.Size())/float64(1024*1024))
		var cutOff int64 = 0
		// calculate cutoff
		for _, log := range logs {
			cutOff += log.Size
		}
		cutOff /= 2
		var sum int64 = 0
		leftSide := make([]Log, 0)
		rightSide := make([]Log, 0)
		// delete the ZIP file
		err = os.Remove(result.outputFile)
		if err != nil {
			return make([]string, 0), err
		}

		for _, log := range result.logsInZip {
			sum += log.Size
			if sum < cutOff {
				leftSide = append(leftSide, log)
			} else {
				rightSide = append(rightSide, log)
			}
		}

		// The logs have been split in two

		leftSizeResult, err := zipUpLogs(leftSide, outDir)
		if err != nil {
			return make([]string, 0), err
		}
		rightSizeResult, err := zipUpLogs(leftSide, outDir)
		if err != nil {
			return make([]string, 0), err
		}
		for _, entry := range leftSizeResult {
			returnArr = append(returnArr, entry)
		}
		for _, entry := range rightSizeResult {
			returnArr = append(returnArr, entry)
		}
	} else {
		color.Green("Chunk %s zipped up successfully (%f MB)", info.Name(), float64(info.Size())/float64(1024*1024))
		returnArr = append(returnArr, result.outputFile)
	}
	return returnArr, nil

}

func zipLogs(logs []Log, outDir string) (ZipResult, error) {
	chunkName := strconv.Itoa(filenameCounter) + ".zip"
	filenameCounter++
	fullResultFilePath := filepath.Join(outDir, chunkName)
	archive, err := os.Create(fullResultFilePath)
	if err != nil {
		color.Red("Failed to build Chunk " + chunkName)
		return ZipResult{}, err
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	for _, entry := range logs {
		fileName := entry.Filepath
		_, fileNameNoPath := filepath.Split(fileName)
		file, err := os.Open(fileName)
		if err != nil {
			color.Red("Failed to open File " + fileName)
			color.Red(err.Error())
			continue
		}
		defer file.Close()

		fileInArchive, err := zipWriter.Create(fileNameNoPath)
		if err != nil {
			color.Red("Failed to create File in Archive " + fileName)
			color.Red(err.Error())
			continue
		}

		_, err = io.Copy(fileInArchive, file)

		if err != nil {
			color.Red("Failed to put File into Archive: " + fileName)
		}
	}
	zipWriter.Close()
	return ZipResult{
		logsInZip:  logs,
		outputFile: fullResultFilePath,
	}, nil
}
