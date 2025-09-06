package role_accesses

import (
	"family_budget/internal/entities/access_groups"
	"family_budget/internal/entities/accesses"
	"family_budget/internal/entities/roles"
)

// AssignAccessesToRole assigns accesses and access groups to roles
func AssignAccessesToRole(r roles.Roles) (ras []RoleAccess) {

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
