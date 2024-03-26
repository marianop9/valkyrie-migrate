package valkyrie

import (
	"database/sql"
	"fmt"
	"path"
	"strings"

	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
	postgresRepo "github.com/marianop9/valkyrie-migrate/internal/repository/postgres"
)

func Init(dbName string) error {

	var db *sql.DB
	var err error

	if strings.HasPrefix(dbName, "postgresql://") {
		if db, err = helpers.GetPostgresDb(dbName); err != nil {
			return err
		}

		return postgresRepo.NewMigrationRepo(db).EnsureCreated()
	} else if path.Ext(dbName) == ".db" {
		if db, err = helpers.GetDb(dbName); err != nil {
			return err
		}

		return repository.EnsureCreated(db)
	}

	return fmt.Errorf("invalid database file extension")
}
