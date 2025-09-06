package financial_event_categories

import (
	"errors"
	"family_budget/internal/utils/response"
	"log"
)

func Create(fec *FinancialEventCategories) (resp response.ResponseModel, err error) {
	_, err = createFec(fec)
	if err != nil {
		response.SetResponseData(&resp, FinancialEventCategories{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, fec, "Категория успешно добавлена", true, 0, 0, 0)
	return
}

func Update(fec *FinancialEventCategories) (resp response.ResponseModel, err error) {
	updated, err := updateFec(fec)
	if err != nil {
		response.SetResponseData(&resp, FinancialEventCategories{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, updated, "Категория успешно обновлена", true, 0, 0, 0)

	return
}

func Delete(id, familyID int) (resp response.ResponseModel, err error) {
	// TODO: надо возвращать конкретные ошибки типо fec не найден или база недоступна

	fec, err := getFec(id)
	if err != nil {
		response.SetResponseData(&resp, FinancialEventCategories{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if fec.FamilyID != familyID {
		response.SetResponseData(&resp, FinancialEventCategories{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		err = errors.New("Family IDs not matched ")
		log.Println("FinancialEventCategories Delete func family ID checking error:", err.Error())
		return
	}

	err = deleteFec(&fec)
	if err != nil {
		response.SetResponseData(&resp, FinancialEventCategories{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}
	response.SetResponseData(&resp, fec, "Категория успешно удалена", true, 0, 0, 0)
	return
}

func Get(id int) (FinancialEventCategories, error) {
	return getFec(id)
}

func GetList(filters Filters) (resp response.ResponseModel, err error) {
	list, total, err := getList(filters)
	if err != nil {
		response.SetResponseData(&resp, []FinancialEventCategories{}, "Что-то пошло не так", false, filters.PageLimit, total, filters.CurrentPage)
		return
	}

	response.SetResponseData(&resp, list, "Успех", true, filters.PageLimit, total, filters.CurrentPage)
	return
}
