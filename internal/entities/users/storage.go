package users

import (
	"errors"
	"family_budget/pkg/database"
	"log"
)

func getMe(userID int) (me Me, err error) {
	userSqlQuery := `select u.id, u.name, u.surname, u.middle_name, u.surname || ' ' || u.name || ' ' || u.middle_name as full_name, u.login, u.phone,
						r.name as role, u.role_id, u.family_id, u.email
						from users u, roles r
						where u.role_id = r.id and u.id = ?`
	if err = database.Postgres().Raw(userSqlQuery, userID).Scan(&me.UserData).Error; err != nil {
		log.Println("getMe func user data query error:", err.Error())
		return Me{}, errors.New("Ошибка сервера ")
	}

	accessGroupSql := `select ag.id as access_group_id, ag.code as access_group_code, ag.name as access_group_name from access_groups ag`

	var accessGroupInfo []struct {
		AccessGroupID   int    `json:"access_group_id"`
		AccessGroupCode string `json:"access_group_code"`
		AccessGroupName string `json:"access_group_name"`
	}

	if err = database.Postgres().Raw(accessGroupSql).Scan(&accessGroupInfo).Error; err != nil {
		log.Println("getMe func accessGroupSql query error:", err.Error())
		return Me{}, errors.New("Ошибка сервера ")
	}

	accessSqlQuery := `select 
    					a.id as access_id, 
    					a.code as access_code, 
    					ra.active as active, 
    					a.name as access_name
						from accesses a, access_groups ag, role_accesses ra, roles r, users u
						where u.role_id = r.id and ra.access_id = a.id and ra.access_group_id = ag.id and ra.role_id = r.id and ag.id = ? and u.id = ?`

	for _, row := range accessGroupInfo {
		var accesses MeAccessGroup
		accesses.AccessGroupID = row.AccessGroupID
		accesses.AccessGroupCode = row.AccessGroupCode
		accesses.AccessGroupName = row.AccessGroupName
		if err = database.Postgres().Raw(accessSqlQuery, row.AccessGroupID, userID).Scan(&accesses.Accesses).Error; err != nil {
			log.Println("getMe func access list query error:", err.Error())
			return Me{}, errors.New("Ошибка сервера ")
		}
		me.AccessList = append(me.AccessList, accesses)
	}

	return
}
