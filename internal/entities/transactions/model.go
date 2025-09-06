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
	FamilyID int
	UserID   int
}
