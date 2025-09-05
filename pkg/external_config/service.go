package external_config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var ExternalConfigs Config

func ExternalSetup(filePath string) {
	fmt.Println("Setting up external config...")
	byteValue, err := os.ReadFile(filePath)
	if err != nil {
		log.Fatal("read config file err: ", err.Error())
		return
	}
	err = json.Unmarshal(byteValue, &ExternalConfigs)
	log.Println("externalConfigs:", ExternalConfigs)
	if err != nil {
		log.Fatal("unmarshall config error: ", err)
		return
	}
	fmt.Println("Successfully set up external config")
}
