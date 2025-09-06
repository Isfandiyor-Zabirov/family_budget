package financial_events

import (
	"family_budget/internal/entities/family"
	"gorm.io/gorm"
	"time"
)

type FinancialEvent struct {
	ID          int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	FamilyID    int            `gorm:"column:family_id" json:"family_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	Inflow      bool           `gorm:"column:inflow" json:"inflow"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK    family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*FinancialEvent) TableName() string {
	return "financial_events"
}
