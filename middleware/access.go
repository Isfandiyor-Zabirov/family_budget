package middleware

import "family_budget/internal/entities/role_accesses"

const (
	FinancialEventCategories = "FINANCIAL_EVENT_CATEGORIES"
	FinancialEvents          = "FINANCIAL_EVENTS"
	Goals                    = "GOALS"
	Roles                    = "ROLES"
	Transactions             = "TRANSACTIONS"
	Users                    = "USERS"
)

// accesses
const (
	CREATE = "CREATE"
	UPDATE = "UPDATE"
	READ   = "READ"
	DELETE = "DELETE"
)

func CheckAccess(accessGroup, access string, userID int) bool {
	return role_accesses.CheckAccess(accessGroup, access, userID)
}
