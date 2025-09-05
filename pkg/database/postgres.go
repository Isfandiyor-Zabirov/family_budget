package database

import (
	"family_budget/pkg/external_config"
	"fmt"
	"github.com/kr/pretty"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
)

var db *gorm.DB

// SetupDB initializes the database instance
func SetupDB() {
	fmt.Println("Connecting to database...")
	var err error
	var conf = external_config.ExternalConfigs

	db, err = gorm.Open(postgres.Open(fmt.Sprintf("host=%s user=%s port=%s dbname=%s password=%s",
		conf.PostgreSQL.Host,
		conf.PostgreSQL.User,
		conf.PostgreSQL.Port,
		conf.PostgreSQL.DbName,
		conf.PostgreSQL.Pass)),
	)
	if err != nil {
		log.Fatalf("open db err:", err.Error())
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(100)
	//миграция БД, отключать если не нужно
	//utils.AutoMigrate()
	pretty.Logln("Database successfully connected! ")
	fmt.Println("Database successfully connected!")
}

// CloseDB closes database connection (unnecessary)
func CloseDB() {
	sqlDB, err := db.DB()
	sqlDB.Close()
	if err != nil {
		pretty.Logln("Error on closing the DB: ", err)
	}
}

func Postgres() *gorm.DB {
	return db
}
