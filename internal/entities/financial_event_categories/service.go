package financial_event_categories

import (
	"errors"
	"family_budget/internal/utils/response"
	"log"
)

func Create(fec *FinancialEventCategories) (resp response.ResponseModel, err error) {
	_, err = createFec(fec)
	if err != nil {
		resp = response.SetResponseData(FinancialEventCategories{}, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(fec, "Категория успешно добавлена", true)
	return
}

func Update(fec *FinancialEventCategories) (resp response.ResponseModel, err error) {
	updated, err := updateFec(fec)
	if err != nil {
		resp = response.SetResponseData(FinancialEventCategories{}, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(updated, "Категория успешно обновлена", true)

	return
}

func Delete(id, familyID int) (resp response.ResponseModel, err error) {
	fec, err := getFec(id)
	if err != nil {
		resp = response.SetResponseData(FinancialEventCategories{}, "Что-то пошло не так", false)
		return
	}

	if fec.FamilyID != familyID {
		resp = response.SetResponseData(FinancialEventCategories{}, "Нет доступа к чужим данным", false)
		err = errors.New("Family IDs not matched ")
		log.Println("FinancialEventCategories Delete func family ID checking error:", err.Error())
		return
	}

	err = deleteFec(&fec)
	if err != nil {
		resp = response.SetResponseData(FinancialEventCategories{}, "Что-то пошло не так", false)
		return
	}
	resp = response.SetResponseData(fec, "Категория успешно удалена", true)
	return
}

func Get(id int) (FinancialEventCategories, error) {
	return getFec(id)
}

func GetList(filters Filters) (resp response.ResponseModel, pagination response.Pagination, err error) {
	list, total, err := getList(filters)
	if err != nil {
		resp = response.SetResponseData([]FinancialEventCategories{}, "Что-то пошло не так", false)
		pagination = response.SetPagination(0, 1, 1)
		return
	}

	resp = response.SetResponseData(list, "Успех", true)
	pagination = response.SetPagination(response.CalculateTotalPages(total, filters.PageLimit), total, filters.CurrentPage)
	return
}
