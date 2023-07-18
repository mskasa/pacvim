package main

import (
	"bytes"
	"testing"

	termbox "github.com/nsf/termbox-go"
)

func TestMoveCross(t *testing.T) {
	initX, initY := 9, 4
	cases := map[string]struct {
		x         int
		y         int
		expectedX int
		expectedY int
		inputNum  int
	}{
		"left": {
			x:         -1,
			y:         0,
			expectedX: initX - 1,
			expectedY: initY,
			inputNum:  0,
		},
		"left with number": {
			x:         -1,
			y:         0,
			expectedX: initX - 3,
			expectedY: initY,
			inputNum:  3,
		},
		"right": {
			x:         1,
			y:         0,
			expectedX: initX + 1,
			expectedY: initY,
			inputNum:  0,
		},
		"right with number": {
			x:         1,
			y:         0,
			expectedX: initX + 3,
			expectedY: initY,
			inputNum:  3,
		},
		"up": {
			x:         0,
			y:         -1,
			expectedX: initX,
			expectedY: initY - 1,
			inputNum:  0,
		},
		"up with number": {
			x:         0,
			y:         -1,
			expectedX: initX,
			expectedY: initY - 2,
			inputNum:  2,
		},
		"down": {
			x:         0,
			y:         1,
			expectedX: initX,
			expectedY: initY + 1,
			inputNum:  0,
		},
		"down with number": {
			x:         0,
			y:         1,
			expectedX: initX,
			expectedY: initY + 2,
			inputNum:  2,
		},
		"right with obstacle": {
			x:         1,
			y:         0,
			expectedX: initX + 4,
			expectedY: initY,
			inputNum:  6,
		},
	}

	p, _, err := playerActionTestInit(t, "files/test/player/moveCross/map01.txt")
	if err != nil {
		t.Error()
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p.x, p.y = initX, initY
			p.inputNum = tt.inputNum
			p.moveCross(tt.x, tt.y)
			if !(p.x == tt.expectedX && p.y == tt.expectedY) {
				t.Errorf("expected %d %d but %d %d", tt.expectedX, tt.expectedY, p.x, p.y)
			}
		})
	}
}

