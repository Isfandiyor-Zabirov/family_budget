package roles

import (
	"errors"
	"family_budget/pkg/database"
	"fmt"
	"gorm.io/gorm"
)

func GetRole(roleID int) (*Roles, error) {
	u, err := getRoleDB(roleID)
	if err != nil {
		return u, err
	}
	return u, nil
}

func CreateRole(role *Roles, tx *gorm.DB) (Roles, error) {
	//повторение названий ролей учитывается для текущего владельцы

	if roleExistsDB(role) {
		return Roles{}, errors.New("Роль уже существует")
	}
	r, err := createRoleDB(role, tx)
	if err != nil {
		return r, err
	}
	return r, nil
}

func UpdateRole(r Roles) error {
	//проверка на существование
	existingRole, err := getRoleDB(r.ID)
	if existingRole == nil {
		return errors.New("Role does not exist ")
	}

	if existingRole.ID == 0 {
		return err
	}
	r.FamilyID = existingRole.FamilyID
	//возможно нужно проверять значения и менять по необходимости (или роли, доступам к определенному полю)
	return updateRoleDB(&r)
}

func DeleteRole(roleID, ownerID int) (*Roles, error) {
	r, err := getRoleDB(roleID)
	if err != nil {
		return r, err
	}
	if r.ID == 0 {
		return r, errors.New("роль уже удалена")
	}
	if r.FamilyID != ownerID {
		return r, errors.New("нет доступа к чужим данным")
	}
	err = deleteRoleDB(r)
	if err != nil {
		return r, err
	}
	return r, nil
}

func GetRoles(filter GetRolesFilter) (*[]GetRolesResp, int, int64, error) {
	query := database.Postgres().Table("roles r").Where("r.owner_id = ? and r.deleted_at is null", *filter.OwnerID).
		Joins(`left join transaction_statuses ts on ts.id = r.default_status_id`)
	//filters with arrays (multiple choice)

	if filter.RoleID != nil {
		query = query.Where("r.id = ?", *filter.RoleID)
	}
	//filters with text part search
	if filter.Search != nil {
		query = query.Where("r.name ilike ? or r.description ilike ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
	}
	//sort by
	if filter.OrderBy != nil {
		if filter.OrderDescending == true {
			query = query.Order(fmt.Sprintf("r.%v desc", *filter.OrderBy))
		} else {
			query = query.Order(fmt.Sprintf("r.%v", *filter.OrderBy))
		}
	}
	//pagination
	var page, pageLimit int
	if filter.Page != nil {
		page = *filter.Page
	} else {
		page = 1
	}
	if filter.PageLimit != nil {
		pageLimit = *filter.PageLimit
	} else {
		pageLimit = 15
	}
	//system filters, depending on roles, accesses etc.

	//нуждается в доработке

	roleList, totalRows, err := getRolesDB(pageLimit, page, query)
	if err != nil {
		return nil, 0, 0, err
	}
	return roleList, page, totalRows, nil
}

func GetRoleWithAccesses(roleID int) (resp GetRoleWithAccessesResp, err error) {
	r, err := getRoleDB(roleID)
	if err != nil {
		return
	}
	return getRoleWithAccessesDB(*r)
}

func CreateRoleWithAccesses(tx *gorm.DB, role *Roles) (int, error) {
	return createRoleWithAccesses(tx, role)
}

func GetRoleWithAccessesV2(roleID int) (resp GetRolesWithAccesses, err error) {
	return getRoleWithAccesses(roleID)
}

func UpdateRoleWithAccesses(request UpdateRoleWithAccessesReq) error {
	tx := database.Postgres().Begin()

	if err := updateRoleWithAccesses(tx, request); err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}
