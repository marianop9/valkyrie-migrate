package models

import (
	"io"
)

type MigrationGroup struct {
	Id             uint
	Name           string
	Files          []io.Reader
	Migrations     []Migration
	MigrationCount int `db:"migrationCount"`
}

func (mg *MigrationGroup) AddFile(f io.Reader) {
	if mg.Files == nil {
		mg.Files = []io.Reader{
			f,
		}
	} else {
		mg.Files = append(mg.Files, f)
	}
}

func (mg *MigrationGroup) AddMigration(mig Migration) {
	if mg.Migrations == nil {
		mg.Migrations = []Migration{
			mig,
		}
	} else {
		mg.Migrations = append(mg.Migrations, mig)
	}
}

type Migration struct {
	Name      string
	GroupName string `db:"groupName"`
	FReader   io.Reader
}

type MigrationStorer interface {
	EnsureCreated() error
	GetMigrations() ([]MigrationGroup, error)
	ExecuteMigrations([]*MigrationGroup) error 
}