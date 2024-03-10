package migrations_test

import (
	"os"
	"path"
	"runtime"
	"strings"
	"testing"

	"github.com/marianop9/valkyrie-migrate/valkyrie-migrate/migrations"
)

func getTestDirPath() string {
	wd, _ := os.Getwd()
	if runtime.GOOS == "windows" {
		wd = strings.ReplaceAll(wd, "\\", "/")
	}
	return path.Join(wd, "../../test")
}

func TestGetMigrations2(t *testing.T) {
	testDirBase := getTestDirPath()

	entries, err := os.ReadDir(testDirBase)
	if err != nil {
		t.Errorf("failed to get test directory path: %v", err)
		return
	} else if len(entries) == 0 {
		t.Errorf("test directory '%v' is empty", testDirBase)
		return
	}

	testCases := []struct {
		desc          string
		testDir       string
		expectedErr   bool
		expectedCount int
	}{
		{
			desc:          "2 valid migrations",
			testDir:       "MigrationDir",
			expectedErr:   false,
			expectedCount: 2,
		},
		{
			desc:          "empty dir",
			testDir:       "EmptyDir",
			expectedErr:   false,
			expectedCount: 0,
		},
		{
			desc:          "invalid date fmt",
			testDir:       "InvalidDateFmt",
			expectedErr:   true,
			expectedCount: 0,
		},
		{
			desc:          "invalid date fmt",
			testDir:       "InvalidNameFmt",
			expectedErr:   true,
			expectedCount: 0,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			migrationEntries, err := os.ReadDir(path.Join(testDirBase, tC.testDir))

			if err != nil {
				t.Errorf("failed to read test directory '%s': %v", tC.testDir, err)
				return
			}

			migs, err := migrations.GetMigrationGroups(migrationEntries)

			if err != nil && !tC.expectedErr {
				t.Errorf("unexpected error: %v", err)
				return
			}
			
			if count := len(migs); count != tC.expectedCount {
				t.Errorf("expected %v migration groups, got %v", tC.expectedCount, count)
			}
		})
	}
}
