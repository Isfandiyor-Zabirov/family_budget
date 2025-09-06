package access_groups

import (
	"family_budget/pkg/database"
	"log"
)

func getList() []AccessGroup {
	var list []AccessGroup
	err := database.Postgres().Find(&list).Error
	if err != nil {
		log.Println("AccessGroup getList func query error:", err.Error())
		return nil
	}

	return list
}
