package financial_events

import (
	"errors"
	"family_budget/internal/utils/response"
	"log"

	"gorm.io/gorm"
)

func Create(finEvent *FinancialEvent) (resp response.ResponseModel, err error) {
	createdFinEvent, err := createFinancialEvent(finEvent)
	if err != nil {
		log.Printf("CreateFinancialEvent create financial event error: %s", err.Error())
		response.SetResponseData(&resp, FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, createdFinEvent, "Финансовая события успешно создана", true, 0, 0, 0)
	return
}

func Update(finEvent *FinancialEvent) (resp response.ResponseModel, err error) {
	f, err := getFinancialEvent(finEvent.ID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SetResponseData(&resp, FinancialEvent{}, "Финансовая события не найдена", false, 0, 0, 0)
			return
		}
		log.Printf("UpdateFinancialEvent get financial event error: %s", err.Error())
		response.SetResponseData(&resp, FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if f.FamilyID != finEvent.FamilyID {
		response.SetResponseData(&resp, FinancialEvent{}, "Доступ к чужим данным запрещен", false, 0, 0, 0)
		return
	}

	updatedFinEvent, err := updateFinancialEvent(finEvent)
	if err != nil {
		log.Printf("UpdateFinancialEvent update financial event error: %s", err.Error())
		response.SetResponseData(&resp, FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, updatedFinEvent, "Финансовая события успешно обновлена", true, 0, 0, 0)
	return
}

func Delete(id, familyID int) (resp response.ResponseModel, err error) {
	f, err := getFinancialEvent(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SetResponseData(&resp, FinancialEvent{}, "Финансовая события не найдена", false, 0, 0, 0)
			return
		}
		log.Printf("DeleteFinancialEvent get financial event error: %s", err.Error())
		response.SetResponseData(&resp, FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if f.FamilyID != familyID {
		response.SetResponseData(&resp, FinancialEvent{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		return
	}

	err = deleteFinancialEvent(&f)
	if err != nil {
		log.Printf("DeleteFinancialEvent delete financial event error: %s", err.Error())
		response.SetResponseData(&resp, FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, f, "Финансовое событие успешно удалено", true, 0, 0, 0)
	return
}

func Get(id int, familyID int) (resp response.ResponseModel, err error) {
	f, err := getFinancialEvent(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.SetResponseData(&resp, FinancialEvent{}, "Финансовая события не найдена", false, 0, 0, 0)
			return
		}
		log.Printf("GetFinancialEvent get financial event error: %s", err.Error())
		response.SetResponseData(&resp, FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	if f.FamilyID != familyID {
		response.SetResponseData(&resp, FinancialEvent{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, f, "Успех", true, 0, 0, 0)
	return
}

func GetList(filters Filters) (resp response.ResponseModel, err error) {
	list, total, err := getFinancialEventList(filters)
	if err != nil {
		log.Printf("GetFinancialEventsList get financial events list error: %s", err.Error())
		response.SetResponseData(&resp, []FinancialEvent{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, list, "Успех", true, filters.PageLimit, total, filters.CurrentPage)
	return
}
