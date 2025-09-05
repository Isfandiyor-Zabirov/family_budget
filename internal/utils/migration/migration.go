package migration

import (
	"family_budget/pkg/database"
	"fmt"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var dbName = "family_budget"
var dsnWithoutDB = "host=localhost user=postgres password=postgres sslmode=disable"
var dsnWithDB = fmt.Sprintf("%s dbname=%s", dsnWithoutDB, dbName)

func AutoMigrate() {
	fmt.Println("Automatically migrating the schemas...")

	serverDB, err := gorm.Open(postgres.Open(dsnWithoutDB), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to postgres server: ", err)
	}

	createDBQuery := fmt.Sprintf("CREATE DATABASE %s;", dbName)
	err = serverDB.Exec(createDBQuery).Error
	if err != nil {
		if err.Error() != fmt.Sprintf("pq: database \"%s\" already exists", dbName) {
			log.Fatal("failed to create database: ", err)
		}
	}

	db, err := gorm.Open(postgres.Open(dsnWithDB), &gorm.Config{})
	if err != nil {
		log.Fatal("failed to connect to target database: ", err)
	}

	err = db.AutoMigrate()
	if err != nil {
		log.Println("error encountered while migrating the schema: ", err)
		return
	}

	err = database.Postgres().AutoMigrate()
	if err != nil {
		log.Println("error encountered while migrating the schema: ", err)
		return
	}
	fmt.Println("Finished migrating the schemas...")
}
