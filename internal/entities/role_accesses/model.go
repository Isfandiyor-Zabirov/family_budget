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
	Role          roles.Roles               `gorm:"foreignKey:RoleID" json:"-"`
	Access        accesses.Access           `gorm:"foreignKey:AccessID" json:"-"`
	AccessGroup   access_groups.AccessGroup `gorm:"foreignKey:AccessGroupID" json:"-"`
}
