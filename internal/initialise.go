package internal

import (
	"fmt"
	"github.com/joho/godotenv"
	"gopkg.in/yaml.v3"
	"os"
)

func InitialiseApplication() {
	_ = godotenv.Load()
	configFilePath := os.Getenv("SHEEPSTOR_CONFIG_FILE_PATH")
	configData, err := os.ReadFile(configFilePath)
	err = yaml.Unmarshal(configData, &Registry)
	if err != nil {
		fmt.Print(err.Error() + "\n")
		fmt.Printf("Halting execution because Config file not loaded from '%s'\n", configFilePath)
		os.Exit(1)
	}
	Log, err = NewZapSugarLogger(true)
	if err != nil {
		fmt.Printf("Unable to initialise logging, halting: %s", err.Error())
		os.Exit(-1)
	}
	err = Registry.Initialise()
	if err != nil {
		Log.Errorf("Unable to initialise registry, halting: %s", err.Error())
		os.Exit(-1)
	}
	Log.Infof("Initialised")
	for _, website := range Registry.WebSites {
		Log.Infof("Website ID: %s", website.ID)
		Log.Infof("Website Webroot: %s", website.WebRoot)
		Log.Infof("BranchRef: %s", website.GitRepo.BranchRef)
	}
}
