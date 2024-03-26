package helpers

import "github.com/marianop9/valkyrie-migrate/internal/models"

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

func FindMigrationGroup(slice []models.MigrationGroup, name string) *models.MigrationGroup {
	for i, mg := range slice {
		if mg.Name == name {
			return &slice[i]
		}
	}

	return nil
}

func FindMigration(migs []models.Migration, name string) *models.Migration {
	for i, m := range migs {
		if m.Name == name {
			return &migs[i]
		}
	}

	return nil
}

