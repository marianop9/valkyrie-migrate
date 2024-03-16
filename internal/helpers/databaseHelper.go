package helpers

import (
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

func GetDb(dsn string) (*sqlx.DB, error) {

	if db, err := sqlx.Open("sqlite3", dsn); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v\n %v", dsn, err.Error())
	} else {
		return db, nil
	}

}

func ConvertDb(db *sql.DB, driver string) *sqlx.DB {
	return sqlx.NewDb(db, driver)
}
