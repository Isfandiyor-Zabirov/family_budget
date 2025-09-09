package transactions

import "family_budget/internal/utils/response"

func CreateTransaction(inputs *Transactions) (resp response.ResponseModel, statusCode int, err error) {
	statusCode, message, err := createTransaction(inputs)
	if err != nil {
		response.SetResponseData(&resp, struct{}{}, message, false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, struct{}{}, message, true, 0, 0, 0)
	return
}

func TransactionList(filters Filters) (resp response.ResponseModel, err error) {
	list, total, err := getTransactionList(filters)
	if err != nil {
		response.SetResponseData(&resp, []Response{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, list, "Что-то пошло не так", false, filters.PageLimit, total, filters.CurrentPage)
	return
}
