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
// enemies には現在の敵一覧を渡す。walk 中の各ステップで接触判定に使用する。
// CmdNum のみ入力を蓄積してリターンし、それ以外は resetInput() で入力をリセットする。
func (p *Player) Apply(cmd input.Command, ch rune, stage *Stage, enemies []Enemy) {
	if p.State != PlayerAlive {
		return
	}

	switch cmd {
	case input.CmdNum:
		// 0 単独（inputNum == 0）は行先頭へのジャンプとして処理する。
		// 数字蓄積中（inputNum > 0）の 0 はカウントの一桁として蓄積する。
		if ch == '0' && p.inputNum == 0 {
			p.jumpTo(p.lineBeginX(stage), p.Y, stage, enemies)
			break
		}
		p.inputNum = p.inputNum*10 + int(ch-'0')
		return // カウント蓄積中はリセットしない

	case input.CmdMoveLeft:
		p.applyWalk(-1, 0, stage, enemies)
	case input.CmdMoveRight:
		p.applyWalk(1, 0, stage, enemies)
	case input.CmdMoveUp:
		p.applyWalk(0, -1, stage, enemies)
	case input.CmdMoveDown:
		p.applyWalk(0, 1, stage, enemies)

	case input.CmdWordForward:
		p.applyWordMove(p.toBeginningOfNextWord, stage, enemies)
	case input.CmdWordEnd:
		p.applyWordMove(p.toEndOfCurrentWord, stage, enemies)
	case input.CmdWordBack:
		p.applyWordMove(p.toBeginningPrevWord, stage, enemies)

	case input.CmdLineBegin:
		p.jumpTo(p.lineBeginX(stage), p.Y, stage, enemies)
	case input.CmdLineEnd:
		p.jumpTo(p.lineEndX(stage), p.Y, stage, enemies)
	case input.CmdLineFirstWord:
		p.jumpTo(p.lineFirstWordX(stage), p.Y, stage, enemies)

	case input.CmdFileTop:
		if p.inputNum != 0 {
			p.gotoLine(p.inputNum, stage, enemies)
		} else {
			p.gotoFirstLine(stage, enemies)
		}
	case input.CmdFileBottom:
		if p.inputNum != 0 {
			p.gotoLine(p.inputNum, stage, enemies)
		} else {
			p.gotoLastLine(stage, enemies)
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
// 移動先が歩行可能なら移動してセルと敵の接触を判定し true を返す。
// 移動不可なら false を返す。
func (p *Player) walkTo(targetX, targetY int, stage *Stage, enemies []Enemy) bool {
	if !stage.Grid.IsWalkable(targetX, targetY) {
		return false
	}
	p.X = targetX
	p.Y = targetY
	p.checkCell(stage, enemies)
	return true
}

// jumpTo は (targetX, targetY) へ直接ジャンプし、目的地のみ判定する。
// 経路上の壁・敵は無視する。目的地が歩行不可な場合は移動しない。
func (p *Player) jumpTo(targetX, targetY int, stage *Stage, enemies []Enemy) {
	if !stage.Grid.IsWalkable(targetX, targetY) {
		return
	}
	p.X = targetX
	p.Y = targetY
	p.checkCell(stage, enemies)
}

// checkCell は現在地のセルと敵接触を判定し、スコア加算・死亡判定を行う。
func (p *Player) checkCell(stage *Stage, enemies []Enemy) {
	// 敵との接触判定（walk の各ステップで発生する）
	for _, e := range enemies {
		ex, ey := e.Position()
		if ex == p.X && ey == p.Y {
			p.State = PlayerDead
			return
		}
	}
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
func (p *Player) applyWalk(dx, dy int, stage *Stage, enemies []Enemy) {
	count := p.repeatCount()
	for i := 0; i < count; i++ {
		if !p.walkTo(p.X+dx, p.Y+dy, stage, enemies) {
			break
		}
		if p.State != PlayerAlive {
			break
		}
	}
}

// applyWordMove は word 移動関数を repeatCount 回繰り返す。
func (p *Player) applyWordMove(fn func(*Stage, []Enemy) bool, stage *Stage, enemies []Enemy) {
	count := p.repeatCount()
	for i := 0; i < count; i++ {
		if !fn(stage, enemies) {
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
func (p *Player) toBeginningOfNextWord(stage *Stage, enemies []Enemy) bool {
	for {
		seenSpace := !isWordKind(stage.Grid.At(p.X, p.Y).Kind)
		if !p.walkTo(p.X+1, p.Y, stage, enemies) {
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
func (p *Player) toBeginningPrevWord(stage *Stage, enemies []Enemy) bool {
	// スペースを左にスキップ
	for stage.Grid.At(p.X-1, p.Y).Kind == CellSpace {
		if !p.walkTo(p.X-1, p.Y, stage, enemies) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	// ワード文字を左にスキップしてワード先頭へ
	for isWordKind(stage.Grid.At(p.X-1, p.Y).Kind) {
		if !p.walkTo(p.X-1, p.Y, stage, enemies) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	return true
}

// e: 現在のワードの末尾へ移動
func (p *Player) toEndOfCurrentWord(stage *Stage, enemies []Enemy) bool {
	// スペースを右にスキップ
	for stage.Grid.At(p.X+1, p.Y).Kind == CellSpace {
		if !p.walkTo(p.X+1, p.Y, stage, enemies) {
			return false
		}
		if p.State != PlayerAlive {
			return false
		}
	}
	// ワード文字を右にスキップしてワード末尾へ
	for isWordKind(stage.Grid.At(p.X+1, p.Y).Kind) {
		if !p.walkTo(p.X+1, p.Y, stage, enemies) {
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
func (p *Player) gotoFirstLine(stage *Stage, enemies []Enemy) {
	for y := 0; y < stage.Grid.Height; y++ {
		if p.rowHasContent(y, stage) {
			p.Y = y
			p.jumpTo(p.lineFirstWordX(stage), p.Y, stage, enemies)
			return
		}
	}
}

// gotoLastLine はコンテンツが存在する最後の行の最初のワードへジャンプする（G）。
func (p *Player) gotoLastLine(stage *Stage, enemies []Enemy) {
	for y := stage.Grid.Height - 1; y >= 0; y-- {
		if p.rowHasContent(y, stage) {
			p.Y = y
			p.jumpTo(p.lineFirstWordX(stage), p.Y, stage, enemies)
			return
		}
	}
}

// gotoLine は n 行目（1-indexed）の最初のワードへジャンプする（Ngg / NG）。
func (p *Player) gotoLine(n int, stage *Stage, enemies []Enemy) {
	y := n - 1
	if y < 0 || y >= stage.Grid.Height {
		return
	}
	if !p.rowHasContent(y, stage) {
		return
	}
	p.Y = y
	p.jumpTo(p.lineFirstWordX(stage), p.Y, stage, enemies)
}
