package main

import (
	"bytes"
	"github.com/fatih/color"
	"io"
	"net/http"
	"strconv"
	"strings"
)

func CheckUpdate() {
	fullUrl := "https://raw.githubusercontent.com/CMDR-WDX/edpvp-log-prepare/master/version"

	buf := new(bytes.Buffer)
	// Stolen from https://golangdocs.com/golang-download-files
	client := http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			req.URL.Opaque = req.URL.Path
			return nil
		},
	}

	resp, err := client.Get(fullUrl)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		color.Blue("Failed to check for new Version")
		return
	}
	_, err = io.Copy(buf, resp.Body)

	versionAsString := strings.Trim(strings.Split(buf.String(), "\n")[0], " ")

	compareWithCurrentVersion(versionAsString, VERSION)
}

func compareWithCurrentVersion(online string, local string) {

	onlineArr := strings.Split(online, ",")
	localArr := strings.Split(local, ",")
	length := len(onlineArr)
	if length != len(localArr) {
		if length < len(localArr) {
			length = len(localArr)
			// Online needs padding
			for len(onlineArr) < length {
				onlineArr = append(onlineArr, "0")
			}
		} else {
			// Local needs padding
			for len(localArr) < length {
				localArr = append(localArr, "0")
			}
		}
	}
	if len(onlineArr) != len(localArr) {
		panic("Should never happen")
	}

	isOnlineMoreRecent := false
	for i := 0; i < length; i++ {
		localVal, _ := strconv.Atoi(localArr[i])
		onlineVal, _ := strconv.Atoi(onlineArr[i])

		if onlineVal > localVal {
			isOnlineMoreRecent = true
			break
		}
	}
	if isOnlineMoreRecent {
		color.Yellow("*** NEW VERSION AVAILABLE ***\nCurrent Version: v%s\nAvailable Version: v%s\n"+
			"Get the Update at https://github.com/CMDR-WDX/edpvp-log-prepare/releases", online, local)
	}
}
