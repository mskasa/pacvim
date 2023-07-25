package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateFiles(t *testing.T) {
	const validateTestMapPath = "files/test/validate/"
	cases := map[string]struct {
		mapFileName string
		expected    string
	}{
		"normal lower limit": {
			"normal_lower_limit.txt",
			"",
		},
		"normal upper limit": {
			"normal_upper_limit.txt",
			"",
		},
		"error columns exceeded": {
			"error_columns_exceeded.txt",
			"Stage Map Validation Error: files/test/validate/error_columns_exceeded.txt; Please keep the stage within 50 columns;",
		},
		"error lines exceeded": {
			"error_lines_exceeded.txt",
			"Stage Map Validation Error: files/test/validate/error_lines_exceeded.txt; Please keep the stage within 20 lines;",
		},
		"error no boundaries": {
			"error_no_boundaries.txt",
			"Stage Map Validation Error: files/test/validate/error_no_boundaries.txt; Create a boundary for the stage map with '+' (line 1,5,10,15);",
		},
		"error uneven length": {
			"error_uneven_length.txt",
			"Stage Map Validation Error: files/test/validate/error_uneven_length.txt; Make the width of the stage map uniform (line 5,10);",
		},
		"error invalid mime type": {
			"error_invalid_mime_type.txt",
			"MIME Type Validation Error: files/test/validate/error_invalid_mime_type.txt; Invalid mime type: application/octet-stream;",
		},
		"error file does not exist": {
			"error_file_does_not_exist.txt",
			"open files/test/validate/error_file_does_not_exist.txt: file does not exist",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			stages := []stage{
				{mapPath: validateTestMapPath + tt.mapFileName},
			}
			if result := validateFiles(stages); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateActualFiles(t *testing.T) {
	stages := initStages()
	if err := validateFiles(stages); err != nil {
		t.Error(err)
	}
}

func TestSplitStages(t *testing.T) {
	cases := map[string]struct {
		level         int
		expectedLevel int
		stages        []stage
	}{
		"head": {
			level:         1,
			expectedLevel: 1,
			stages: []stage{
				{level: 1},
				{level: 2},
			},
		},
		"middle": {
			level:         2,
			expectedLevel: 2,
			stages: []stage{
				{level: 1},
				{level: 2},
				{level: 3},
			},
		},
		"default": {
			level:         3,
			expectedLevel: 1,
			stages: []stage{
				{level: 1},
				{level: 2},
			},
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			stages := splitStages(tt.stages, &tt.level)
			if stages[0].level != tt.expectedLevel {
				t.Errorf("expected %d but %d", tt.expectedLevel, stages[0].level)
			}
		})
	}
}
