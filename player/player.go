package player

import (
	"pacvim/buffer"
	"pacvim/window"

	game "pacvim/game"

	termbox "github.com/nsf/termbox-go"
)

type Player struct {
	X int
	Y int
}

// プレイヤー初期化
func Initialize(b *buffer.Buffer) *Player {
	p := new(Player)
	p.initPosition(b)
	return p
}
func (p *Player) initPosition(b *buffer.Buffer) {
	// マップ中央座標をセット
	p.Y = b.NumOfLines()/2 - 1
	p.X = len(b.GetTextOnLine(p.Y)) / 2
	for {
		if buffer.IsSpace(p.X, p.Y) || buffer.IsTarget(p.X, p.Y) {
			// スペースかターゲットの場合は確定
			p.moveOneSquare(0, 0)
			termbox.SetCursor(p.X, p.Y)
			break
		}
		// 適当に右へ TODO
		p.X++
	}
}

// プレイヤーの制御
func (p *Player) Control(ch chan bool, b *buffer.Buffer, w *window.Window) {
	for game.IsContinuing() {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			if game.IsFirstInput_g() {
				// 「g」が入力済みの場合
				if ev.Ch == 'g' {
					// 最初の行の行頭の単語の先頭にワープ
					p.warpWord(warpBeginningFirstWordFirstLine, b)
				}
				// 入力情報の初期化
				game.InitInputNum()
				game.InitInput_g()
			} else if ev.Ch == 'g' {
				// 「g」入力済みに
				game.FirstInput_g()
			} else if s, ok := game.IsInputNum(ev.Ch); ok {
				// 数字が入力された場合
				game.SetInputNum(s)
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
					game.Quit()
				}
				// 入力数値の初期化
				game.InitInputNum()
			}
		}
		termbox.SetCursor(p.X, p.Y)
		b.DisplayScore()
		termbox.Flush()
	}
	ch <- true
}

// 移動（十字）
func (p *Player) moveCross(xDirection, yDirection int) {
	if game.GetInputNum() != 0 {
		for i := 0; i < game.GetInputNum(); i++ {
			if !p.moveOneSquare(xDirection, yDirection) {
				break
			}
		}
	} else {
		p.moveOneSquare(xDirection, yDirection)
	}
}

// 移動（１マス）
func (p *Player) moveOneSquare(xDirection, yDirection int) bool {
	x := p.X + xDirection
	y := p.Y + yDirection
	if !buffer.IsWall(x, y) {
		p.X = x
		p.Y = y
	} else {
		return false
	}
	p.checkActionResult()
	return true
}

// 移動（単語単位）
func (p *Player) moveWordByWord(fn func(*Player) bool) {
	if game.GetInputNum() != 0 {
		for i := 0; i < game.GetInputNum(); i++ {
			if !fn(p) {
				break
			}
		}
	} else {
		fn(p)
	}
}

// 次の単語の先頭に移動
func moveBeginningNextWord(p *Player) bool {
	spaceFlg := false
	for {
		if buffer.IsSpace(p.X, p.Y) {
			spaceFlg = true
		}
		if !p.moveOneSquare(1, 0) {
			return false
		}
		if spaceFlg {
			if buffer.IsTarget(p.X, p.Y) {
				return true
			}
		}
	}
}

// ワープ（行頭・行末）
func (p *Player) warpLine(fn func(*Player)) {
	fn(p)
	p.checkActionResult()
}

// ワープ（単語の先頭）
func (p *Player) warpWord(fn func(*Player, *buffer.Buffer), b *buffer.Buffer) {
	fn(p, b)
	p.checkActionResult()
}

// 移動結果の判定
func (p *Player) checkActionResult() {
	if buffer.IsGhost(p.X, p.Y) || buffer.IsPoison(p.X, p.Y) {
		game.Lose()
	} else {
		p.turnGreen()
	}
}

// b:現在の単語もしくは前の単語の先頭に移動
func moveBeginningPrevWord(p *Player) bool {
	for buffer.IsSpace(p.X-1, p.Y) {
		if !p.moveOneSquare(-1, 0) {
			break
		}
	}
	for !buffer.IsSpace(p.X-1, p.Y) {
		if !p.moveOneSquare(-1, 0) {
			return false
		}
	}
	return true
}

// e:単語の最後の文字に移動
func moveLastWord(p *Player) bool {
	for buffer.IsSpace(p.X+1, p.Y) {
		if !p.moveOneSquare(1, 0) {
			break
		}
	}
	for !buffer.IsSpace(p.X+1, p.Y) {
		if !p.moveOneSquare(1, 0) {
			return false
		}
	}
	return true
}

// 0:行頭にワープ
func warpBeginningLine(p *Player) {
	p.X = 0
	for {
		if !buffer.IsWall(p.X, p.Y) {
			p.X++
		} else {
			break
		}
	}
	for {
		if buffer.IsWall(p.X, p.Y) {
			p.X++
		} else {
			break
		}
	}
}

// $:行末にワープ
func warpEndLine(p *Player) {
	p.X, _ = termbox.Size()
	for {
		if !buffer.IsWall(p.X, p.Y) {
			p.X--
		} else {
			break
		}
	}
	for {
		if buffer.IsWall(p.X, p.Y) {
			p.X--
		} else {
			break
		}
	}
}

// ^:行頭の単語の先頭にワープ
func warpBeginningWord(p *Player, b *buffer.Buffer) {
	warpBeginningLine(p)
	width := len(b.GetTextOnLine(p.Y)) + b.Offset
	x := p.X
	for {
		if x > width {
			return
		}
		if buffer.IsTarget(x, p.Y) {
			p.X = x
			break
		}
		x++
	}
}

// gg:最初の行の行頭の単語の先頭にワープ（入力数値があれば、その行が対象）
func warpBeginningFirstWordFirstLine(p *Player, b *buffer.Buffer) {
	if game.GetInputNum() == 0 {
		p.Y = buffer.FirstTargetY
	} else if game.GetInputNum() > buffer.LastTargetY {
		p.Y = buffer.LastTargetY
	} else {
		p.Y = game.GetInputNum() - 1
	}
	warpBeginningWord(p, b)
}

// G:最後の行の行頭の単語の先頭にワープ（入力数値があれば、その行が対象）
func warpBeginningFirstWordLastLine(p *Player, b *buffer.Buffer) {
	if game.GetInputNum() == 0 || game.GetInputNum() > buffer.LastTargetY {
		p.Y = buffer.LastTargetY
	} else if game.GetInputNum() <= buffer.FirstTargetY {
		p.Y = buffer.FirstTargetY
	} else {
		p.Y = game.GetInputNum() - 1
	}
	warpBeginningWord(p, b)
}

// ターゲットの色を変更（白→緑）
func (p *Player) turnGreen() {
	winWidth, _ := termbox.Size()
	cell := termbox.CellBuffer()[(winWidth*p.Y)+p.X]
	if rune(cell.Ch) == buffer.ChTarget && cell.Fg == termbox.ColorWhite {
		termbox.SetCell(p.X, p.Y, rune(cell.Ch), termbox.ColorGreen, termbox.ColorBlack)
		game.AddScore()
	}
}
