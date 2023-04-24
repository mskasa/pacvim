package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateMimeType(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal/map01": {
			mapPath:  "files/test/validateMimeType/map01.txt",
			expected: "",
		},
		"error/map02": {
			mapPath:  "files/test/validateMimeType/map02.txt",
			expected: "MIME Type Validation Error: files/test/validateMimeType/map02.txt; Invalid mime type: application/octet-stream;",
		},
		"error/map03": {
			mapPath:  "files/test/validateMimeType/map03.txt",
			expected: "open files/test/validateMimeType/map03.txt: file does not exist",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if result := validateMimeType(tt.mapPath); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateFileSize(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal/map01": {
			mapPath:  "files/test/validateFileSize/map01.txt",
			expected: "",
		},
		"error/map02": {
			mapPath:  "files/test/validateFileSize/map02.txt",
			expected: "File Size Validation Error: files/test/validateFileSize/map02.txt; File size exceeded: 1049 (Max file size is 1024);",
		},
		"error/map03": {
			mapPath:  "files/test/validateFileSize/map03.txt",
			expected: "open files/test/validateFileSize/map03.txt: file does not exist",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if result := validateFileSize(tt.mapPath); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateStageMap(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal/map01": {
			mapPath:  "files/test/validateStageMap/map01.txt",
			expected: "",
		},
		"normal/map02": {
			mapPath:  "files/test/validateStageMap/map02.txt",
			expected: "",
		},
		"error/map03": {
			mapPath:  "files/test/validateStageMap/map03.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map03.txt; Please keep the stage within 50 columns;",
		},
		"error/map04": {
			mapPath:  "files/test/validateStageMap/map04.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map04.txt; Please keep the stage within 20 lines;",
		},
		"error/map05": {
			mapPath:  "files/test/validateStageMap/map05.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map05.txt; Create a boundary for the stage map with '+' (line 1,15);",
		},
		"error/map06": {
			mapPath:  "files/test/validateStageMap/map06.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map06.txt; Make the width of the stage map uniform (line 5,10);",
		},
		"error/map07": {
			mapPath:  "files/test/validateStageMap/map07.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map07.txt; Create a boundary for the stage map with '+' (line 5,10);",
		},
		"error/map08": {
			mapPath:  "files/test/validateStageMap/map08.txt",
			expected: "open files/test/validateStageMap/map08.txt: file does not exist",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if result := validateStageMap(tt.mapPath); result != nil {
				assert.EqualErrorf(t, result, tt.expected, "Error should be: %v, got: %v", tt.expected, result)
			}
		})
	}
}

func TestValidateFiles(t *testing.T) {
	cases := map[string]struct {
		mapPath  string
		expected string
	}{
		"normal": {
			mapPath:  "files/test/validateMimeType/map01.txt",
			expected: "",
		},
		"error/validateMimeType": {
			mapPath:  "files/test/validateMimeType/map02.txt",
			expected: "MIME Type Validation Error: files/test/validateMimeType/map02.txt; Invalid mime type: application/octet-stream;",
		},
		"error/validateFileSize": {
			mapPath:  "files/test/validateFileSize/map02.txt",
			expected: "File Size Validation Error: files/test/validateFileSize/map02.txt; File size exceeded: 1049 (Max file size is 1024);",
		},
		"error/validateStageMap": {
			mapPath:  "files/test/validateStageMap/map05.txt",
			expected: "Stage Map Validation Error: files/test/validateStageMap/map05.txt; Create a boundary for the stage map with '+' (line 1,15);",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			stages := []stage{
				{mapPath: tt.mapPath},
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
		t.Error()
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
				t.Error()
			}
		})
	}
}

func TestWinOrLose(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	cases := map[string]struct {
		i            int
		life         int
		expectedI    int
		expectedLife int
		gameState    int
	}{
		"win": {
			i:            1,
			life:         2,
			expectedI:    2,
			expectedLife: 2,
			gameState:    win,
		},
		"lose": {
			i:            1,
			life:         2,
			expectedI:    1,
			expectedLife: 1,
			gameState:    lose,
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			gameState = tt.gameState
			if err := winOrLose(&tt.i, &tt.life); err != nil {
				t.Error()
			}
			if tt.i != tt.expectedI || tt.life != tt.expectedLife {
				t.Error()
			}
		})
	}
}
