package helpers

import "github.com/marianop9/valkyrie-migrate/app/repository"

func All[T any](slice []T, predicate func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if !predicate(slice[i]) {
			return false
		}
	}
	return true
}

func Any[T any](slice []T, predicate func(T) bool) bool {
	for i := 0; i < len(slice); i++ {
		if predicate(slice[i]) {
			return true
		}
	}
	return false
}

func FindByName(slice []repository.MigrationGroup, name string) *repository.MigrationGroup {
	for i, mg := range slice {
		if mg.Name == name {
			return &slice[i]
		}
	}

	return nil
}