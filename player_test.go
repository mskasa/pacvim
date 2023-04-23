package main

import (
	"bytes"
	"testing"

	termbox "github.com/nsf/termbox-go"
)

func TestJumpOnCurrentLine(t *testing.T) {
	cases := map[string]struct {
		initX        int
		initY        int
		toLeftEdgeX  int
		toRightEdgeX int
	}{
		"blank": {
			initX:        16,
			initY:        1,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
		"target": {
			initX:        16,
			initY:        2,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
		"obstacle": {
			initX:        16,
			initY:        3,
			toLeftEdgeX:  6,
			toRightEdgeX: 26,
		},
		"enemy": {
			initX:        16,
			initY:        4,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
	}

	if err := termbox.Init(); err != nil {
		t.Error(err)
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		t.Error(err)
	}
	defer termbox.Close()
	stage := stage{
		mapPath:       "files/test/player/jumpOnCurrentLine/map01.txt",
		hunterBuilder: newEnemyBuilder().defaultHunter(),
		ghostBuilder:  newEnemyBuilder().defaultGhost(),
	}
	f, err := static.ReadFile(stage.mapPath)
	if err != nil {
		t.Error()
	}
	b := createBuffer(bytes.NewReader(f))
	w := createWindow(b)
	if err = w.show(b); err != nil {
		t.Error()
	}
	p := new(player)
	stage.plot(b, p)

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p.x = tt.initX
			p.y = tt.initY
			p.jumpOnCurrentLine(p.toLeftEdge)
			if !(p.x == tt.toLeftEdgeX && p.y == tt.initY) {
				t.Error("expected:", tt.toLeftEdgeX, tt.initY, "result:", p.x, p.y)
			}
			p.jumpOnCurrentLine(p.toRightEdge)
			if !(p.x == tt.toRightEdgeX && p.y == tt.initY) {
				t.Error("expected:", tt.toRightEdgeX, tt.initY, "result:", p.x, p.y)
			}
		})
	}
}
