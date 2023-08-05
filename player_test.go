package main

import (
	"bytes"
	"testing"

	termbox "github.com/nsf/termbox-go"
)

const playerTestMapPath = "files/test/player/"

func TestMoveCross(t *testing.T) {
	const initX, initY = 7, 5
	cases := map[string]struct {
		inputNum  int
		inputG    bool
		x         int
		y         int
		expectedX int
		expectedY int
	}{
		"left":                    {0, false, -1, 0, initX - 1, initY},
		"left with input number":  {2, false, -1, 0, initX - 2, initY},
		"left to the obstacle":    {6, false, -1, 0, initX - 4, initY},
		"right":                   {0, false, 1, 0, initX + 1, initY},
		"right with input number": {2, false, 1, 0, initX + 2, initY},
		"right to the obstacle":   {6, false, 1, 0, initX + 4, initY},
		"up":                      {0, false, 0, -1, initX, initY - 1},
		"up with input number":    {2, false, 0, -1, initX, initY - 2},
		"up to the obstacle":      {6, false, 0, -1, initX, initY - 2},
		"down":                    {0, false, 0, 1, initX, initY + 1},
		"down with input number":  {2, false, 0, 1, initX, initY + 2},
		"down to the obstacle":    {6, false, 0, 1, initX, initY + 2},
		// Check the input value is processed after inputNum and inputG are reset.
		"with input number and G": {2, true, -1, 0, initX - 1, initY},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p := &player{
				x:        initX,
				y:        initY,
				inputNum: tt.inputNum,
				inputG:   tt.inputG,
			}
			_, offset, err := playerActionTestInit(t, playerTestMapPath+"move_cross.txt", p)
			if err != nil {
				t.Error(err)
			}
			p.x += offset
			tt.expectedX += offset
			p.moveCross(tt.x, tt.y)
			if !(p.x == tt.expectedX && p.y == tt.expectedY) {
				t.Errorf("expected %d %d but %d %d", tt.expectedX, tt.expectedY, p.x, p.y)
			}
			if !(p.inputNum == 0 && p.inputG == false) {
				t.Errorf("expected %d %t but %d %t", 0, false, p.inputNum, p.inputG)
			}
		})
	}
}

func TestMoveByWord(t *testing.T) {
	cases := map[string]struct {
		inputNum      int
		inputChar     rune
		inputG        bool
		initX         int
		initY         int
		expectedX     int
		expectedY     int
		expectedState int
	}{
		// Check if coordinates are as expected.
		"w: from space":             {0, 'w', false, 14, 1, 16, 1, continuing},
		"w: from word":              {0, 'w', false, 18, 1, 21, 1, continuing},
		"b: from space":             {0, 'b', false, 14, 1, 9, 1, continuing},
		"b: from middle of word":    {0, 'b', false, 11, 1, 9, 1, continuing},
		"b: from beginning of word": {0, 'b', false, 9, 1, 5, 1, continuing},
		"e: from space":             {0, 'e', false, 14, 1, 19, 1, continuing},
		"e: from middle of word":    {0, 'e', false, 17, 1, 19, 1, continuing},
		"e: from end of word":       {0, 'e', false, 19, 1, 23, 1, continuing},
		// Check the player doesn't cross the obstacle.
		"w: to the obstacle": {0, 'w', false, 14, 4, 19, 4, continuing},
		"b: to the obstacle": {0, 'b', false, 14, 4, 9, 4, continuing},
		"e: to the obstacle": {0, 'e', false, 14, 4, 19, 4, continuing},
		// Check the player doesn't cross the boundary.
		"w: to the boundary": {4, 'w', false, 14, 1, 27, 1, continuing},
		"b: to the boundary": {4, 'b', false, 14, 1, 1, 1, continuing},
		"e: to the boundary": {4, 'e', false, 14, 1, 27, 1, continuing},
		// Check poison is a string.
		"w: beyond the poison": {2, 'w', false, 1, 2, 16, 2, lose},
		"b: beyond the poison": {0, 'b', false, 14, 2, 2, 2, lose},
		"e: beyond the poison": {0, 'e', false, 1, 2, 4, 2, lose},
		// Check an enemy is not a string.
		"w: beyond the enemy": {2, 'w', false, 14, 2, 18, 2, lose},
		"b: beyond the enemy": {2, 'b', false, 27, 2, 16, 2, lose},
		"e: beyond the enemy": {2, 'e', false, 14, 2, 18, 2, lose},
		// Check the input value is processed after inputNum and inputG are reset.
		"with input number and G": {4, 'w', true, 14, 1, 16, 1, continuing},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p := &player{
				x:        tt.initX,
				y:        tt.initY,
				inputNum: tt.inputNum,
				inputG:   tt.inputG,
				state:    continuing,
			}
			_, offset, err := playerActionTestInit(t, playerTestMapPath+"move_by_word.txt", p)
			if err != nil {
				t.Error(err)
			}
			p.x += offset
			tt.expectedX += offset
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
			if !(p.inputNum == 0 && p.inputG == false) {
				t.Errorf("expected %d %t but %d %t", 0, false, p.inputNum, p.inputG)
			}
		})
	}
}

