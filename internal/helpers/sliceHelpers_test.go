package helpers_test

import (
	"testing"

	"github.com/marianop9/valkyrie-migrate/internal/helpers"
	"github.com/marianop9/valkyrie-migrate/internal/repository"
)

func nonEmptyStartWithA(s string) bool {
	return len(s) > 0 && s[0] == 'a'
}

func TestAll(t *testing.T) {
	

	testCases := []struct {
		desc           string
		slice          []string
		expectedResult bool
	}{
		{
			desc: "all true",
			slice: []string{
				"abc",
				"acdc",
				"abds",
			},
			expectedResult: true,
		},
		{
			desc: "some true",
			slice: []string{
				"abc",
				"bcdc",
				"abds",
			},
			expectedResult: false,
		},
		{
			desc: "none true",
			slice: []string{
				"dbc",
				"bcdc",
				"cbds",
			},
			expectedResult: false,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := helpers.All[string](tC.slice, nonEmptyStartWithA)

			if result != tC.expectedResult {
				t.Errorf("%s - expected '%v', got '%v'", tC.desc, tC.expectedResult, result)
			}
		})
	}
}

func TestAny(t *testing.T) {
	testCases := []struct {
		desc           string
		slice          []string
		expectedResult bool
	}{
		{
			desc: "one matches",
			slice: []string{
				"abc",
				"bcdc",
				"bbds",
			},
			expectedResult: true,
		},
		{
			desc: "none match",
			slice: []string{
				"xbc",
				"bcdc",
				"xbds",
			},
			expectedResult: false,
		},
		{
			desc: "more than one match",
			slice: []string{
				"abc",
				"acdc",
				"cbds",
			},
			expectedResult: true,
		},
	}

	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := helpers.Any[string](tC.slice, nonEmptyStartWithA)

			if result != tC.expectedResult {
				t.Errorf("%s - expected '%v', got '%v'", tC.desc, tC.expectedResult, result)
			}
		})
	}
}

func TestFindByName(t *testing.T) {
	slice := []repository.MigrationGroup {
		{
			Name: "name123",
		},
		{
			Name: "anothername123",
		},
		{
			Name: "anothertestname123",
		},
	}
	
	testCases := []struct {
		desc	string
		targetName string
		expectedNil bool
	}{
		{
			desc: "existing name",
			targetName: "anothername123",
			expectedNil: false,			
		},
		{
			desc: "non-existing name",
			targetName: "nonexisting123",
			expectedNil: true,			
		},
	}
	for _, tC := range testCases {
		t.Run(tC.desc, func(t *testing.T) {
			result := helpers.FindMigrationGroup(slice, tC.targetName)

			if (result == nil && !tC.expectedNil) || (result != nil && tC.expectedNil) {
				t.Errorf("%s - expected nil: '%v', got '%v'", tC.desc, tC.expectedNil, result)
			}
		})
	}
}