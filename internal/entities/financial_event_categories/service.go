package financial_event_categories

import "family_budget/internal/utils/response"

func Create(fec *FinancialEventCategories) (resp response.ResponseModel, err error) {
	_, err = createFec(fec)
	if err != nil {
		response.SetResponseData(FinancialEventCategories{}, "Что-то пошло не так", false)
		return
	}

	response.SetResponseData(fec, "Категория успешно добавлена", true)
	return
}

func Update(fec *FinancialEventCategories) (FinancialEventCategories, error) {
	return updateFec(fec)
}

func Delete(id int) error {
	fec := &FinancialEventCategories{ID: id}
	return deleteFec(fec)
}

func Get(id int) (FinancialEventCategories, error) {
	return getFec(id)
}

func GetList(filters Filters) (resp response.ResponseModel, pagination response.Pagination, err error) {
	list, total, err := getList(filters)
	if err != nil {
		response.SetResponseData([]FinancialEventCategories{}, "Что-то пошло не так", false)
		response.SetPagination(0, 1, 1)
		return
	}

	response.SetResponseData(list, "Успех", true)
	response.SetPagination(response.CalculateTotalPages(total, filters.PageLimit), total, filters.CurrentPage)
	return
}
