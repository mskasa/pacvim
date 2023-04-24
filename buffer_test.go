package main

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetDigit(t *testing.T) {
	cases := map[string]struct {
		linenum  int
		expected int
	}{
		"1digit": {
			linenum:  1,
			expected: 1,
		},
		"2digit": {
			linenum:  10,
			expected: 2,
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := getDigit(tt.linenum)
			if tt.expected != result {
				t.Error("expected:", tt.expected, "result:", result)
			}
		})
	}
}

func TestMakeLineNum(t *testing.T) {
	cases := map[string]struct {
		num      int
		maxDigit int
		expected []rune
	}{
		"1digit": {
			num:      1,
			maxDigit: 2,
			expected: []rune{' ', '1', ' '},
		},
		"2digit": {
			num:      10,
			maxDigit: 2,
			expected: []rune{'1', '0', ' '},
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			result := makeLineNum(tt.num, tt.maxDigit, tt.maxDigit+1)
			if !reflect.DeepEqual(tt.expected, result) {
				t.Error("expected:", tt.expected, "result:", result)
			}
		})
	}
}

func TestSwitchScene(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping test in short mode.")
	}
	cases := map[string]struct {
		scene    string
		expected string
	}{
		"normal": {
			scene:    sceneStart,
			expected: "",
		},
		"error": {
			scene:    "foo",
			expected: "open foo: file does not exist",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			if err := switchScene(tt.scene); err != nil {
				assert.EqualErrorf(t, err, tt.expected, "Error should be: %v, got: %v", tt.expected, err)
			}
		})
	}
}
