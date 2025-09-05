package financial_event_categories

import (
	"family_budget/internal/entities/family"
	"gorm.io/gorm"
	"time"
)

type FinancialEventCategories struct {
	ID          int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	FamilyID    int            `gorm:"column:family_id" json:"family_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	CreatedAt   *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK    family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*FinancialEventCategories) TableName() string {
	return "financial_event_categories"
}

type Filters struct {
	FamilyID    int     `form:"family_id"` // for internal use only
	CurrentPage *uint64 `form:"current_page"`
	PageLimit   *uint64 `form:"page_limit"`
}
