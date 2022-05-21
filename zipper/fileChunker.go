package zipper

import (
	"archive/zip"
	"edpvp-log-prepare/config"
	"github.com/fatih/color"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"time"
)

// List of files is split into chunks of files
// Each chunk is < 10MB big. This way, after compression, the chunk will be (hopefully) < 8MB after ZIP compression
func createFileChunks(fileNames []string) [][]string {
	returnArray := make([][]string, 0)
	currentChunk := make([]string, 0)
	var currentChunkSize int64 = 0
	for _, file := range fileNames {
		fileStat, _ := os.Stat(file)
		sizeWithCurrentFile := currentChunkSize + fileStat.Size()
		if sizeWithCurrentFile > 1024*1024*10 && len(currentChunk) > 0 { // 10 MByte
			returnArray = append(returnArray, currentChunk)
			currentChunk = make([]string, 0)
			currentChunkSize = 0
		}

		currentChunk = append(currentChunk, file)
		currentChunkSize = sizeWithCurrentFile
	}
	// if there are still files, create the tailing chunk
	if len(currentChunk) > 0 {
		returnArray = append(returnArray, currentChunk)
	}

	return returnArray
}

func ZipUpFiles(directory string, cfg config.AppConfig) string {
	dirStat, err := ioutil.ReadDir(directory)

	if err != nil {
		panic(err)
	}

	fileNames := make([]string, 0)
	for _, file := range dirStat {
		fileNames = append(fileNames, filepath.Join(directory, file.Name()))
	}

	namePrefix := generateNamePrefix()
	outputPath := filepath.Join(cfg.OutputDir, namePrefix)
	err = os.MkdirAll(outputPath, 0633)
	if err != nil {
		panic(err)
	}

	chunks := createFileChunks(fileNames)
	for i, chunk := range chunks {

		name := filepath.Join(cfg.OutputDir, namePrefix, strconv.Itoa(i+1)+".zip")
		color.Green("Building chunk " + name)
		zipUpChunk(chunk, name)
	}
	return outputPath
}

func generateNamePrefix() string {
	t := time.Now()
	return t.Format("2006-01-02-15-04-05")
}

func zipUpChunk(chunk []string, name string) {
	archive, err := os.Create(name)
	if err != nil {
		color.Red("Failed to build Chunk " + name)
	}
	defer archive.Close()

	zipWriter := zip.NewWriter(archive)
	for _, fileName := range chunk {
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
}
