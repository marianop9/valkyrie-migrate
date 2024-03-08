package main

import (
	"database/sql"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/marianop9/valkyrie-migrate/db"
	_ "github.com/mattn/go-sqlite3"
)

const (
	dateFmt = "20060102"
)

func main() {
	folderNames := []string {
		"./migrations",
		"../migrations",
	}
	var dirEntries []os.DirEntry

	for _, folderName := range folderNames {
		entries, err := os.ReadDir(folderName)

		if err == nil && len(entries) > 0 {
			dirEntries = entries
			break
		}
	}

	fmt.Println("files: ", len(dirEntries))

	migrationGroups := getMigrationGroups(dirEntries)

	migrationDb := getDb()

	if err := db.EnsureCreated(migrationDb); err != nil {
		fmt.Println(err.Error())
		fmt.Println("exiting...")
		os.Exit(1)
	}
	
	fmt.Printf("%+v\n", migrationGroups)
}

type MigrationGroup struct {
	Name string;
	Date time.Time
	Files []string
}

func getDb() *sql.DB {
	dsn := "./test.db"
	db, err := sql.Open("sqlite3", dsn)

	if err != nil {
		panic(fmt.Sprintf("Failed to connect to db: %v\n %s", dsn, err.Error()))
	}

	return db

}

func getMigrationGroups(dirEntries []os.DirEntry) []MigrationGroup {
	migrationGroups := make([]MigrationGroup, 0)
	
	for _, entry := range dirEntries {
		fmt.Printf("%+v\n", entry.Name())
		if (entry.IsDir()) {
			entryParts := strings.Split(entry.Name(), "_")

			if (len(entryParts) < 2) {
				panic("folder name doesn't match expected format")
			}

			date, err := time.Parse(dateFmt, entryParts[0]) // yyyymmdd
			if err != nil {
				panic("folder date doesn't match expected format")
			}

			group := MigrationGroup{
				Name: strings.Join(entryParts[1:], "_"),
				Date: date,
			}
			
			migrationGroups = append(migrationGroups, group)
		}
	}

	return migrationGroups
}

