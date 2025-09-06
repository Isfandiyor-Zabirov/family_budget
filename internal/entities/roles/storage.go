package roles

import (
	"errors"
	"family_budget/pkg/database"
	"gorm.io/gorm"
	"log"
	"time"
)

func getRoleDB(id int) (role *Roles, err error) {
	err = database.Postgres().Find(&role, id).Error
	if err != nil {
		log.Println("getRole err: ", err.Error())
		return role, err
	}
	return role, nil
}

func createRoleDB(role *Roles, tx *gorm.DB) (Roles, error) {
	err := tx.Create(&role).Error
	if err != nil {
		log.Println("createRoleDB err:", err.Error())
		return *role, err
	}
	return *role, nil
}

func updateRoleDB(role *Roles) error {
	err := database.Postgres().Updates(&role).Error
	if err != nil {
		log.Println("updateRoleDB err:", err.Error())
		return err
	}
	return nil
}

func deleteRoleDB(role *Roles) error {
	err := database.Postgres().Delete(&role, role.ID).Error
	if err != nil {
		log.Println("deleteRoleDB err:", err.Error())
		return err
	}
	return nil
}

func getRolesDB(limit, page int, query *gorm.DB) (roles *[]GetRolesResp, totalRows int64, err error) {
	err = query.Count(&totalRows).Error
	if err != nil {
		log.Println("failed to count roles", err.Error())
		return nil, 0, err
	}
	err = query.Select(`r.*, ts.description as default_status, ts.color as status_color`).Offset(limit * (page - 1)).Limit(limit).Order("r.name").Scan(&roles).Error
	if err != nil {
		log.Println("db find roles err:", err.Error())
		return nil, 0, err
	}
	return roles, totalRows, nil
}

func roleExistsDB(role *Roles) bool {
	count := 0
	err := database.Postgres().Raw("SELECT count(*) FROM roles WHERE name = ? and family_id = ?", role.Name, role.FamilyID).Scan(&count).Error
	if err != nil {
		log.Println("check role existence err: ", err.Error())
		return true
	}
	if count > 0 {
		return true
	}
	return false
}

func getRoleWithAccessesDB(r Roles) (resp GetRoleWithAccessesResp, err error) {
	log.Println("role: ", r)
	resp.RoleID = r.ID
	resp.RoleName = r.Name
	resp.RoleOwnerID = r.FamilyID
	var groups []AccessByGroup
	err = database.Postgres().Raw(`SELECT ag.id group_id, ag.name group_name, ag.description group_description FROM access_groups ag order by group_id`).Scan(&groups).Error
	if err != nil {
		return
	}

	var accessData []RoleAccessesByGroup
	for _, v := range groups {
		err = database.Postgres().Raw(`SELECT ra.id, ra.role_id, ra.active, ra.access_id, a.name, a.description FROM role_accesses ra, roles r, accesses a
	WHERE ra.role_id = r.id and a.id = ra.access_id and r.id = ?`, resp.RoleID).
			Scan(&accessData).Error
		if err != nil {
			return
		}
		var group AccessByGroups
		group.GroupID = v.GroupID
		group.GroupDescription = v.GroupDescription
		group.GroupName = v.GroupName
		group.RoleAccessesByGroup = nil
		resp.AccessByGroups = append(resp.AccessByGroups, group)
	}
	return
}

func createRoleWithAccesses(tx *gorm.DB, role *Roles) (int, error) {
	if err := tx.Create(&role).Error; err != nil {
		log.Println("createRoleWithAccesses func create role query error:", err.Error())
		return 0, errors.New("Ошибка сервера ")
	}

	return role.ID, nil
}

func getRoleWithAccesses(roleID int) (resp GetRolesWithAccesses, err error) {
	roleInfoSql := `select r.id, r.name as name, r.description, r.family_id from roles r where  r.id = ?`

	if err = database.Postgres().Raw(roleInfoSql, roleID).Scan(&resp.Role).Error; err != nil {
		log.Println("getRoleWithAccesses func role info query error:", err.Error())
		return GetRolesWithAccesses{}, errors.New("Ошибка сервера ")
	}

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