func TestMoveByWord(t *testing.T) {
	cases := map[string]struct {
		initX         int
		initY         int
		expectedX     int
		expectedY     int
		expectedState int
		inputNum      int
		inputChar     rune
		mapPath       string
	}{
		// Check if coordinates are as expected.
		"w: from blank": {
			initX:         16,
			initY:         1,
			expectedX:     18,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'w',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"w: from word": {
			initX:         20,
			initY:         1,
			expectedX:     23,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'w',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"b: from blank": {
			initX:         16,
			initY:         1,
			expectedX:     11,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'b',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"b: from middle of word": {
			initX:         13,
			initY:         1,
			expectedX:     11,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'b',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"b: from beginning of word": {
			initX:         11,
			initY:         1,
			expectedX:     7,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'b',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"e: from blank": {
			initX:         16,
			initY:         1,
			expectedX:     21,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'e',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"e: from middle of word": {
			initX:         19,
			initY:         1,
			expectedX:     21,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'e',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"e: from end of word": {
			initX:         21,
			initY:         1,
			expectedX:     25,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      0,
			inputChar:     'e',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		// Check if the player crosses the boundary.
		"w: to the boundary": {
			initX:         16,
			initY:         1,
			expectedX:     29,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      4,
			inputChar:     'w',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"b: to the boundary": {
			initX:         16,
			initY:         1,
			expectedX:     3,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      4,
			inputChar:     'b',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"e: to the boundary": {
			initX:         16,
			initY:         1,
			expectedX:     29,
			expectedY:     1,
			expectedState: continuing,
			inputNum:      4,
			inputChar:     'e',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		// Check that poison is a string.
		"w: beyond the poison": {
			initX:         3,
			initY:         2,
			expectedX:     18,
			expectedY:     2,
			expectedState: lose,
			inputNum:      2,
			inputChar:     'w',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"b: beyond the poison": {
			initX:         16,
			initY:         2,
			expectedX:     4,
			expectedY:     2,
			expectedState: lose,
			inputNum:      0,
			inputChar:     'b',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"e: beyond the poison": {
			initX:         3,
			initY:         2,
			expectedX:     6,
			expectedY:     2,
			expectedState: lose,
			inputNum:      0,
			inputChar:     'e',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		// Check that an enemy is not a string.
		"w: beyond the enemy": {
			initX:         16,
			initY:         2,
			expectedX:     20,
			expectedY:     2,
			expectedState: lose,
			inputNum:      2,
			inputChar:     'w',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"b: beyond the enemy": {
			initX:         29,
			initY:         2,
			expectedX:     18,
			expectedY:     2,
			expectedState: lose,
			inputNum:      2,
			inputChar:     'b',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
		"e: beyond the enemy": {
			initX:         16,
			initY:         2,
			expectedX:     20,
			expectedY:     2,
			expectedState: lose,
			inputNum:      2,
			inputChar:     'e',
			mapPath:       "files/test/player/moveByWord/map01.txt",
		},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p, _, err := playerActionTestInit(t, tt.mapPath)
			if err != nil {
				t.Error(err)
			}
			p.x, p.y = tt.initX, tt.initY
			p.inputNum = tt.inputNum
			p.state = continuing
			switch tt.inputChar {
			case 'w':
				p.moveByWord(p.toBeginningOfNextWord)
			case 'b':
				p.moveByWord(p.toBeginningPrevWord)
			case 'e':
				p.moveByWord(p.toEndOfCurrentWord)
			}
			if !(p.x == tt.expectedX && p.y == tt.expectedY) {
				t.Errorf("expected %d %d but %d %d", tt.expectedX, tt.expectedY, p.x, p.y)
			}
			if p.state != tt.expectedState {
				t.Errorf("expected %d but %d", tt.expectedState, p.state)
			}
		})
	}
}

// TODO
func TestJumpOnCurrentLine(t *testing.T) {
	initX := 16
	cases := map[string]struct {
		initY        int
		toLeftEdgeX  int
		toRightEdgeX int
	}{
		"blank": {
			initY:        1,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
		"target": {
			initY:        2,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
		"obstacle": {
			initY:        3,
			toLeftEdgeX:  6,
			toRightEdgeX: 26,
		},
		"enemy": {
			initY:        4,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
		"poison": {
			initY:        5,
			toLeftEdgeX:  3,
			toRightEdgeX: 29,
		},
	}

	p, _, err := playerActionTestInit(t, "files/test/player/jumpOnCurrentLine/map01.txt")
	if err != nil {
		t.Error()
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p.x, p.y = initX, tt.initY
			p.jumpOnCurrentLine(p.toLeftEdge)
			if !(p.x == tt.toLeftEdgeX && p.y == tt.initY) {
				t.Errorf("expected %d %d but %d %d", tt.toLeftEdgeX, tt.initY, p.x, p.y)
			}
			p.jumpOnCurrentLine(p.toRightEdge)
			if !(p.x == tt.toRightEdgeX && p.y == tt.initY) {
				t.Errorf("expected %d %d but %d %d", tt.toRightEdgeX, tt.initY, p.x, p.y)
			}
		})
	}
}

func TestJumpAcrossLine(t *testing.T) {
	cases := map[string]struct {
		expectedX        int
		expectedY        int
		expectedInputG   bool
		expectedInputNum int
		inputG           bool
		inputNum         int
		inputChar        rune
		mapPath          string
	}{
		"gg with target": {
			expectedX:        29,
			expectedY:        1,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           true,
			inputNum:         0,
			inputChar:        'g',
			mapPath:          "files/test/player/jumpAcrossLine/map01.txt",
		},
		"gg no target": {
			expectedX:        3,
			expectedY:        1,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           true,
			inputNum:         0,
			inputChar:        'g',
			mapPath:          "files/test/player/jumpAcrossLine/map02.txt",
		},
		"Ng": {
			expectedX:        16,
			expectedY:        3,
			expectedInputG:   true,
			expectedInputNum: 3,
			inputG:           false,
			inputNum:         3,
			inputChar:        'g',
			mapPath:          "files/test/player/jumpAcrossLine/map01.txt",
		},
		"Ngg with target": {
			expectedX:        29,
			expectedY:        1,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           true,
			inputNum:         2,
			inputChar:        'g',
			mapPath:          "files/test/player/jumpAcrossLine/map01.txt",
		},
		"Ngg no target": {
			expectedX:        3,
			expectedY:        3,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           true,
			inputNum:         4,
			inputChar:        'g',
			mapPath:          "files/test/player/jumpAcrossLine/map02.txt",
		},
		"G with target": {
			expectedX:        29,
			expectedY:        5,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           false,
			inputNum:         0,
			inputChar:        'G',
			mapPath:          "files/test/player/jumpAcrossLine/map01.txt",
		},
		"G no target": {
			expectedX:        3,
			expectedY:        5,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           false,
			inputNum:         0,
			inputChar:        'G',
			mapPath:          "files/test/player/jumpAcrossLine/map02.txt",
		},
		"NG with target": {
			expectedX:        29,
			expectedY:        1,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           false,
			inputNum:         2,
			inputChar:        'G',
			mapPath:          "files/test/player/jumpAcrossLine/map01.txt",
		},
		"NG no target": {
			expectedX:        3,
			expectedY:        3,
			expectedInputG:   false,
			expectedInputNum: 0,
			inputG:           false,
			inputNum:         4,
			inputChar:        'G',
			mapPath:          "files/test/player/jumpAcrossLine/map02.txt",
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p, stage, err := playerActionTestInit(t, tt.mapPath)
			if err != nil {
				t.Error(err)
			}
			p.inputG = tt.inputG
			p.inputNum = tt.inputNum
			switch tt.inputChar {
			case 'g':
				p.jumpAcrossLine(p.toFirstLine, stage, tt.inputChar)
			case 'G':
				p.jumpAcrossLine(p.toLastLine, stage, tt.inputChar)
			}
			if !(p.x == tt.expectedX && p.y == tt.expectedY) {
				t.Errorf("expected %d %d but %d %d", tt.expectedX, tt.expectedY, p.x, p.y)
			}
			if !(p.inputG == tt.expectedInputG || p.inputNum != tt.expectedInputNum) {
				t.Errorf("expected %t %d but %t %d", tt.expectedInputG, tt.expectedInputNum, p.inputG, p.inputNum)
			}
		})
	}
}

func TestJudgeMoveResult(t *testing.T) {
	initX, initY := 9, 4
	cases := map[string]struct {
		x             int
		y             int
		expectedState int
	}{
		"left": {
			x:             -1,
			y:             0,
			expectedState: lose,
		},
		"right": {
			x:             1,
			y:             0,
			expectedState: lose,
		},
		"up": {
			x:             0,
			y:             -1,
			expectedState: lose,
		},
		"down": {
			x:             0,
			y:             1,
			expectedState: win,
		},
	}

	p, _, err := playerActionTestInit(t, "files/test/player/judgeMoveResult/map01.txt")
	if err != nil {
		t.Error()
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p.x, p.y = initX, initY
			p.state = continuing
			p.moveCross(tt.x, tt.y)
			if p.state != tt.expectedState {
				t.Errorf("expected %d but %d", tt.expectedState, p.state)
			}
		})
	}
}

func playerActionTestInit(t *testing.T, mapPath string) (*player, stage, error) {
	t.Helper()
	if err := termbox.Init(); err != nil {
		t.Error(err)
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		termbox.Close()
	})
	s := stage{
		mapPath:       mapPath,
		hunterBuilder: newEnemyBuilder().defaultHunter(),
		ghostBuilder:  newEnemyBuilder().defaultGhost(),
	}
	f, err := static.ReadFile(s.mapPath)
	if err != nil {
		return nil, s, err
	}
	b := createBuffer(bytes.NewReader(f))
	w := createWindow(b)
	if err = w.show(b); err != nil {
		return nil, s, err
	}
	p := new(player)
	s.plot(b, p)
	return p, s, nil
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
				t.Errorf("expected %s %t but %s %t", tt.expected.s, tt.expected.b, s, b)
			}
		})
	}
}
