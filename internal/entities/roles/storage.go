package roles

import (
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
)

func createRole(r *Roles) (Roles, error) {
	repo := crud.NewRepository[Roles]()
	db := database.Postgres()
	return repo.Create(db, r)
}

func getRole(id int) (Roles, error) {
	repo := crud.NewRepository[Roles]()
	db := database.Postgres()
	return repo.Get(db, id)
}

func updateRole(r *Roles) (Roles, error) {
	repo := crud.NewRepository[Roles]()
	db := database.Postgres()
	return repo.Update(db, r)
}

func deleteRole(r *Roles) error {
	repo := crud.NewRepository[Roles]()
	db := database.Postgres()
	return repo.Delete(db, r)
}
