package financial_events

type FinancialEvent struct {
	ID   int    `gorm:"column:id;primary_key;autoIncrement" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

func (*FinancialEvent) TableName() string {
	return "financial_events"
}
