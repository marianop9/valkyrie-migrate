package repository

import (
	"database/sql"
	"fmt"
	"os"
)

func EnsureCreated(db *sql.DB) error {
	migrationTables := []string{
		"migration_group",
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

	// create migration_group table
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

func createMigrationGroupTable(db *sql.DB) error {
	fmt.Println("creating table 'migration_group'...")

	buf,err := os.ReadFile("./db/schema/sqlite/cr_migrationGroup.sql")
	if err != nil {
		return err
	}

	if _, sqlErr := db.Exec(string(buf)); sqlErr != nil {
		return sqlErr
	}
	
	return nil
}

func createMigrationTable(db *sql.DB) error {
	fmt.Println(`creating table 'migration'...`)

	buf,err := os.ReadFile("./db/schema/sqlite/cr_migration.sql")
	if err != nil {
		return err
	}
	
	if _, sqlErr := db.Exec(string(buf)); sqlErr != nil {
		return sqlErr
	}
	
	return nil
}