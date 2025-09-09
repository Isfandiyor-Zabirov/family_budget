package role_accesses

import (
	"family_budget/internal/entities/access_groups"
	"family_budget/internal/entities/accesses"
	"family_budget/internal/entities/roles"
	"gorm.io/gorm"
	"time"
)

// RoleAccess - each role accesses to diff routes (access - method, group - entity)
type RoleAccess struct {
	ID            int                       `gorm:"column:id;primary_key;autoIncrement" json:"id,omitempty"`
	RoleID        int                       `gorm:"column:role_id" json:"role_id"`
	AccessID      int                       `gorm:"column:access_id" json:"access_id"`
	AccessGroupID int                       `gorm:"column:access_group_id" json:"access_group_id"`
	Active        bool                      `gorm:"column:active" json:"active"`
	CreatedAt     *time.Time                `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     *time.Time                `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     gorm.DeletedAt            `gorm:"index" json:"-"`
	Role          roles.Role                `gorm:"foreignKey:RoleID" json:"-"`
	Access        accesses.Access           `gorm:"foreignKey:AccessID" json:"-"`
	AccessGroup   access_groups.AccessGroup `gorm:"foreignKey:AccessGroupID" json:"-"`
}

type GetRoleAccessFilter struct {
	OwnerID     *int    `form:"owner_id"`
	RoleID      *int    `form:"role_id"`
	AccessID    *int    `form:"access_id"`
	Name        *string `form:"name"`
	Description *string `form:"description"`
	Page        *int    `form:"page"`
	PageLimit   *int    `form:"page_limit"`
}

type GetRoleAccessResp struct {
	ID          int    `json:"id"`
	Role        string `json:"role"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Active      bool   `json:"active"`
}

type CreateRoleWithAccessesReq struct {
	Role     roles.Role                `json:"role"`
	Accesses []AccessGroupWithAccesses `json:"accesses"`
}

type AccessGroupWithAccesses struct {
	AccessGroupID int   `json:"access_group_id"`
	AccessIDs     []int `json:"access_ids"`
}
