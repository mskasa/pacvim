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
