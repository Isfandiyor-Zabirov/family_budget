package financial_event_types

type FinancialEventType struct {
	ID          int    `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name        string `gorm:"column:name" json:"name"`
	Description string `gorm:"column:description" json:"description"`
}

func (*FinancialEventType) TableName() string {
	return "financial_event_types"
}
