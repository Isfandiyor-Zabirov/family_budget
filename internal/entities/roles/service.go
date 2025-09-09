package roles

import (
	"errors"
	"family_budget/internal/utils/response"
	"family_budget/pkg/database"
	"gorm.io/gorm"
)

func GetRole(roleID int) (resp response.ResponseModel, err error) {
	role, err := getRole(roleID)
	if err != nil {
		response.SetResponseData(&resp, Role{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}
	response.SetResponseData(&resp, role, "Что-то пошло не так", false, 0, 0, 0)
	return
}

func DeleteRole(roleID, familyID int) (resp response.ResponseModel, err error) {
	r, err := getRole(roleID)
	if err != nil {
		response.SetResponseData(&resp, Role{}, "Роль не найден", false, 0, 0, 0)
		return resp, err
	}
	if r.ID == 0 {
		response.SetResponseData(&resp, Role{}, "Роль не найден или удален", false, 0, 0, 0)
		return resp, errors.New("role not found or deleted")
	}
	if r.FamilyID != familyID {
		response.SetResponseData(&resp, Role{}, "Нет доступа к чужим данным", false, 0, 0, 0)
		return resp, errors.New("deleting role failed: family ID mismatch")
	}
	err = deleteRole(&r)
	if err != nil {
		response.SetResponseData(&resp, Role{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}
	response.SetResponseData(&resp, r, "Успех", true, 0, 0, 0)

	return
}

func GetRoles(filter GetRolesFilter) (resp response.ResponseModel, err error) {

	list, total, err := getRoles(filter)
	if err != nil {
		response.SetResponseData(&resp, []GetRolesResp{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}
	response.SetResponseData(&resp, list, "Успех", true, filter.PageLimit, total, filter.CurrentPage)
	return
}

func CreateRoleWithAccesses(tx *gorm.DB, role *Role) (int, error) {
	return createRoleWithAccesses(tx, role)
}

func GetRoleWithAccesses(roleID int) (resp response.ResponseModel, err error) {
	role, err := getRoleWithAccesses(roleID)

	if err != nil {
		response.SetResponseData(&resp, GetRolesWithAccesses{}, "Что-то пошло не так", false, 0, 0, 0)
		return
	}

	response.SetResponseData(&resp, role, "Успех", true, 0, 0, 0)
	return
}

func UpdateRoleWithAccesses(request UpdateRoleWithAccessesReq) (resp response.ResponseModel, err error) {
	tx := database.Postgres().Begin()

	if err = updateRoleWithAccesses(tx, request); err != nil {
		response.SetResponseData(&resp, UpdateRoleWithAccessesReq{}, "Что-то пошло не так", false, 0, 0, 0)
		tx.Rollback()
		return
	}
	tx.Commit()
	response.SetResponseData(&resp, request, "Успех", true, 0, 0, 0)
	return
}
