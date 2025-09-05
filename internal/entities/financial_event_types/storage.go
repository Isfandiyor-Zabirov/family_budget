package financial_event_types

import (
	"family_budget/pkg/database"
	"log"
)

func getFinancialEventTypes() ([]FinancialEventType, error) {
	var response []FinancialEventType
	err := database.Postgres().Find(&response).Order("id").Error
	if err != nil {
		log.Println("getFinancialEventTypes func query error:", err.Error())
		return []FinancialEventType{}, err
	}
	return response, nil
}
