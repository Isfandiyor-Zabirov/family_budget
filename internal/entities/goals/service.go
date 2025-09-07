package goals

import (
	"errors"
	"family_budget/internal/utils/response"

	"gorm.io/gorm"
)

func CreateGoal(goal *Goals) (resp response.ResponseModel, err error) {
	createdGoal, err := createGoal(goal)
	if err != nil {
		response.SetResponseData(&resp, createdGoal, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, createdGoal, "Цель успешно создана", true, 0, 0, 0)
	return
}

func UpdateGoal(goal *Goals) (resp response.ResponseModel, err error) {
	g, err := getGoal(goal.ID)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		response.SetResponseData(&resp, Goals{}, "Цель для обновления не найдена", false, 0, 0, 0)
		return
	}

	if g.FamilyID != goal.FamilyID {
		response.SetResponseData(&resp, Goals{}, "Доступ к чужим данным запрещен", false, 0, 0, 0)
		return
	}

	updatedGoal, err := updateGoal(goal)
	if err != nil {
		response.SetResponseData(&resp, Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, updatedGoal, "Цель успешно обновлена", true, 0, 0, 0)
	return
}

func DeleteGoal(id, familyID int) (resp response.ResponseModel, err error) {
	goal, err := getGoal(id)
	if err != nil {
		response.SetResponseData(&resp, Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if goal.FamilyID != familyID {
		response.SetResponseData(&resp, Goals{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		return
	}

	err = deleteGoal(&goal)
	if err != nil {
		response.SetResponseData(&resp, Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, goal, "Цель успешно удалена", true, 0, 0, 0)
	return
}

func GetGoalsList(filters Filters) (resp response.ResponseModel, err error) {
	list, total, err := getGoals(filters)
	if err != nil {
		response.SetResponseData(&resp, []Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, list, "Успех", true, filters.PageLimit, total, filters.CurrentPage)
	return
}

func GetGoal(id int, familyID int) (resp response.ResponseModel, err error) {
	goal, err := getGoal(id)
	if err != nil {
		response.SetResponseData(&resp, Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if goal.FamilyID != familyID {
		response.SetResponseData(&resp, Goals{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, goal, "Успех", true, 0, 0, 0)
	return
}
