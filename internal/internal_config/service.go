package internal_config

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
)

var InternalConfigs Config

func InternalSetup(F string) {
	fmt.Println("Setting up internal config...")
	byteValue, err := os.ReadFile(F)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	err = json.Unmarshal(byteValue, &InternalConfigs)
	if err != nil {
		log.Fatalf("%v", err)
		return
	}
	fmt.Println("Successfully set up internal config...")
}
