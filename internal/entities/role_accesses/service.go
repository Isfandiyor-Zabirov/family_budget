package role_accesses

import (
	"family_budget/internal/entities/access_groups"
	"family_budget/internal/entities/accesses"
	"family_budget/internal/entities/roles"
	"family_budget/internal/utils/response"
	"family_budget/pkg/database"
)

// AssignAccessesToRole assigns accesses and access groups to roles
func AssignAccessesToRole(r roles.Role) (ras []RoleAccess) {

	as := accesses.GetAccessList()
	ags := access_groups.GetAccessGroupList()

	for _, a := range as {
		for _, ag := range ags {
			var ra RoleAccess
			ra.RoleID = r.ID
			ra.AccessID = a.ID
			ra.AccessGroupID = ag.ID

			if r.ID == 1 {
				ra.Active = true
			} else {
				ra.Active = false
			}

			ras = append(ras, ra)
		}
	}
	return ras
}

func CheckAccess(accessGroup, access string, userID int) bool {
	return checkAccess(accessGroup, access, userID)
}

func CreateRoleWithAccesses(req *CreateRoleWithAccessesReq) (resp response.ResponseModel, err error) {

	tx := database.Postgres().Begin()

	roleID, err := roles.CreateRoleWithAccesses(tx, &req.Role)
	if err != nil {
		response.SetResponseData(&resp, struct{}{}, "Не удалось добавить роль", false, 0, 0, 0)
		tx.Rollback()
		return
	}

	accessList := accesses.GetAccessList()

	accessGroupList := access_groups.GetAccessGroupList()

	var ras []RoleAccess

	for _, ag := range accessGroupList {
		for _, a := range accessList {
			var ra RoleAccess
			ra.RoleID = roleID
			ra.AccessID = a.ID
			ra.AccessGroupID = ag.ID
			for _, v := range req.Accesses {
				for _, ad := range v.AccessIDs {
					if a.ID == ad && ag.ID == v.AccessGroupID {
						ra.Active = true
					}

				}
			}
			ras = append(ras, ra)
		}
	}

	if err = createRoleWithAccesses(tx, ras); err != nil {
		response.SetResponseData(&resp, struct{}{}, "Не удалось добавить доступи к роли", false, 0, 0, 0)
		tx.Rollback()
		return
	}

	tx.Commit()

	response.SetResponseData(&resp, struct{}{}, "Роль успешно добавлен", true, 0, 0, 0)
	return
}
