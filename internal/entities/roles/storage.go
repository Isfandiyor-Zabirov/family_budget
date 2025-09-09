package roles

import (
	"errors"
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
	"gorm.io/gorm"
	"log"
	"time"
)

func getRole(id int) (Role, error) {
	repo := crud.NewRepository[Role]()
	db := database.Postgres()
	return repo.Get(db, id)
}

func createRole(role *Role, tx *gorm.DB) (Role, error) {
	repo := crud.NewRepository[Role]()
	//	db := database.Postgres()
	return repo.Create(tx, role)
}

func updateRole(role *Role) (Role, error) {
	repo := crud.NewRepository[Role]()
	db := database.Postgres()
	return repo.Create(db, role)
}

func deleteRole(role *Role) error {
	repo := crud.NewRepository[Role]()
	db := database.Postgres()
	return repo.Delete(db, role)
}

func getRoles(filter GetRolesFilter) (roles *[]GetRolesResp, totalRows int64, err error) {
	roles = &[]GetRolesResp{}
	query := database.Postgres().Table("roles r").Where("r.family_id = ? and r.deleted_at is null", *filter.FamilyID)

	if filter.RoleID != nil {
		query = query.Where("r.id = ?", *filter.RoleID)
	}

	if filter.Search != nil {
		query = query.Where("r.name ilike ? or r.description ilike ?", "%"+*filter.Search+"%", "%"+*filter.Search+"%")
	}

	if filter.CurrentPage != 0 {
		filter.CurrentPage = 1
	}

	if filter.PageLimit == 0 {
		filter.PageLimit = 20
	}

	err = query.Count(&totalRows).Error
	if err != nil {
		log.Println("failed to count roles", err.Error())
		return nil, 0, err
	}
	err = query.Select(`r.*`).Offset(filter.PageLimit * (filter.CurrentPage - 1)).
		Limit(filter.PageLimit).Order("r.name").Scan(&roles).Error
	if err != nil {
		log.Println("getRoles func query error:", err.Error())
		return nil, 0, err
	}
	return roles, totalRows, nil
}

func createRoleWithAccesses(tx *gorm.DB, role *Role) (int, error) {
	if err := tx.Create(&role).Error; err != nil {
		log.Println("createRoleWithAccesses func create role query error:", err.Error())
		return 0, errors.New("Ошибка сервера ")
	}

	return role.ID, nil
}

func getRoleWithAccesses(roleID int) (resp GetRolesWithAccesses, err error) {
	roleInfoSql := `select r.id, r.name as name, r.description, r.family_id from roles r where r.id = ?`

	if err = database.Postgres().Raw(roleInfoSql, roleID).Scan(&resp.Role).Error; err != nil {
		log.Println("getRoleWithAccesses func role info query error:", err.Error())
		return GetRolesWithAccesses{}, errors.New("Ошибка сервера ")
	}

	resp.Role, _ = getRole(roleID)

	var accessGroup []struct {
		ID          int
		Name        string
		Code        string
		Description string
	}

	if err = database.Postgres().Raw("select * from access_groups").Scan(&accessGroup).Error; err != nil {
		log.Println("getRoleWithAccesses func access groups info query error:", err.Error())
		return GetRolesWithAccesses{}, errors.New("Ошибка сервера ")
	}

	for _, row := range accessGroup {
		var accesses GetRoleAccessGroups
		accesses.AccessGroupID = row.ID
		accesses.AccessGroupCode = row.Code
		accesses.AccessGroupName = row.Name
		accesses.AccessGroupDescription = row.Description

		sqlQuery := `select a.id as access_id, a.code as access_code, a.name as access_name, a.description as access_description, ra.active 
						from accesses a, access_groups ag, role_accesses ra 
						where ra.access_id = a.id and ra.access_group_id = ag.id and ag.id = ? and ra.role_id = ? order by access_id`

		if err = database.Postgres().Raw(sqlQuery, row.ID, roleID).Scan(&accesses.Accesses).Error; err != nil {
			log.Println("getRoleWithAccesses func accesses info query error:", err.Error())
			return GetRolesWithAccesses{}, errors.New("Ошибка сервера ")
		}

		resp.AccessList = append(resp.AccessList, accesses)
	}

	return
}

func updateRoleWithAccesses(tx *gorm.DB, req UpdateRoleWithAccessesReq) error {

	if err := tx.Updates(req.Role).Error; err != nil {
		log.Println("updateRoleWithAccesses func update role query error:", err.Error())
		return errors.New("Ошибка сервера ")
	}

	sqlQuery := `update role_accesses set active = ?, updated_at = ? where access_group_id = ? and access_id = ? and role_id = ?`

	for _, row := range req.AccessList {
		for _, access := range row.Accesses {
			err := tx.Exec(sqlQuery, access.Active, time.Now(), row.AccessGroupID, access.AccessID, req.Role.ID).Error
			if err != nil {
				log.Println("updateRoleWithAccesses func update role query error:", err.Error())
				return errors.New("Ошибка сервера ")
			}
		}
	}

	return nil
}
