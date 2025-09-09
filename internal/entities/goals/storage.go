package goals

import (
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
	"log"
)

func createGoal(goal *Goals) (Goals, error) {
	repo := crud.NewRepository[Goals]()
	db := database.Postgres()
	return repo.Create(db, goal)
}

func updateGoal(goal *Goals) (Goals, error) {
	repo := crud.NewRepository[Goals]()
	db := database.Postgres()
	return repo.Update(db, goal)
}

func deleteGoal(goal *Goals) error {
	repo := crud.NewRepository[Goals]()
	db := database.Postgres()
	return repo.Delete(db, goal)
}

func getGoal(id int) (Goals, error) {
	repo := crud.NewRepository[Goals]()
	db := database.Postgres()
	return repo.Get(db, id)
}

func getGoals(filters Filters) (resp []Goals, totalRows int64, err error) {
	resp = []Goals{}
	query := database.Postgres().Table("goals g").
		Where("g.deleted_at IS NULL and g.family_id = ?", filters.FamilyID)

	if filters.Search != nil {
		search := "%" + *filters.Search + "%"
		query = query.Where("g.name ilike ? or g.description ilike ?", search, search)
	}

	if filters.Status != nil {
		query = query.Where("g.status = ?", filters.Status)
	}

	if filters.DueDateFrom != nil {
		query = query.Where("g.due_date::date >= ?", filters.DueDateFrom)
	}

	if filters.DueDateTo != nil {
		query = query.Where("g.due_date::date <= ?", filters.DueDateTo)
	}

	err = query.Count(&totalRows).Error
	if err != nil {
		log.Println("Failed to count Goals", err.Error())
		return []Goals{}, 0, err
	}

	if filters.CurrentPage == 0 {
		filters.CurrentPage = 1
	}

	if filters.PageLimit == 0 {
		pageLimit := 20
		filters.PageLimit = pageLimit
	}

	err = query.Select("g.*").Offset(filters.PageLimit * (filters.CurrentPage - 1)).
		Limit(filters.PageLimit).Order("g.id desc").Scan(&resp).Error

	return
}
