package reports

import "family_budget/internal/utils/response"

func GetMainReport(familyID int) (resp response.ResponseModel, err error) {
	report, err := getMainReport(familyID)
	if err != nil {
		response.SetResponseData(&resp, MainReport{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, report, "Успех", true, 0, 0, 0)
	return
}

func GetGraphReport(filter Filter) (resp response.ResponseModel, err error) {
	report, err := getGraphReport(filter)
	if err != nil {
		response.SetResponseData(&resp, MainReport{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, report, "Успех", true, 0, 0, 0)
	return
}
