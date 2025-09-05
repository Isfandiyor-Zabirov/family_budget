package financial_event_categories

import (
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
	"log"
)

func createFec(fec *FinancialEventCategories) (FinancialEventCategories, error) {
	return crud.NewRepository[FinancialEventCategories]().Create(database.Postgres(), fec)
}

func updateFec(fec *FinancialEventCategories) (FinancialEventCategories, error) {
	return crud.NewRepository[FinancialEventCategories]().Update(database.Postgres(), fec)
}

func deleteFec(fec *FinancialEventCategories) error {
	return crud.NewRepository[FinancialEventCategories]().Delete(database.Postgres(), fec)
}

func getFec(id int) (FinancialEventCategories, error) {
	return crud.NewRepository[FinancialEventCategories]().Get(database.Postgres(), id)
}

func getList(familyID int) (resp []FinancialEventCategories, totalRows int64, err error) {
	err = database.Postgres().Find(&resp, "family_id = ?", familyID).Error
	if err != nil {
		log.Println("FinancialEventCategories getList func query error:", err.Error())
		return []FinancialEventCategories{}, 0, err
	}

	return
}
