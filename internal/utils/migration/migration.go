package migration

import (
	"family_budget/internal/entities/access_groups"
	"family_budget/internal/entities/accesses"
	"family_budget/internal/entities/family"
	"family_budget/internal/entities/financial_event_categories"
	"family_budget/internal/entities/financial_events"
	"family_budget/internal/entities/goals"
	"family_budget/internal/entities/role_accesses"
	"family_budget/internal/entities/roles"
	"family_budget/internal/entities/transactions"
	"family_budget/internal/entities/users"
	"family_budget/pkg/database"
	"fmt"
	"log"
)

func AutoMigrate() {
	fmt.Println("Automatically migrating the schemas...")

	err := database.Postgres().AutoMigrate(
		&family.Family{},
		&roles.Role{},
		&users.User{},
		&accesses.Access{},
		&access_groups.AccessGroup{},
		&role_accesses.RoleAccess{},
		&financial_event_categories.FinancialEventCategories{},
		&financial_events.FinancialEvent{},
		&goals.Goals{},
		&transactions.Transactions{},
	)
	if err != nil {
		log.Println("AutoMigrate func error: ", err.Error())
		return
	}

	initDml()
	fmt.Println("Finished migrating the schemas...")
}
