package family

import (
	"gorm.io/gorm"
	"time"
)

type Family struct {
	ID        int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name      string         `gorm:"column:name" json:"name"`
	Phone     string         `gorm:"column:phone" json:"phone"`
	CreatedAt *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
