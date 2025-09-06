package family

import (
	"family_budget/internal/utils/response"
)

func Create(family *Family) (resp response.ResponseModel, err error) {
	_, err = createFamily(family)
	if err != nil {
		resp = response.SetResponseData(Family{}, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(family, "Семья успешно добавлена", true)
	return
}

func Update(family *Family) (resp response.ResponseModel, err error) {
	updated, err := updateFamily(family)
	if err != nil {
		resp = response.SetResponseData(Family{}, "Что-то пошло не так", false)
		return
	}

	resp = response.SetResponseData(updated, "Семья успешно обновлена", true)

	return
}

func Delete(family *Family) (resp response.ResponseModel, err error) {
	err = deleteFamily(family)
	if err != nil {
		resp = response.SetResponseData(Family{}, "Что-то пошло не так", false)
		return
	}
	resp = response.SetResponseData(family, "Семья успешно удалена", true)
	return
}

func Get(id int) (Family, error) {
	return getFamily(id)
}
