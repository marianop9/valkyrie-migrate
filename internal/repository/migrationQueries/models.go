// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.25.0

package queries

import (
	"time"
)

type Migration struct {
	ID               int64
	MigrationGroupID int64
	Name             string
	ExecutedAt       time.Time
}

type MigrationGroup struct {
	ID   int64
	Name string
}