package main

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type player struct {
	x int
	y int
}

func (p *player) initPosition(b *buffer) {
	// マップ中央座標をセット
	p.y = len(b.lines)/2 - 1
	p.x = len(b.lines[p.y].text) / 2
	for {
		if isCharSpace(p.x, p.y) || isCharTarget(p.x, p.y) {
			// スペースかターゲットの場合は確定
			p.moveOneSquare(0, 0)
			termbox.SetCursor(p.x, p.y)
			break
		}
		// 適当に右へ TODO
		p.x++
	}
}

// プレイヤーの制御
func (p *player) action(b *buffer, w *window) error {
	isLowercaseGEntered := false
	for gameState == continuing {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if isLowercaseGEntered {
				if ev.Ch == 'g' {
					// 最初の行の行頭の単語の先頭にワープ
					p.warpWord(warpBeginningFirstWordFirstLine, b)
				}
				// 入力情報の初期化
				inputNum = 0
				isLowercaseGEntered = false
			} else if ev.Ch == 'g' {
				isLowercaseGEntered = true
			} else if s, ok := isInputNum(ev.Ch); ok {
				if inputNum != 0 {
					// 既に入力数値がある場合は文字列として数値を足算する（例：1+2=12）
					s = strconv.Itoa(inputNum) + s
				}
				inputNum, _ = strconv.Atoi(s)
			} else {
				switch ev.Ch {
				// 上に移動
				case 'k':
					p.moveCross(0, -1)
				// 下に移動
				case 'j':
					p.moveCross(0, 1)
				// 左に移動
				case 'h':
					p.moveCross(-1, 0)
				// 右に移動
				case 'l':
					p.moveCross(1, 0)
				// 次の単語の先頭に移動
				case 'w':
					p.moveWordByWord(moveBeginningNextWord)
				// 現在の単語もしくは前の単語の先頭に移動
				case 'b':
					p.moveWordByWord(moveBeginningPrevWord)
				// 単語の最後の文字に移動
				case 'e':
					p.moveWordByWord(moveLastWord)
				// 行頭にワープ
				case '0':
					p.warpLine(warpBeginningLine)
				// 行末にワープ
				case '$':
					p.warpLine(warpEndLine)
				// 行頭の単語の先頭にワープ
				case '^':
					p.warpWord(warpBeginningWord, b)
				// 最後の行の行頭の単語の先頭にワープ
				case 'G':
					p.warpWord(warpBeginningFirstWordLastLine, b)
				// ゲームをやめる
				case 'q':
					gameState = quit
				}
				// 入力数値の初期化
				inputNum = 0
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
func isInputNum(r rune) (string, bool) {
	s := string(r)
	i, err := strconv.Atoi(s)
	if err == nil && (i != 0 || (i == 0 && inputNum != 0)) {
		// 数値変換成功かつ入力数値が「0」でない場合
		return s, true
	}
	return s, false
}

// 移動（十字）
func (p *player) moveCross(xDirection, yDirection int) {
	if inputNum != 0 {
		for i := 0; i < inputNum; i++ {
			if !p.moveOneSquare(xDirection, yDirection) {
				break
			}
		}
	} else {
		p.moveOneSquare(xDirection, yDirection)
	}
}

// 移動（１マス）
func (p *player) moveOneSquare(xDirection, yDirection int) bool {
	x := p.x + xDirection
	y := p.y + yDirection
	if !isCharWall(x, y) {
		p.x = x
		p.y = y
	} else {
		return false
	}
	p.checkActionResult()
	return true
}

// 移動（単語単位）
func (p *player) moveWordByWord(fn func(*player) bool) {
	if inputNum != 0 {
		for i := 0; i < inputNum; i++ {
			if !fn(p) {
				break
			}
		}
	} else {
		fn(p)
	}
}

// 次の単語の先頭に移動
func moveBeginningNextWord(p *player) bool {
	spaceFlg := false
	for {
		if isCharSpace(p.x, p.y) {
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

// ワープ（行頭・行末）
func (p *player) warpLine(fn func(*player)) {
	fn(p)
	p.checkActionResult()
}

// ワープ（単語の先頭）
func (p *player) warpWord(fn func(*player, *buffer), b *buffer) {
	fn(p, b)
	p.checkActionResult()
}

// 移動結果の判定
func (p *player) checkActionResult() {
	if isCharGhost(p.x, p.y) || isCharPoison(p.x, p.y) {
		gameState = lose
	} else {
		p.turnGreen()
	}
}

// b:現在の単語もしくは前の単語の先頭に移動
func moveBeginningPrevWord(p *player) bool {
	for isCharSpace(p.x-1, p.y) {
		if !p.moveOneSquare(-1, 0) {
			break
		}
	}
	for !isCharSpace(p.x-1, p.y) {
		if !p.moveOneSquare(-1, 0) {
			return false
		}
	}
	return true
}

// e:単語の最後の文字に移動
func moveLastWord(p *player) bool {
	for isCharSpace(p.x+1, p.y) {
		if !p.moveOneSquare(1, 0) {
			break
		}
	}
	for !isCharSpace(p.x+1, p.y) {
		if !p.moveOneSquare(1, 0) {
			return false
		}
	}
	return true
}

// 0:行頭にワープ
func warpBeginningLine(p *player) {
	p.x = 0
	for {
		if !isCharWall(p.x, p.y) {
			p.x++
		} else {
			break
		}
	}
	for {
		if isCharWall(p.x, p.y) {
			p.x++
		} else {
			break
		}
	}
}

// $:行末にワープ
func warpEndLine(p *player) {
	p.x, _ = termbox.Size()
	for {
		if !isCharWall(p.x, p.y) {
			p.x--
		} else {
			break
		}
	}
	for {
		if isCharWall(p.x, p.y) {
			p.x--
		} else {
			break
		}
	}
}

// ^:行頭の単語の先頭にワープ
func warpBeginningWord(p *player, b *buffer) {
	warpBeginningLine(p)
	width := len(b.lines[p.y].text) + b.offset
	x := p.x
	for {
		if x > width {
			return
		}
		if isCharTarget(x, p.y) {
			p.x = x
			break
		}
		x++
	}
}

// gg:最初の行の行頭の単語の先頭にワープ（入力数値があれば、その行が対象）
func warpBeginningFirstWordFirstLine(p *player, b *buffer) {
	if inputNum == 0 {
		p.y = b.firstTargetY
	} else if inputNum > b.lastTargetY {
		p.y = b.lastTargetY
	} else {
		p.y = inputNum - 1
	}
	warpBeginningWord(p, b)
}

// G:最後の行の行頭の単語の先頭にワープ（入力数値があれば、その行が対象）
func warpBeginningFirstWordLastLine(p *player, b *buffer) {
	if inputNum == 0 || inputNum > b.lastTargetY {
		p.y = b.lastTargetY
	} else if inputNum <= b.firstTargetY {
		p.y = b.firstTargetY
	} else {
		p.y = inputNum - 1
	}
	warpBeginningWord(p, b)
}

// ターゲットの色を変更（白→緑）
func (p *player) turnGreen() {
	winWidth, _ := termbox.Size()
	cell := termbox.CellBuffer()[(winWidth*p.y)+p.x]
	if rune(cell.Ch) == chTarget && cell.Fg == termbox.ColorWhite {
		termbox.SetCell(p.x, p.y, rune(cell.Ch), termbox.ColorGreen, termbox.ColorBlack)
		score++
		if score == targetScore {
			// 目標スコアに達した場合、ステージクリア
			gameState = win
		}
	}
}
