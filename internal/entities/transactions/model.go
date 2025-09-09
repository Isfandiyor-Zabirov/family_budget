package transactions

import (
	"family_budget/internal/entities/family"
	"gorm.io/gorm"
	"time"
)

type Transactions struct {
	ID               int            `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	FamilyID         int            `gorm:"column:family_id" json:"family_id"` // к какой семье принадлежит
	UserID           int            `gorm:"column:user_id" json:"user_id"`     // кто создал операцию
	FinancialEventID int            `gorm:"column:financial_event_id" json:"financial_event_id"`
	GoalID           int            `gorm:"column:goal_id" json:"goal_id"`
	Amount           float64        `gorm:"column:amount" json:"amount"`
	Description      string         `gorm:"column:description" json:"description"`
	CreatedAt        *time.Time     `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt        *time.Time     `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt        gorm.DeletedAt `gorm:"index" json:"-"`
	FamilyFK         family.Family  `gorm:"foreignKey:FamilyID" json:"-"`
}

func (*Transactions) TableName() string {
	return "transactions"
}

type Filters struct {
	FamilyID         int     // for internal use only
	Search           *string `form:"search"`
	UserID           *int    `form:"user_id"`
	FinancialEventID *int    `form:"financial_event_id"`
	GoalID           *int    `form:"goal_id"`
	DateFrom         *string `form:"date_from"`
	DateTo           *string `form:"date_to"`
	CurrentPage      int     `form:"current_page"`
	PageLimit        int     `form:"page_limit"`
}

type Response struct {
	Transactions
	UserName           string `json:"user_name"`
	FinancialEventName string `json:"financial_event_name"`
	GoalName           string `json:"goal_name"`
}
