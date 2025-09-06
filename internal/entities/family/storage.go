package family

import (
	"family_budget/internal/utils/crud"
	"family_budget/pkg/database"
)

func createFamily(fec *Family) (Family, error) {
	repo := crud.NewRepository[Family]()
	db := database.Postgres()
	return repo.Create(db, fec)
}

func updateFamily(fec *Family) (Family, error) {
	repo := crud.NewRepository[Family]()
	db := database.Postgres()
	return repo.Update(db, fec)
}

func deleteFamily(family *Family) error {
	repo := crud.NewRepository[Family]()
	db := database.Postgres()
	return repo.Delete(db, family)
}

func getFamily(id int) (Family, error) {
	repo := crud.NewRepository[Family]()
	db := database.Postgres()
	return repo.Get(db, id)
}

