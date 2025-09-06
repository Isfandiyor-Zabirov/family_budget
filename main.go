package main

import (
	_ "family_budget/docs"
	"family_budget/internal/api"
	"family_budget/internal/internal_config"
	"family_budget/internal/logger"
	"family_budget/internal/utils/migration"
	"family_budget/pkg/database"
	"family_budget/pkg/external_config"
)

// @title           Family Budget API
// @version         1.0
// @description     Welcome to family budget API.
// @termsOfService  http://swagger.io/terms/

// @contact.name   Isfandiyor Zabirov
// @contact.url    http://www.swagger.io/support
// @contact.email  isfandiyor.zabirov.sh@gmail.com

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	logger.InitLogger()
	external_config.ExternalSetup("./pkg/external_config/external_configs.json")
	internal_config.InternalSetup("./internal/internal_config/internal_configs.json")

	database.SetupDB()
	migration.AutoMigrate()
	api.Init()
}
