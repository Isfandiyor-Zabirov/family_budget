package financial_event_types

type FinancialEventType struct{}

func (*FinancialEventType) TableName() string {
	return "financial_event_types"
}
