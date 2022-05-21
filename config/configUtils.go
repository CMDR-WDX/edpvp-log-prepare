package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

func BuildConfigDirIfNotExistsAndReturnDir() string {
	allConfigDir, err := os.UserConfigDir() // Appdata on Windows, /home/$USER/.config on Linux
	if err != nil {
		panic(err)
	}
	appConfigDir := filepath.Join(allConfigDir, "edpvpLogPrepare")

	_, err = os.Stat(appConfigDir)

	if os.IsNotExist(err) {
		fmt.Printf("%s does not exist. Trying to create\n", appConfigDir)
		err := os.Mkdir(appConfigDir, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}

		// Try again after creating dir
		_, err = os.Stat(appConfigDir)
		if err != nil {
			panic(err) // should not happen. Could happen in case of weird perms
		}
	}

	return appConfigDir
}

func GetConfig() (bool, AppConfig) {

	appConfigDir := BuildConfigDirIfNotExistsAndReturnDir()

	// Check if config exists
	appConfigFilePath := filepath.Join(appConfigDir, "config.json")

	file, err := os.OpenFile(appConfigFilePath, os.O_RDONLY, os.ModePerm)
	if err != nil {
		if os.IsNotExist(err) {
			fmt.Printf("%s does not exist\n", appConfigFilePath)
			return false, AppConfig{}
		} else {
			panic(err)
		}
	}

	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)

	// Parse file as json
	var config AppConfig

	err = json.Unmarshal(fileContent, &config)
	if err != nil {
		panic(err)
	}

	return true, config
}

func SetConfig(config AppConfig) error {
	configDir := BuildConfigDirIfNotExistsAndReturnDir()
	configAsBytes, err := json.MarshalIndent(config, "", "\t")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filepath.Join(configDir, "config.json"), configAsBytes, 0600)

	if err != nil {
		return err
	}
	return nil
}
