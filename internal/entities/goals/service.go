package goals

import (
	"errors"
	"family_budget/internal/utils/response"
	"log"
)

func CreateGoal(goal *Goals) (resp response.ResponseModel, err error) {
	createdGoal, err := createGoal(goal)
	if err != nil {
		resp = response.SetResponseData(createdGoal, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(createdGoal, "Цель успешно создана", true)
	return
}

func UpdateGoal(goal *Goals) (resp response.ResponseModel, err error) {
	updatedGoal, err := updateGoal(goal)
	if err != nil {
		resp = response.SetResponseData(updatedGoal, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(updatedGoal, "Цель успешно обновлена", true)
	return
}

func DeleteGoal(id, familyID int) (resp response.ResponseModel, err error) {
	goal, err := getGoal(id)
	if err != nil {
		resp = response.SetResponseData(Goals{}, "Что-то пошло не так", false)
		return
	}

	if goal.FamilyID != familyID {
		resp = response.SetResponseData(Goals{}, "Нет доступа к чужим данным", false)
		err = errors.New("Family IDs not matched ")
		log.Println("DeleteGoal func family ID checking error:", err.Error())
		return
	}

	err = deleteGoal(&goal)
	if err != nil {
		resp = response.SetResponseData(Goals{}, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(goal, "Цель успешно удалена", true)
	return
}

func GetGoalsList(filters Filters) (resp response.ResponseModel, pagination response.Pagination, err error) {
	list, total, err := getGoals(filters)
	if err != nil {
		resp = response.SetResponseData([]Goals{}, "Что-то пошло не так", false)
		pagination = response.SetPagination(0, 1, 1)
		return
	}

	resp = response.SetResponseData(list, "Успех", true)
	pagination = response.SetPagination(response.CalculateTotalPages(total, filters.PageLimit), total, filters.CurrentPage)
	return
}
