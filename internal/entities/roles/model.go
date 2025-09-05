package roles

import (
	"time"

	"gorm.io/gorm"
)

type Roles struct {
	ID          int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
}

func (*Roles) TableName() string {
	return "roles"
}
