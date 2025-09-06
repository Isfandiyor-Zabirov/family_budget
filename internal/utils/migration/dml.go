package migration

import (
	"family_budget/pkg/database"
	"fmt"
	"log"
)

const (
	defaultFamily = `INSERT INTO families (id, name, phone, created_at, updated_at) 
	VALUES (1, 'Family Budget', '00000000', now(), now())`

	headOfFamilyRole = `INSERT INTO roles (id, family_id, name, description, created_at, updated_at) 
	VALUES (1, 1, 'Глава семьи', 'Основная роль для глави семьи', now(), now())`

	defaultAccesses = `INSERT INTO accesses (id, code, name, description, created_at, updated_at) 
	VALUES
	(1, 'READ', 'Просмотр', 'Доступ просмотра', now(), now()),
	(2, 'CREATE', 'Создание', 'Доступ создания', now(), now()),
	(3, 'UPDATE', 'Изменение', 'Доступ изменения', now(), now()),
	(4, 'DELETE', 'Удаление', 'Доступ удаления', now(), now())`

	defaultAccessGroups = `INSERT INTO access_groups (id, code, name, description, created_at, updated_at) values
        (1, 'FINANCIAL_EVENT_CATEGORIES', 'Категории расходов и доходов', 'Категории расходов или доходов',  now(), now()),
        (2, 'FINANCIAL_EVENTS', 'Расходы и доходы', 'Расходы и доходы', now(), now()),
        (3, 'GOALS', 'Цели', 'Семейные цели', now(), now()),
        (4, 'ROLES', 'Роли', 'Роли для доступа к функционалу системы (можно создавать самостоятельно)', now(), now()),
		(5, 'TRANSACTIONS', 'Операции', 'Создани операции по расходам и доходам', now(), now()),
		(6, 'USERS', 'Семейство', 'Семейство', now(), now())`
)

// initDml to insert minimum default values to start the app
func initDml() {
	fmt.Println("Starting inserting default values...")
	db := database.Postgres()
	var errorList []string
	var err error

	err = db.Exec(defaultFamily).Error
	if err != nil {
		errorList = append(errorList, err.Error())
	}

	err = db.Exec(headOfFamilyRole).Error
	if err != nil {
		errorList = append(errorList, err.Error())
	}

	err = db.Exec(defaultAccesses).Error
	if err != nil {
		errorList = append(errorList, err.Error())
	}

	err = db.Exec(defaultAccessGroups).Error
	if err != nil {
		errorList = append(errorList, err.Error())
	}

	if errorList != nil {
		log.Println("Error on adding default values: ", errorList)
	}
	fmt.Println("Inserting default values successfully completed...")
}
