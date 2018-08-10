package leanix

import (
	"testing"
)

func TestPackageTagSets(t *testing.T) {
	input := [][]string{
		[]string{"a", "b"},
		[]string{"c", "d"},
	}
	expectedOutput := []map[string][]string{
		map[string][]string{
			"tags": []string{"a", "b"},
		},
		map[string][]string{
			"tags": []string{"c", "d"},
		},
	}

	actualOutput := packageTagSets(input)
	assertEqual(t, actualOutput, expectedOutput)
}
