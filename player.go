package main

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type player struct {
	x           int
	y           int
	inputNum    int
	inputG      rune
	score       int
	targetScore int
}

func (p *player) action(stage stage) error {
	for gameState == continuing {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if v, ok := p.isInputNum(ev.Ch); ok {
				p.inputNum, _ = strconv.Atoi(strconv.Itoa(p.inputNum) + v)
			} else {
				switch ev.Ch {
				// quit
				case 'q':
					gameState = quit
				// Move cursor
				// to upward direction by one line
				case 'k':
					p.moveCross(0, -1)
				// to downward direction by one line
				case 'j':
					p.moveCross(0, 1)
				// to left by one position
				case 'h':
					p.moveCross(-1, 0)
				// to right by one position
				case 'l':
					p.moveCross(1, 0)
				// to the beginning of the next word
				case 'w':
					p.moveByWord(p.toBeginningOfNextWord)
				// to the beginning of the previous word
				case 'b':
					p.moveByWord(p.toBeginningPrevWord)
				// to the end of the current word
				case 'e':
					p.moveByWord(p.toEndOfCurrentWord)
				// to the beginning of the current line
				case '0':
					p.jumpOnCurrentLine(p.toLeftEdge)
				// to the end of the current line
				case '$':
					p.jumpOnCurrentLine(p.toRightEdge)
				// to the beginning of the first word on the current line
				case '^':
					p.jumpOnCurrentLine(p.toBeginningOfFirstWord)
				// to the beginning of the first word on the last line
				case 'G':
					p.jumpAcrossLine(p.toLastLine, stage)
				case 'g':
					if p.inputG == 'g' {
						// to the beginning of the first word on the first line
						p.jumpAcrossLine(p.toFirstLine, stage)
					} else {
						p.inputG = 'g'
						continue
					}
				}
				p.inputNum = 0
				p.inputG = 0
			}
		}
		termbox.SetCursor(p.x, p.y)
		p.plotScore(stage)
		if err := termbox.Flush(); err != nil {
			return err
		}
	}
	return nil
}
func (p *player) isInputNum(r rune) (string, bool) {
	s := string(r)
	i, err := strconv.Atoi(s)
	if err == nil && (i != 0 || (i == 0 && p.inputNum != 0)) {
		return s, true
	}
	return s, false
}

func (p *player) moveCross(x, y int) {
	if p.inputNum != 0 {
		for i := 0; i < p.inputNum; i++ {
			if !p.moveOneSquare(x, y) {
				break
			}
		}
	} else {
		p.moveOneSquare(x, y)
	}
}
func (p *player) moveOneSquare(x, y int) bool {
	tmpX := p.x + x
	tmpY := p.y + y
	if !isCharWall(tmpX, tmpY) {
		p.x = tmpX
		p.y = tmpY
	} else {
		return false
	}
	p.judgeMoveResult()
	return true
}

func (p *player) moveByWord(fn func() bool) {
	if p.inputNum != 0 {
		for i := 0; i < p.inputNum; i++ {
			if !fn() {
				break
			}
		}
	} else {
		fn()
	}
}

func (p *player) jumpOnCurrentLine(fn func()) {
	fn()
	p.judgeMoveResult()
}

func (p *player) jumpAcrossLine(fn func(stage), s stage) {
	fn(s)
	p.judgeMoveResult()
}

func (p *player) judgeMoveResult() {
	if isCharEnemy(p.x, p.y) || isCharPoison(p.x, p.y) {
		gameState = lose
	} else {
		// Change target color (white → green)
		winWidth, _ := termbox.Size()
		cell := termbox.CellBuffer()[(winWidth*p.y)+p.x]
		if cell.Ch == chTarget && cell.Fg == termbox.ColorWhite {
			termbox.SetCell(p.x, p.y, cell.Ch, termbox.ColorGreen, termbox.ColorBlack)
			p.score++
			if p.score == p.targetScore {
				gameState = win
			}
		}
	}
}

// w: Move cursor to the beginning of the next word
func (p *player) toBeginningOfNextWord() bool {
	spaceFlg := false
	for {
		if isCharSpace(p.x, p.y) || isCharEnemy(p.x, p.y) {
			spaceFlg = true
		}
		if !p.moveOneSquare(1, 0) {
			return false
		}
		if spaceFlg {
			if isCharTarget(p.x, p.y) {
				return true
			}
		}
	}
}

// b: Move cursor to the beginning of the previous word
func (p *player) toBeginningPrevWord() bool {
	for isCharSpace(p.x-1, p.y) || isCharEnemy(p.x-1, p.y) {
		if !p.moveOneSquare(-1, 0) {
			break
		}
	}
	for !isCharSpace(p.x-1, p.y) && !isCharEnemy(p.x-1, p.y) {
		if !p.moveOneSquare(-1, 0) {
			return false
		}
	}
	return true
}

// e: Move cursor to the end of the current word
func (p *player) toEndOfCurrentWord() bool {
	for isCharSpace(p.x+1, p.y) || isCharEnemy(p.x+1, p.y) {
		if !p.moveOneSquare(1, 0) {
			break
		}
	}
	for !isCharSpace(p.x+1, p.y) && !isCharEnemy(p.x+1, p.y) {
		if !p.moveOneSquare(1, 0) {
			return false
		}
	}
	return true
}

// 0: Move cursor to the beginning of the current line
func (p *player) toLeftEdge() {
	x := 0
	for {
		x++
		if isCharBorder(x, p.y) {
			break
		}
	}
	for {
		x++
		if !isCharWall(x, p.y) {
			break
		}
	}
	p.x = x
}

// $: Move cursor to the end of the current line
func (p *player) toRightEdge() {
	x, _ := termbox.Size()
	for {
		x--
		if isCharBorder(x, p.y) {
			break
		}
	}
	for {
		x--
		if !isCharWall(x, p.y) {
			break
		}
	}
	p.x = x
}

// ^: Move cursor to the beginning of the first word on the current line
func (p *player) toBeginningOfFirstWord() {
	p.toLeftEdge()
	x := p.x
	for {
		if isCharTarget(x, p.y) || isCharPoison(x, p.y) {
			p.x = x
			break
		}
		x++
	}
}

// gg: Move cursor to the beginning of the first word on the first line
func (p *player) toFirstLine(s stage) {
	if p.inputNum == 0 {
		for y := 1; y < s.height; y++ {
			if canMove(s, y) {
				p.y = y
				p.toBeginningOfFirstWord()
				break
			}
		}
	} else {
		y := p.inputNum - 1
		if canMove(s, y) {
			p.y = y
			p.toBeginningOfFirstWord()
		}
	}
}

// G: Move cursor to the beginning of the first word on the last line
func (p *player) toLastLine(s stage) {
	if p.inputNum == 0 {
		for y := s.height - 1; y > 0; y-- {
			if canMove(s, y) {
				p.y = y
				p.toBeginningOfFirstWord()
				break
			}
		}
	} else {
		y := p.inputNum - 1
		if canMove(s, y) {
			p.y = y
			p.toBeginningOfFirstWord()
		}
	}
}
func canMove(s stage, y int) bool {
	x := getOffset(s.height)
	for x < s.width {
		if !isCharWall(x, y) && !isCharEnemy(x, y) {
			return true
		}
		x++
	}
	return false
}

func (p *player) plotScore(s stage) {
	position := s.height
	text := []rune("score: " + strconv.Itoa(p.score) + "/" + strconv.Itoa(p.targetScore))
	for x, r := range text {
		termbox.SetCell(x, position, r, termbox.ColorGreen, termbox.ColorBlack)
	}
}
