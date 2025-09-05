package financial_events

type FinancialEvent struct{}

func (*FinancialEvent) TableName() string {
	return "financial_events"
}
