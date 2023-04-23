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

func TestIsInputNum(t *testing.T) {
	type expectedValues struct {
		s string
		b bool
	}
	cases := map[string]struct {
		player   player
		arg      rune
		expected expectedValues
	}{
		"Argument cannot be converted to a number.": {
			player: player{
				inputNum: 0,
			},
			arg: 'k',
			expected: expectedValues{
				s: "k",
				b: false,
			},
		},
		"Argument can be converted to a number.": {
			player: player{
				inputNum: 0,
			},
			arg: '2',
			expected: expectedValues{
				s: "2",
				b: true,
			},
		},
		"0 is the number 0.": {
			player: player{
				inputNum: 1,
			},
			arg: '0',
			expected: expectedValues{
				s: "0",
				b: true,
			},
		},
		"0 is the special string 0": {
			player: player{
				inputNum: 0,
			},
			arg: '0',
			expected: expectedValues{
				s: "0",
				b: false,
			},
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			t.Parallel()
			s, b := tt.player.isInputNum(tt.arg)
			if s != tt.expected.s || b != tt.expected.b {
				t.Error("expected:", tt.expected.s, tt.expected.b, "result:", s, b)
			}
		})
	}
}
