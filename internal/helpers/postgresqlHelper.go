package helpers

import (
	"database/sql"
	"fmt"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func GetPostgresDb(dbUrl string) (*sql.DB, error) {

	if db, err := sql.Open("pgx", dbUrl); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v\n %v", dbUrl, err.Error())
	} else {
		return db, nil
	}
}