package migration

import (
	"family_budget/pkg/client"
	"fmt"
	"log"
)

func AutoMigrate() {
	fmt.Println("Automatically migrating the schemas...")
	err := client.Postgres().AutoMigrate()
	if err != nil {
		log.Println("error encountered while migrating the schema: ", err)
		return
	}
	fmt.Println("Finished migrating the schemas...")
}
