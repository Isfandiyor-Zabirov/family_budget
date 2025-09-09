package role_accesses

import (
	"errors"
	"family_budget/pkg/database"
	"gorm.io/gorm"
	"log"
)

func checkAccess(accessGroup, access string, userID int) bool {
	var active bool
	sqlQuery := `select ra.active
					from role_accesses ra, roles r, access_groups ag, accesses a, users u
					where u.deleted_at is null and ra.access_id = a.id and ra.access_group_id = ag.id and r.id = ra.role_id and u.role_id = r.id
					and u.id = ? and ag.code = ? and a.code = ?`
	if err := database.Postgres().Raw(sqlQuery, userID, accessGroup, access).Row().Scan(&active); err != nil {
		log.Println("checkAccess func query error:", err.Error())
		return false
	}
	return active
}

func createRoleWithAccesses(tx *gorm.DB, ras []RoleAccess) error {
	if err := tx.Create(&ras).Error; err != nil {
		log.Println("createRoleWithAccesses func create role accesses query error:", err.Error())
		return errors.New("Ошибка сервера ")
	}
	return nil
}
