package main

import (
	"family_budget/internal/internal_config"
	"family_budget/pkg/external_config"
)

func main() {
	external_config.ExternalSetup("./pkg/external_config/external_configs.json")
	internal_config.InternalSetup("./internal/internal_config/internal_configs.json")
}
