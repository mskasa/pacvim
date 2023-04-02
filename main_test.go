package main

import (
	"testing"

	termbox "github.com/nsf/termbox-go"
)

func Test_switchScene(t *testing.T) {
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
		sceneCongrats,
		sceneGoodbye,
	}
	for _, s := range scenes {
		if err := switchScene(s); err != nil {
			t.Error(err)
		}
	}
}
