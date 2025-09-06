package roles

import (
	"family_budget/internal/entities/family"
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	ID          int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	FamilyID    int            `gorm:"column:family_id" json:"family_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK    family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*Roles) TableName() string {
	return "roles"
}
