package state

import "github.com/masahiro-kasatani/pacvim/input"

// PlayerState はプレイヤーの生死状態を表す。
type PlayerState int

const (
	PlayerAlive PlayerState = iota
	PlayerDead
	PlayerWon
)

// Player はプレイヤーの状態と Vim コマンドのロジックを保持する。
type Player struct {
	X, Y        int
	Score       int
	TargetScore int
	State       PlayerState
	inputNum    int  // カウント蓄積（3w の "3" など）
	inputG      bool // g 入力待ち（gg コマンド用、input.Handler で管理するため通常は未使用）
}

// Apply はコマンドをプレイヤー状態に適用する。
// CmdNum のみ入力を蓄積してリターンし、それ以外は resetInput() で入力をリセットする。
func (p *Player) Apply(cmd input.Command, ch rune, stage *Stage) {
	if p.State != PlayerAlive {
		return
	}

	switch cmd {
	case input.CmdNum:
		p.inputNum = p.inputNum*10 + int(ch-'0')
		return // カウント蓄積中はリセットしない

	case input.CmdMoveLeft:
		p.applyWalk(-1, 0, stage)
	case input.CmdMoveRight:
		p.applyWalk(1, 0, stage)
	case input.CmdMoveUp:
		p.applyWalk(0, -1, stage)
	case input.CmdMoveDown:
		p.applyWalk(0, 1, stage)

	case input.CmdWordForward:
		p.applyWordMove(p.toBeginningOfNextWord, stage)
	case input.CmdWordEnd:
		p.applyWordMove(p.toEndOfCurrentWord, stage)
	case input.CmdWordBack:
		p.applyWordMove(p.toBeginningPrevWord, stage)

	case input.CmdLineBegin:
		p.jumpTo(p.lineBeginX(stage), p.Y, stage)
	case input.CmdLineEnd:
		p.jumpTo(p.lineEndX(stage), p.Y, stage)
	case input.CmdLineFirstWord:
		p.jumpTo(p.lineFirstWordX(stage), p.Y, stage)

	case input.CmdFileTop:
		if p.inputNum != 0 {
			p.gotoLine(p.inputNum, stage)
		} else {
			p.gotoFirstLine(stage)
		}
	case input.CmdFileBottom:
		if p.inputNum != 0 {
			p.gotoLine(p.inputNum, stage)
		} else {
			p.gotoLastLine(stage)
		}
	}

	p.resetInput()
}

// repeatCount はカウント入力の回数を返す。入力がなければ 1。
func (p *Player) repeatCount() int {
	if p.inputNum == 0 {
		return 1
	}
	return p.inputNum
}

// resetInput はカウント・g 入力待ち状態をリセットする。
func (p *Player) resetInput() {
	p.inputNum = 0
	p.inputG = false
}

// walkTo は現在地から (targetX, targetY) へ1ステップ移動する。
// 移動先が歩行可能なら移動してセルを判定し true を返す。
// 移動不可なら false を返す。
func (p *Player) walkTo(targetX, targetY int, stage *Stage) bool {
	if !stage.Grid.IsWalkable(targetX, targetY) {
		return false
	}
	p.X = targetX
	p.Y = targetY
	p.checkCell(stage)
	return true
}

// jumpTo は (targetX, targetY) へ直接ジャンプし、目的地のみ判定する。
// 経路上の壁・敵は無視する。
func (p *Player) jumpTo(targetX, targetY int, stage *Stage) {
	p.X = targetX
	p.Y = targetY
	p.checkCell(stage)
}

// checkCell は現在地のセルを判定し、スコア加算・死亡判定を行う。
func (p *Player) checkCell(stage *Stage) {
	switch stage.Grid.At(p.X, p.Y).Kind {
	case CellApple:
		p.Score++
		stage.Grid.Set(p.X, p.Y, CellAppleEaten)
		if p.Score >= p.TargetScore {
			p.State = PlayerWon
		}
	case CellPoison:
		p.State = PlayerDead
	}
}

// applyWalk は (dx, dy) 方向へ repeatCount 回 walk を繰り返す。
func (p *Player) applyWalk(dx, dy int, stage *Stage) {
	count := p.repeatCount()
	for i := 0; i < count; i++ {
		if !p.walkTo(p.X+dx, p.Y+dy, stage) {
			break
		}
		if p.State != PlayerAlive {
			break
		}
	}
}

