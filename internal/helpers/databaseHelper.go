package helpers

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

func GetDb(dsn string) (*sql.DB, error) {

	if db, err := sql.Open("sqlite3", dsn); err != nil {
		return nil, fmt.Errorf("failed to connect to db: %v\n %v", dsn, err.Error())
	} else {
		return db, nil
	}

}