func TestJumpOnCurrentLine(t *testing.T) {
	const initX = 14
	cases := map[string]struct {
		initY        int
		toLeftEdgeX  int
		toRightEdgeX int
	}{
		"to the space":    {1, 1, 27},
		"to the apple":    {2, 1, 27},
		"to the obstacle": {3, 4, 24},
		"to the enemy":    {4, 1, 27},
		"to the poison":   {5, 1, 27},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p := &player{
				x: initX,
				y: tt.initY,
			}
			_, offset, err := playerActionTestInit(t, playerTestMapPath+"jump_on_current_line.txt", p)
			if err != nil {
				t.Error(err)
			}
			p.x += offset
			tt.toLeftEdgeX += offset
			tt.toRightEdgeX += offset
			p.jumpOnCurrentLine(p.toLeftEdge)
			if !(p.x == tt.toLeftEdgeX && p.y == tt.initY) {
				t.Errorf("expected %d %d but %d %d", tt.toLeftEdgeX, tt.initY, p.x, p.y)
			}
			if !(p.inputNum == 0 && p.inputG == false) {
				t.Errorf("expected %d %t but %d %t", 0, false, p.inputNum, p.inputG)
			}
			p.jumpOnCurrentLine(p.toRightEdge)
			if !(p.x == tt.toRightEdgeX && p.y == tt.initY) {
				t.Errorf("expected %d %d but %d %d", tt.toRightEdgeX, tt.initY, p.x, p.y)
			}
		})
	}
}

func TestJumpAcrossLine(t *testing.T) {
	const initX, initY = 14, 3
	cases := map[string]struct {
		inputNum         int
		inputG           bool
		inputChar        rune
		expectedInputNum int
		expectedInputG   bool
		expectedX        int
		expectedY        int
		mapFileName      string
	}{
		"gg":  {0, true, 'g', 0, false, 27, 1, "jump_across_line.txt"},
		"Ngg": {2, true, 'g', 0, false, 27, 1, "jump_across_line.txt"},
		"G":   {0, false, 'G', 0, false, 27, 5, "jump_across_line.txt"},
		"NG":  {2, false, 'G', 0, false, 27, 1, "jump_across_line.txt"},
		// Check player behavior when there is no apples in the target row.
		"gg no apples":  {0, true, 'g', 0, false, 1, 1, "jump_across_line_no_apples.txt"},
		"Ngg no apples": {4, true, 'g', 0, false, 1, 3, "jump_across_line_no_apples.txt"},
		"G no apples":   {0, false, 'G', 0, false, 1, 5, "jump_across_line_no_apples.txt"},
		"NG no apples":  {4, false, 'G', 0, false, 1, 3, "jump_across_line_no_apples.txt"},
		// Check that inputs are saved.
		"Ng": {3, false, 'g', 3, true, 14, 3, "jump_across_line.txt"},
		// Check that inputs are reset.
		"Ng*": {3, true, 'r', 0, false, 14, 3, "jump_across_line.txt"},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p := &player{
				x:        initX,
				y:        initY,
				inputG:   tt.inputG,
				inputNum: tt.inputNum,
			}
			stage, offset, err := playerActionTestInit(t, playerTestMapPath+tt.mapFileName, p)
			if err != nil {
				t.Error(err)
			}
			p.x += offset
			tt.expectedX += offset
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
	const initX, initY = 7, 4
	cases := map[string]struct {
		x             int
		y             int
		expectedState int
	}{
		"win":                  {-1, 0, win},
		"lost to the enemy(H)": {1, 0, lose},
		"lost to the enemy(G)": {0, -1, lose},
		"lost by poison":       {0, 1, lose},
	}

	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p := &player{
				x:     initX,
				y:     initY,
				state: continuing,
			}
			_, offset, err := playerActionTestInit(t, playerTestMapPath+"judge_move_result.txt", p)
			if err != nil {
				t.Error(err)
			}
			p.x += offset
			p.moveCross(tt.x, tt.y)
			if p.state != tt.expectedState {
				t.Errorf("expected %d but %d", tt.expectedState, p.state)
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
				t.Errorf("expected %s %t but %s %t", tt.expected.s, tt.expected.b, s, b)
			}
		})
	}
}

func playerActionTestInit(t *testing.T, mapPath string, p *player) (stage, int, error) {
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
		return s, 0, err
	}
	b := createBuffer(bytes.NewReader(f))
	w := createWindow(b)
	if err = w.show(b); err != nil {
		return s, b.offset, err
	}
	s.plot(b, p)
	return s, b.offset, nil
}
