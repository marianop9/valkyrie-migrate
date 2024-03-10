package main

import (
	"fmt"

	valkyrieMigrate "github.com/marianop9/valkyrie-migrate/valkyrie-migrate"
)

func main() {
	dsn := "./test.db"
	app := valkyrieMigrate.MigrateApp{}
	if err := app.Run(dsn); err != nil {
		fmt.Println(err.Error())
	}
}


