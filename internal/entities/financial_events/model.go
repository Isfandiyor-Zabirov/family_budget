package financial_events

import (
	"family_budget/internal/entities/family"
	"family_budget/internal/entities/financial_event_categories"
	"time"

	"gorm.io/gorm"
)

type FinancialEvent struct {
	ID              int                                                 `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	EventCategoryID int                                                 `gorm:"column:category_id" json:"category_id"`
	FamilyID        int                                                 `gorm:"column:family_id" json:"family_id"`
	Name            string                                              `gorm:"column:name" json:"name"`
	Description     string                                              `gorm:"column:description" json:"description"`
	Inflow          bool                                                `gorm:"column:inflow" json:"inflow"`
	CreatedAt       *time.Time                                          `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       *time.Time                                          `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt                                      `gorm:"index" json:"-"`
	FamilyFK        family.Family                                       `gorm:"foreignKey:FamilyID" json:"-"`
	CategoryFK      financial_event_categories.FinancialEventCategories `gorm:"foreignKey:EventCategoryID" json:"-"`
}

func (*FinancialEvent) TableName() string {
	return "financial_events"
}

type Filters struct {
	Search      *string `form:"search"`
	CategoryID  *int    `form:"category_id"`
	Inflow      *bool   `form:"inflow"`
	FamilyID    int     `form:"family_id"` // for internal use only
	CurrentPage int     `form:"current_page"`
	PageLimit   int     `form:"page_limit"`
}
