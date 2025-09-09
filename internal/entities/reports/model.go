package reports

type MainReport struct {
	TotalBudget  float64 `json:"total_budget"`
	TotalExpense float64 `json:"total_expense"`
	TotalIncome  float64 `json:"total_income"`
}

type GraphReport struct {
	Date     string  `json:"date"`
	Expenses float64 `json:"expenses"`
	Incomes  float64 `json:"incomes"`
}

type Filter struct {
	FamilyID int    `json:"family_id"`
	DateFrom string `json:"date_from"`
	DateTo   string `json:"date_to"`
}
