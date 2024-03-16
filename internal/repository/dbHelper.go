package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

func EnsureCreated(db *sqlx.DB) error {
	migrationTables := []string{
		"migrationGroup",
		"migration",
	}

	query := `SELECT name 
		FROM sqlite_master 
		WHERE type='table' 
			AND name IN ($1, $2)`

	rows, err := db.Query(query, migrationTables[0], migrationTables[1])
	if err != nil {
		return err
	}

	defer rows.Close()

	foundTables := make([]string, 0)
	for rows.Next() {
		var name string
        if err := rows.Scan(&name); err != nil {
            return err
        }
        
		foundTables = append(foundTables, name)
    }

	if len(foundTables) == len(migrationTables) {
		fmt.Println("migrations tables exist")
		return nil
	}

	// create migrationGroup table
	if !sliceContains(foundTables, migrationTables[0]) {
		if err = createMigrationGroupTable(db); err != nil {
			return err
		}
	}

	// create migration table
	if !sliceContains(foundTables, migrationTables[1]) {
		if err = createMigrationTable(db); err != nil {
			return err
		}
	}

	return nil
}

func sliceContains(slice []string, s string) bool {
	for _, ss := range slice {
		if ss == s {
			return true
		}
	}

	return false
}

func createMigrationGroupTable(db *sqlx.DB) error {
	fmt.Println("creating table 'migrationGroup'...")
	
	cmd := `CREATE TABLE "migrationGroup" (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name VARCHAR(255) NOT NULL
	)`

	if _, sqlErr := db.Exec(cmd); sqlErr != nil {
		return sqlErr
	}
	
	return nil
}

func createMigrationTable(db *sqlx.DB) error {
	fmt.Println(`creating table 'migration'...`)

	cmd := `CREATE TABLE "migration" (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		migration_group_id INTEGER,
		name VARCHAR(255) NOT NULL,
		executed_at TIMESTAMP NOT NULL
	)`
	
	if _, sqlErr := db.Exec(cmd); sqlErr != nil {
		return sqlErr
	}
	
	return nil
}