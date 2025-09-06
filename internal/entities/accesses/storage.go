package accesses

import (
	"family_budget/pkg/database"
	"log"
)

func getList() []Access {
	var list []Access
	err := database.Postgres().Find(&list).Error
	if err != nil {
		log.Println("Access getList func query error:", err.Error())
		return nil
	}
	return list
}
