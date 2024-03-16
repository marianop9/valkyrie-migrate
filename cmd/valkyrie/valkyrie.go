package main

import (
	"fmt"

	"github.com/marianop9/valkyrie-migrate/pkg/cmd"
)

func main() {
	err := cmd.NewValkyrieCmd().Execute()

	if err != nil {
		fmt.Printf("command failed:\n %v\n", err)
	}
}