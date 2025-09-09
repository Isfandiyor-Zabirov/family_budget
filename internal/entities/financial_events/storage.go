package financial_events

import (
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
	"log"
)

func createFinancialEvent(event *FinancialEvent) (FinancialEvent, error) {
	repo := crud.NewRepository[FinancialEvent]()
	db := database.Postgres()
	return repo.Create(db, event)
}

func updateFinancialEvent(event *FinancialEvent) (FinancialEvent, error) {
	repo := crud.NewRepository[FinancialEvent]()
	db := database.Postgres()
	return repo.Update(db, event)
}

func deleteFinancialEvent(event *FinancialEvent) error {
	repo := crud.NewRepository[FinancialEvent]()
	db := database.Postgres()
	return repo.Delete(db, event)
}

func getFinancialEvent(id int) (FinancialEvent, error) {
	repo := crud.NewRepository[FinancialEvent]()
	db := database.Postgres()
	return repo.Get(db, id)
}

func getFinancialEventList(filters Filters) (resp []FinancialEvent, totalRows int64, err error) {
	resp = []FinancialEvent{}
	query := database.Postgres().Table("financial_events f").
		Where("f.deleted_at IS NULL and f.family_id = ?", filters.FamilyID)

	if filters.CategoryID != nil {
		query = query.Where("f.category_id = ?", *(filters.CategoryID))
	}

	if filters.Inflow != nil {
		query = query.Where("f.inflow = ?", *(filters.Inflow))
	}

	if filters.Search != nil {
		search := "%" + *filters.Search + "%"
		query = query.Where("f.name ilike ? or f.description ilike ?", search, search)
	}

	err = query.Count(&totalRows).Error
	if err != nil {
		log.Println("Failed to count FinancialEvents", err.Error())
		return []FinancialEvent{}, 0, err
	}

	if filters.CurrentPage == 0 {
		filters.CurrentPage = 1
	}

	if filters.PageLimit == 0 {
		filters.PageLimit = 20
	}

	err = query.Select("f.*").Offset(filters.PageLimit * (filters.CurrentPage - 1)).
		Limit(filters.PageLimit).Order("f.id desc").Scan(&resp).Error

	return
}
