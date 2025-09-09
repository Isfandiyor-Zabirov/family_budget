package transactions

import (
	"errors"
	"family_budget/internal/entities/goals"
	"family_budget/internal/entities/users"
	"family_budget/pkg/database"
	"log"
	"net/http"
)

func createTransaction(inputs *Transactions) (statusCode int, message string, err error) {
	tx := database.Postgres().Begin()

	// check user remaining limit
	user, err := users.GetUserInternal(inputs.UserID)
	if err != nil {
		tx.Rollback()
		return http.StatusInternalServerError, "Ошибка в получении даных пользователя", err
	}

	if user.RemainingLimit < inputs.Amount {
		tx.Rollback()
		return http.StatusBadRequest, "Ваш лимит закончился", errors.New("user limit reached")
	}

	// create the transaction
	err = tx.Create(inputs).Error
	if err != nil {
		tx.Rollback()
		log.Println("createTransaction func create transaction query error:", err.Error())
		return http.StatusInternalServerError, "Не удалось создать операцию", err
	}

	// update user remaining limit
	updateUserLimitQuery := `UPDATE users set remaining_limit = remaining_limit - ? WHERE id = ?`
	err = tx.Exec(updateUserLimitQuery, inputs.Amount, inputs.UserID).Error
	if err != nil {
		tx.Rollback()
		log.Println("createTransaction func update user limit query error:", err.Error())
		return http.StatusInternalServerError, "Ошибка при обновлении лимита пользователя", err
	}

	if inputs.GoalID != 0 {
		// check goal remaining budget
		goal, err := goals.GetGoalInternal(inputs.GoalID)
		if err != nil {
			tx.Rollback()
			return http.StatusInternalServerError, "Ошибка в получении данных цели", err
		}

		if goal.RemainingBudget < inputs.Amount {
			tx.Rollback()
			log.Println("createTransaction func goal remaining budget less than amount error")
			return http.StatusBadRequest, "Остаток бюджета ", errors.New("goal.RemainingBudget < inputs.Amount")
		}

		updateGoalBudgetQuery := `UPDATE goals SET remaining_budget = remaining_budget - ? WHERE id = ?`
		err = tx.Exec(updateGoalBudgetQuery, inputs.Amount, inputs.GoalID).Error
		if err != nil {
			tx.Rollback()
			log.Println("createTransaction func update goals remaining_budget query error:", err.Error())
			return http.StatusInternalServerError, "Ошибка при обновлении бюджета цели", err
		}
	}

	tx.Commit()
	return http.StatusOK, "Операция успешно создана", nil
}

func getTransactionList(filters Filters) (resp []Response, totalRows int64, err error) {
	resp = []Response{}
	query := database.Postgres().Table("transactions t").
		Where("t.family_id = ? AND t.deleted_at IS NULL", filters.FamilyID)

	joinsSql := `LEFT JOIN users u ON u.id = t.user_id
				 LEFT JOIN goals g ON g.id = t.goal_id
				 LEFT JOIN financial_events f ON f.id = t.financial_event_id`

	query = query.Joins(joinsSql)

	if filters.Search != nil {
		searchText := "%" + *filters.Search + "%"
		query = query.Where(`t.description ILIKE ? OR g.name ILIKE ? OR f.name ILIKE ?`, searchText, searchText, searchText)
	}

	if filters.UserID != nil {
		query = query.Where(`t.user_id = ?`, filters.UserID)
	}

	if filters.GoalID != nil {
		query = query.Where(`t.goal_id = ?`, filters.GoalID)
	}

	if filters.FinancialEventID != nil {
		query = query.Where(`t.financial_event_id = ?`, filters.FinancialEventID)
	}

	if filters.DateFrom != nil {
		query = query.Where(`t.created_at::date >= ?`, filters.DateFrom)
	}

	if filters.DateTo != nil {
		query = query.Where(`t.created_at::date <= ?`, filters.DateTo)
	}

	err = query.Count(&totalRows).Error
	if err != nil {
		log.Println("getTransactionList query count error:", err.Error())
		return
	}

	if filters.CurrentPage == 0 {
		filters.CurrentPage = 1
	}

	if filters.PageLimit == 0 {
		filters.PageLimit = 20
	}

	selectQuery := `t.*, u.name || u.surname as user_name, f.name as financial_event_name, g.name as goal_name`

	err = query.Select(selectQuery).Offset(filters.PageLimit * (filters.CurrentPage - 1)).
		Limit(filters.PageLimit).Order("t.id desc").Scan(&resp).Error

	return
}
