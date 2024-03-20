package repository_test

// import (
// 	"database/sql"
// 	"fmt"
// 	"testing"

// 	"github.com/marianop9/valkyrie-migrate/internal/repository"
// 	_ "github.com/mattn/go-sqlite3"
// )

// func TestGetMigrations(t *testing.T) {
// 	repo := repository.NewMigrationRepo(getDb())
// 	if err := repo.EnsureCreated(); err != nil {
// 		t.Fatal(err.Error())
// 	}

// 	mgs, err := repo.GetMigrations()

// 	if err != nil {
// 		t.Fatal(err.Error())
// 	}

// 	if len(mgs) == 0 {
// 		t.Error("found no rows!")
// 	}

// 	t.Logf("found %v rows\n", len(mgs))
// }

// func getDb() *sql.DB {
// 	db, err := sql.Open("sqlite3", "../../test.db")

// 	if err != nil {
// 		panic(fmt.Sprintf("Failed to connect to db: %v\n %s", "", err.Error()))
// 	}

// 	return db

// }