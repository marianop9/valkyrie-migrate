package repository_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/jmoiron/sqlx"
	"github.com/marianop9/valkyrie-migrate/valkyrie-migrate/repository"
	_ "github.com/mattn/go-sqlite3"
)

func TestGetMigrations(t *testing.T) {
	repo := repository.NewMigrationRepo(getDb())
	
	mgs, err := repo.GetMigrations()

	if err != nil {
		t.Error(err.Error())
	}

	if len(mgs) == 0 {
		t.Error("found no rows")
	}

	t.Logf("found %v rows\n", len(mgs))
}

func getDb() *sqlx.DB {
	db, err := sqlx.Open("sqlite3", "../../test.db")

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %v\n %s", "", err.Error()))
	}

	return db

}

func TestXxx(t *testing.T) {
	wd, _ := os.Getwd()
	fmt.Println(wd)
}