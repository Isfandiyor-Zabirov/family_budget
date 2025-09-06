package migration

import (
	"family_budget/internal/entities/family"
	"family_budget/internal/entities/roles"
	"family_budget/internal/entities/users"
	"family_budget/pkg/database"
	"fmt"
	"log"
)

func AutoMigrate() {
	fmt.Println("Automatically migrating the schemas...")

	err := database.Postgres().AutoMigrate(
		&family.Family{},
		&roles.Roles{},
		&users.User{},
	)
	if err != nil {
		log.Println("AutoMigrate func error: ", err.Error())
		return
	}
	fmt.Println("Finished migrating the schemas...")
}
