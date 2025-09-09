package roles

import (
	"family_budget/internal/entities/family"
	"time"

	"gorm.io/gorm"
)

type Role struct {
	ID          int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	FamilyID    int            `gorm:"column:family_id" json:"family_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK    family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*Role) TableName() string {
	return "roles"
}

type GetRolesResp struct {
	ID              int    `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
	Name            string `gorm:"column:name" json:"name,omitempty"`
	Description     string `gorm:"column:description" json:"description,omitempty"`
	FamilyID        int    `gorm:"column:family_id" json:"family_id,omitempty"`
	Default         bool   `gorm:"column:default;default:false" json:"default"`
	DefaultStatusID int    `gorm:"column:default_status_id" json:"default_status_id"`
	DefaultStatus   string `gorm:"column:default_status" json:"default_status"`
	StatusColor     string `gorm:"column:status_color" json:"status_color"`
}

type GetRolesFilter struct {
	FamilyID    *int    `form:"family_id"`
	RoleID      *int    `form:"role_id"`
	Search      *string `form:"search"`
	CurrentPage int     `form:"current_page"`
	PageLimit   int     `form:"page_limit"`
}

type GetRoleWithAccessesResp struct {
	RoleID         int              `json:"role_id"`
	RoleName       string           `json:"role_name"`
	RoleFamilyID   int              `json:"role_family_id"`
	AccessByGroups []AccessByGroups `json:"access_by_groups"`
}

type AccessByGroups struct {
	GroupID             int                   `json:"group_id"`
	GroupName           string                `json:"group_name"`
	GroupDescription    string                `json:"group_description"`
	RoleAccessesByGroup []RoleAccessesByGroup `json:"role_accesses_by_group"`
}

type AccessByGroup struct {
	GroupID          int    `json:"group_id"`
	GroupName        string `json:"group_name"`
	GroupDescription string `json:"group_description"`
}

type RoleAccessesByGroup struct {
	ID          int    `json:"id"`
	AccessID    int    `json:"access_id"`
	Active      bool   `json:"active"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type GetRoleResp struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	FamilyID    int    `json:"family_id"`
	Family      string `json:"Family"`
	Default     bool   `json:"default"`
}

type GetRolesWithAccesses struct {
	Role       Role                  `json:"role"`
	AccessList []GetRoleAccessGroups `json:"access_list" json:"accessList,omitempty"`
}

type GetRoleAccessGroups struct {
	AccessGroupID          int               `json:"group_id"`
	AccessGroupCode        string            `json:"group_code"`
	AccessGroupName        string            `json:"group_name"`
	AccessGroupDescription string            `json:"group_description"`
	Accesses               []GetRoleAccesses `json:"accesses"`
}

type GetRoleAccesses struct {
	AccessID          int    `json:"id"`
	AccessCode        string `json:"code"`
	AccessName        string `json:"name"`
	AccessDescription string `json:"description"`
	Active            bool   `json:"active"`
}

type UpdateRoleWithAccessesReq struct {
	Role       Role                      `json:"role"`
	AccessList []AccessGroupWithAccesses `json:"access_list"`
}

type AccessGroupWithAccesses struct {
	AccessGroupID int      `json:"id"`
	Accesses      []Access `json:"accesses"`
}

type Access struct {
	AccessID int  `json:"id"`
	Active   bool `json:"active"`
}
