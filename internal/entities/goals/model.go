package goals

import (
	"family_budget/internal/entities/family"
	"gorm.io/gorm"
	"time"
)

const (
	StatusPlanned   = "PLANNED"
	StatusApproved  = "APPROVED"
	StatusRejected  = "REJECTED"
	StatusCancelled = "CANCELLED"
	StatusPending   = "PENDING"
	StatusComplete  = "COMPLETE"
)

type Goals struct {
	ID              int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	FamilyID        int            `gorm:"column:family_id" json:"family_id"` // к какой семье принадлежит
	Name            string         `gorm:"column:name" json:"name"`
	Description     string         `gorm:"column:description" json:"description"`
	TotalBudget     float64        `gorm:"column:total_budget" json:"total_budget"`         // сколько бюджета нужен для выполнение цели
	RemainingBudget float64        `gorm:"column:remaining_budget" json:"remaining_budget"` // сколько бюджета осталось
	Status          string         `gorm:"column:status" json:"status"`
	DueDate         string         `gorm:"column:due_date" json:"due_date"`
	CreatedAt       *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt       *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt       gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK        family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*Goals) TableName() string {
	return "goals"
}

type Filters struct {
	FamilyID    int     `form:"family_id"` // for internal use only
	Search      *string `form:"search"`
	Status      *string `form:"status"`
	DueDateFrom *string `form:"due_date_from"` // YYYY-MM-DD
	DueDateTo   *string `form:"due_date_to"`   // YYYY-MM-DD
	CurrentPage int     `form:"current_page"`
	PageLimit   int     `form:"page_limit"`
}
