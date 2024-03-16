package main

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/marianop9/valkyrie-migrate/internal/app"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
	_ "github.com/mattn/go-sqlite3"
)

func main() {
	dsn := "./test.db"

	migrationDb := getDb(dsn)

	migrationRepo := repository.NewMigrationRepo(migrationDb)
	
	app := app.NewMigrateApp(migrationRepo)
	
	if err := app.Run(dsn); err != nil {
		fmt.Println(err.Error())
	}
}

func getDb(dsn string) *sqlx.DB {
	db, err := sqlx.Open("sqlite3", dsn)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %v\n %s", dsn, err.Error()))
	}

	return db

}
