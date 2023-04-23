package main

import (
	"reflect"
	"testing"

	termbox "github.com/nsf/termbox-go"
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
	if err := termbox.Init(); err != nil {
		t.Error(err)
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		t.Error(err)
	}
	defer termbox.Close()
	scenes := []string{
		sceneStart,
		sceneYouwin,
		sceneYoulose,
		sceneGoodbye,
	}
	for _, s := range scenes {
		if err := switchScene(s); err != nil {
			t.Error(err)
		}
	}
}
