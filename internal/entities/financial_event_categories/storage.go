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

func getList(filters Filters) (resp []FinancialEventCategories, totalRows int64, err error) {
	query := database.Postgres().Table("financial_event_categories fec").
		Where("fec.family_id = ? and fec.deleted_at is null", filters.FamilyID)

	if filters.CurrentPage != 0 {
		page := 1
		filters.CurrentPage = page
	}

	if filters.PageLimit == 0 {
		pageLimit := 20
		filters.PageLimit = pageLimit
	}

	if filters.Search != nil {
		searchText := "%" + *filters.Search + "%"
		query = query.Where("fec.name ilike ? or description ilike ?", searchText, searchText)
	}

	err = query.Count(&totalRows).Error
	if err != nil {
		log.Println("Failed to count FinancialEventCategories", err.Error())
		return []FinancialEventCategories{}, 0, err
	}

	err = query.Select("fec.*").Offset(filters.PageLimit*filters.CurrentPage - 1).
		Limit(filters.PageLimit).Error
	if err != nil {
		log.Println("FinancialEventCategories getList func query error:", err.Error())
		return []FinancialEventCategories{}, 0, err
	}

	return
}