// applyWordMove は word 移動関数を repeatCount 回繰り返す。
func (p *Player) applyWordMove(fn func(*Stage) bool, stage *Stage) {
	count := p.repeatCount()
	for i := 0; i < count; i++ {
		if !fn(stage) {
			break
		}
		if p.State != PlayerAlive {
			break
		}
	}
}

// isWordKind はセル種別がワード文字（リンゴ・収集済みリンゴ・毒）かを返す。
// 収集済みリンゴも元のワード境界を維持するためワード文字として扱う。
func isWordKind(k CellKind) bool {
	return k == CellApple || k == CellAppleEaten || k == CellPoison
}

// w: 次のワードの先頭へ移動
func (p *Player) toBeginningOfNextWord(stage *Stage) bool {
	for {
		seenSpace := !isWordKind(stage.Grid.At(p.X, p.Y).Kind)
		if !p.walkTo(p.X+1, p.Y, stage) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
		if seenSpace && isWordKind(stage.Grid.At(p.X, p.Y).Kind) {
			return true
		}
	}
}

// b: 前のワードの先頭へ移動
func (p *Player) toBeginningPrevWord(stage *Stage) bool {
	// スペースを左にスキップ
	for stage.Grid.At(p.X-1, p.Y).Kind == CellSpace {
		if !p.walkTo(p.X-1, p.Y, stage) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	// ワード文字を左にスキップしてワード先頭へ
	for isWordKind(stage.Grid.At(p.X-1, p.Y).Kind) {
		if !p.walkTo(p.X-1, p.Y, stage) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	return true
}

// e: 現在のワードの末尾へ移動
func (p *Player) toEndOfCurrentWord(stage *Stage) bool {
	// スペースを右にスキップ
	for stage.Grid.At(p.X+1, p.Y).Kind == CellSpace {
		if !p.walkTo(p.X+1, p.Y, stage) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	// ワード文字を右にスキップしてワード末尾へ
	for isWordKind(stage.Grid.At(p.X+1, p.Y).Kind) {
		if !p.walkTo(p.X+1, p.Y, stage) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	return true
}

// lineBeginX は現在行の最も左側の歩行可能セルの x 座標を返す。
func (p *Player) lineBeginX(stage *Stage) int {
	for x := 0; x < stage.Grid.Width; x++ {
		if stage.Grid.IsWalkable(x, p.Y) {
			return x
		}
	}
	return p.X
}

// lineEndX は現在行の最も右側の歩行可能セルの x 座標を返す。
func (p *Player) lineEndX(stage *Stage) int {
	for x := stage.Grid.Width - 1; x >= 0; x-- {
		if stage.Grid.IsWalkable(x, p.Y) {
			return x
		}
	}
	return p.X
}

// lineFirstWordX は現在行の最初のワード文字（リンゴ・毒）の x 座標を返す。
// ワード文字がない場合は lineBeginX にフォールバックする。
func (p *Player) lineFirstWordX(stage *Stage) int {
	for x := 0; x < stage.Grid.Width; x++ {
		if isWordKind(stage.Grid.At(x, p.Y).Kind) {
			return x
		}
	}
	return p.lineBeginX(stage)
}

// rowHasContent は行 y に歩行可能なセルが存在するかを返す。
func (p *Player) rowHasContent(y int, stage *Stage) bool {
	for x := 0; x < stage.Grid.Width; x++ {
		if stage.Grid.IsWalkable(x, y) {
			return true
		}
	}
	return false
}

// gotoFirstLine はコンテンツが存在する最初の行の最初のワードへジャンプする（gg）。
func (p *Player) gotoFirstLine(stage *Stage) {
	for y := 0; y < stage.Grid.Height; y++ {
		if p.rowHasContent(y, stage) {
			p.Y = y
			p.jumpTo(p.lineFirstWordX(stage), p.Y, stage)
			return
		}
	}
}

// gotoLastLine はコンテンツが存在する最後の行の最初のワードへジャンプする（G）。
func (p *Player) gotoLastLine(stage *Stage) {
	for y := stage.Grid.Height - 1; y >= 0; y-- {
		if p.rowHasContent(y, stage) {
			p.Y = y
			p.jumpTo(p.lineFirstWordX(stage), p.Y, stage)
			return
		}
	}
}

// gotoLine は n 行目（1-indexed）の最初のワードへジャンプする（Ngg / NG）。
func (p *Player) gotoLine(n int, stage *Stage) {
	y := n - 1
	if y < 0 || y >= stage.Grid.Height {
		return
	}
	if !p.rowHasContent(y, stage) {
		return
	}
	p.Y = y
	p.jumpTo(p.lineFirstWordX(stage), p.Y, stage)
}
