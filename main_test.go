package main

import (
	"reflect"
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

func Test_dirwalk(t *testing.T) {
	expected := map[int]string{
		1:  "files/stage/map01.txt",
		2:  "files/stage/map02.txt",
		3:  "files/stage/map03.txt",
		4:  "files/stage/map04.txt",
		5:  "files/stage/map05.txt",
		6:  "files/stage/map06.txt",
		7:  "files/stage/map07.txt",
		8:  "files/stage/map08.txt",
		9:  "files/stage/map09.txt",
		10: "files/stage/map10.txt"}
	result, err := dirwalk(stageDir)
	if err != nil {
		t.Error(err)
	}
	if !reflect.DeepEqual(expected, result) {
		t.Error(err)
	}
}
