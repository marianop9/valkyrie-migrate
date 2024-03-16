package valkyrie

import (
	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
)

func Init(dbName string) error {

	db, err := helpers.GetDb(dbName)
	if err != nil {
		return err
	}

	return repository.EnsureCreated(db)
}