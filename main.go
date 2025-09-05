package main

import (
	"family_budget/internal/api"
	"family_budget/internal/internal_config"
	"family_budget/internal/logger"
	"family_budget/internal/utils/migration"
	"family_budget/pkg/database"
	"family_budget/pkg/external_config"
)

// @title Family Budget API документация
// @version 1.0.0
// @description Документация к сервису Family Budget.
// @termsOfService http://swagger.io/terms/
// @contact.name Family Budget
// @contact.email Family Budget
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8000
// @securitydefinitions.apikey  JWT
// @in                          header
// @name                        Authorization
// @accept json
// @produce json
// @schemes https
func main() {
	logger.InitLogger()
	external_config.ExternalSetup("./pkg/external_config/external_configs.json")
	internal_config.InternalSetup("./internal/internal_config/internal_configs.json")

	database.SetupDB()
	migration.AutoMigrate()
	api.Init()
}
