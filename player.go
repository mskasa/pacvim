package main

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type player struct {
	x        int
	y        int
	inputNum int
	inputG   rune
}

func (p *player) action(b *buffer) error {
	for gameState == continuing {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if p.inputG == 'g' {
				if ev.Ch == 'g' {
					// Regex: *gg
					// Move cursor to the beginning of the first word on the first line
					p.jumpWord(p.jumpBeginningFirstWordFirstLine, b)
				}
				// Regex: *g.
				p.inputNum = 0
				p.inputG = 0
			} else {
				if ev.Ch == 'g' {
					p.inputG = 'g'
				} else if v, ok := p.isInputNum(ev.Ch); ok {
					p.inputNum, _ = strconv.Atoi(strconv.Itoa(p.inputNum) + v)
				} else {
					// Move cursor
					switch ev.Ch {
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
						p.moveByWord(p.moveBeginningNextWord)
					// to the beginning of the previous word
					case 'b':
						p.moveByWord(p.moveBeginningPrevWord)
					// to the end of the current word
					case 'e':
						p.moveByWord(p.moveLastWord)
					// to the beginning of the current line
					case '0':
						p.jumpLine(p.jumpBeginningLine)
					// to the end of the current line
					case '$':
						p.jumpLine(p.jumpEndLine)
					// to the beginning of the first word on the current line
					case '^':
						p.jumpWord(p.jumpBeginningWord, b)
					// to the beginning of the first word on the last line
					case 'G':
						p.jumpWord(p.jumpBeginningFirstWordLastLine, b)
					// quit
					case 'q':
						gameState = quit
					}
					p.inputNum = 0
					p.inputG = 0
				}
			}
		}
		termbox.SetCursor(p.x, p.y)
		b.plotScore()
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
		// 数値変換成功かつ入力数値が「0」でない場合
		return s, true
	}
	return s, false
}

// Move (cross)
func (p *player) moveCross(xDirection, yDirection int) {
	if p.inputNum != 0 {
		for i := 0; i < p.inputNum; i++ {
			if !p.moveOneSquare(xDirection, yDirection) {
				break
			}
		}
	} else {
		p.moveOneSquare(xDirection, yDirection)
	}
}
func (p *player) moveOneSquare(xDirection, yDirection int) bool {
	x := p.x + xDirection
	y := p.y + yDirection
	if !isCharWall(x, y) {
		p.x = x
		p.y = y
	} else {
		return false
	}
	p.judgeMoveResult()
	return true
}

// Move (by word)
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

// Jump (to beginning/end of line)
func (p *player) jumpLine(fn func()) {
	fn()
	p.judgeMoveResult()
}

// Jump (to beginning of word)
func (p *player) jumpWord(fn func(*buffer), b *buffer) {
	fn(b)
	p.judgeMoveResult()
}

// Judge the movement result
func (p *player) judgeMoveResult() {
	if isCharEnemy(p.x, p.y) || isCharPoison(p.x, p.y) {
		gameState = lose
	} else {
		// Change target color (white → green)
		winWidth, _ := termbox.Size()
		cell := termbox.CellBuffer()[(winWidth*p.y)+p.x]
		if cell.Ch == chTarget && cell.Fg == termbox.ColorWhite {
			termbox.SetCell(p.x, p.y, cell.Ch, termbox.ColorGreen, termbox.ColorBlack)
			score++
			if score == targetScore {
				gameState = win
			}
		}
	}
}

// w: Move cursor to the beginning of the next word
func (p *player) moveBeginningNextWord() bool {
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
func (p *player) moveBeginningPrevWord() bool {
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
func (p *player) moveLastWord() bool {
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
func (p *player) jumpBeginningLine() {
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
func (p *player) jumpEndLine() {
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
func (p *player) jumpBeginningWord(b *buffer) {
	p.jumpBeginningLine()
	width := len(b.lines[p.y].text) + b.offset
	x := p.x
	for {
		if x > width {
			return
		}
		if isCharTarget(x, p.y) || isCharPoison(x, p.y) {
			p.x = x
			break
		}
		x++
	}
}

// gg: Move cursor to the beginning of the first word on the first line
func (p *player) jumpBeginningFirstWordFirstLine(b *buffer) {
	if p.inputNum == 0 {
		p.y = b.firstLine
	} else {
		p.jumpToSelectedLine(b)
	}
	p.jumpBeginningWord(b)
}

// G: Move cursor to the beginning of the first word on the last line
func (p *player) jumpBeginningFirstWordLastLine(b *buffer) {
	if p.inputNum == 0 {
		p.y = b.endLine
	} else {
		p.jumpToSelectedLine(b)
	}
	p.jumpBeginningWord(b)
}

// gg or G: Move cursor to the beginning of the first word on the selected line
func (p *player) jumpToSelectedLine(b *buffer) {
	switch {
	case p.inputNum < b.firstLine:
		p.y = b.firstLine
	case p.inputNum > b.endLine:
		p.y = b.endLine
	default:
		p.y = p.inputNum - 1
	}
}
