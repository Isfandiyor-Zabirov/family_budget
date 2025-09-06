package goals

import (
	"errors"
	"family_budget/internal/utils/response"
	"log"
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
	updatedGoal, err := updateGoal(goal)
	if err != nil {
		response.SetResponseData(&resp, updatedGoal, "Что-то пошло не так", false, 0, 0, 0)
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
		err = errors.New("Family IDs not matched ")
		log.Println("DeleteGoal func family ID checking error:", err.Error())
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

func GetGoalsList(filters Filters) (resp response.ResponseModel, pagination response.Pagination, err error) {
	list, total, err := getGoals(filters)
	if err != nil {
		response.SetResponseData(&resp, []Goals{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, list, "Успех", true, filters.PageLimit, total, filters.CurrentPage)
	return
}
