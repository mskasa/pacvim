package state

import (
	"testing"

	"github.com/masahiro-kasatani/pacvim/input"
)

// buildGrid はテスト用の Grid を文字列スライスから生成する。
// 文字規則は state/files/stage/*.txt と同じ（+ ! - | o X）。
func buildGrid(rows []string) *Grid {
	height := len(rows)
	width := 0
	for _, r := range rows {
		if n := len([]rune(r)); n > width {
			width = n
		}
	}
	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
	}
	g := &Grid{Cells: cells, Width: width, Height: height}
	for y, row := range rows {
		for x, ch := range []rune(row) {
			switch ch {
			case '+':
				g.Cells[y][x] = Cell{Kind: CellBoundary}
			case '!', '-', '|':
				g.Cells[y][x] = Cell{Kind: CellWall}
			case 'o':
				g.Cells[y][x] = Cell{Kind: CellApple}
			case 'X':
				g.Cells[y][x] = Cell{Kind: CellPoison}
			default:
				g.Cells[y][x] = Cell{Kind: CellSpace}
			}
		}
	}
	return g
}

func newStage(rows []string) *Stage {
	return &Stage{Grid: buildGrid(rows)}
}

// --- h j k l（移動） ---

func TestPlayerMoveRight(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.X != 2 || p.Y != 1 {
		t.Errorf("want (2,1), got (%d,%d)", p.X, p.Y)
	}
}

func TestPlayerMoveLeft(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	p := &Player{X: 3, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdMoveLeft, 0, stage)
	if p.X != 2 || p.Y != 1 {
		t.Errorf("want (2,1), got (%d,%d)", p.X, p.Y)
	}
}

func TestPlayerMoveUp(t *testing.T) {
	stage := newStage([]string{
		"+ +",
		"+ +",
		"+ +",
	})
	p := &Player{X: 1, Y: 2, State: PlayerAlive}
	p.Apply(input.CmdMoveUp, 0, stage)
	if p.X != 1 || p.Y != 1 {
		t.Errorf("want (1,1), got (%d,%d)", p.X, p.Y)
	}
}

func TestPlayerMoveDown(t *testing.T) {
	stage := newStage([]string{
		"+ +",
		"+ +",
		"+ +",
	})
	p := &Player{X: 1, Y: 0, State: PlayerAlive}
	p.Apply(input.CmdMoveDown, 0, stage)
	if p.X != 1 || p.Y != 1 {
		t.Errorf("want (1,1), got (%d,%d)", p.X, p.Y)
	}
}

func TestPlayerMoveBlockedByWall(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ ! +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdMoveRight, 0, stage)
	// ! は CellWall なので移動できない
	if p.X != 1 {
		t.Errorf("should be blocked, got X=%d", p.X)
	}
}

func TestPlayerMoveWithCount(t *testing.T) {
	stage := newStage([]string{
		"++++++++",
		"+      +",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	// 3l = 3 マス右
	p.Apply(input.CmdNum, '3', stage)
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.X != 4 {
		t.Errorf("want X=4, got X=%d", p.X)
	}
}

func TestPlayerMoveCountStopsAtWall(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	// 10l = 壁手前で止まる
	p.Apply(input.CmdNum, '1', stage)
	p.Apply(input.CmdNum, '0', stage)
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.X != 3 {
		t.Errorf("want X=3 (stopped at wall), got X=%d", p.X)
	}
}

// --- リンゴ収集・勝敗 ---

func TestPlayerEatsApple(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ o +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 1}
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.Score != 1 {
		t.Errorf("want Score=1, got %d", p.Score)
	}
	if stage.Grid.At(2, 1).Kind != CellAppleEaten {
		t.Errorf("apple should be eaten")
	}
}

func TestPlayerWinsWhenAllApplesEaten(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ o +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 1}
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.State != PlayerWon {
		t.Errorf("want PlayerWon, got %v", p.State)
	}
}

func TestPlayerDiesOnPoison(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ X +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.State != PlayerDead {
		t.Errorf("want PlayerDead, got %v", p.State)
	}
}

func TestPlayerDeadIgnoresCommands(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerDead}
	p.Apply(input.CmdMoveRight, 0, stage)
	if p.X != 1 {
		t.Errorf("dead player should not move, got X=%d", p.X)
	}
}

// --- w e b（ワード移動） ---

func TestPlayerWordForward(t *testing.T) {
	// "oo  oo" → w で次のワード先頭へ
	stage := newStage([]string{
		"++++++++",
		"+oo  oo+",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdWordForward, 0, stage)
	if p.X != 5 {
		t.Errorf("want X=5 (next word), got X=%d", p.X)
	}
}

func TestPlayerWordForwardCollectsApples(t *testing.T) {
	// w でワード内を通過するときリンゴを収集する
	stage := newStage([]string{
		"+++++++++",
		"+oo   oo+",
		"+++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 10}
	p.Apply(input.CmdWordForward, 0, stage)
	// 通過したリンゴ（x=1,x=2）は収集される
	if p.Score < 1 {
		t.Errorf("want apples collected during w, got Score=%d", p.Score)
	}
}

func TestPlayerWordBack(t *testing.T) {
	// "oo  oo" → x=6 から b で前のワード先頭へ
	stage := newStage([]string{
		"++++++++",
		"+oo  oo+",
		"++++++++",
	})
	p := &Player{X: 6, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdWordBack, 0, stage)
	if p.X != 5 {
		t.Errorf("want X=5 (prev word start), got X=%d", p.X)
	}
}

func TestPlayerWordEnd(t *testing.T) {
	// "oo  oo" → x=1 から e で現在ワード末尾へ
	stage := newStage([]string{
		"++++++++",
		"+oo  oo+",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdWordEnd, 0, stage)
	if p.X != 2 {
		t.Errorf("want X=2 (word end), got X=%d", p.X)
	}
}

func TestPlayerWordForwardWithCount(t *testing.T) {
	// "oo  oo  oo" → 2w で 2番目のワード先頭へ
	stage := newStage([]string{
		"++++++++++++",
		"+oo  oo  oo+",
		"++++++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdNum, '2', stage)
	p.Apply(input.CmdWordForward, 0, stage)
	if p.X != 9 {
		t.Errorf("want X=9 (2nd next word), got X=%d", p.X)
	}
}

// --- 0 $ ^（行内ジャンプ） ---

func TestPlayerLineBegin(t *testing.T) {
	stage := newStage([]string{
		"++++++++",
		"+  ooo +",
		"++++++++",
	})
	p := &Player{X: 6, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdLineBegin, 0, stage)
	if p.X != 1 {
		t.Errorf("want X=1 (line begin), got X=%d", p.X)
	}
}

func TestPlayerLineEnd(t *testing.T) {
	stage := newStage([]string{
		"++++++++",
		"+ ooo  +",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdLineEnd, 0, stage)
	if p.X != 6 {
		t.Errorf("want X=6 (line end), got X=%d", p.X)
	}
}

func TestPlayerLineFirstWord(t *testing.T) {
	// "   ooo" → ^ で最初のリンゴへ
	stage := newStage([]string{
		"++++++++",
		"+   ooo+",
		"++++++++",
	})
	p := &Player{X: 7, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdLineFirstWord, 0, stage)
	if p.X != 4 {
		t.Errorf("want X=4 (first word), got X=%d", p.X)
	}
}

func TestPlayerLineBeginJumpsNotWalks(t *testing.T) {
	// jump コマンドは毒を飛び越える（経路判定なし）
	stage := newStage([]string{
		"+++++++",
		"+ XXXX+",
		"+++++++",
	})
	p := &Player{X: 5, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdLineBegin, 0, stage)
	// 毒を飛び越えて行先頭へ（行先頭は Space なので死なない）
	if p.X != 1 {
		t.Errorf("want X=1, got X=%d", p.X)
	}
	if p.State != PlayerAlive {
		t.Errorf("jump should skip poison, want PlayerAlive, got %v", p.State)
	}
}

// --- gg G（行間ジャンプ） ---

func TestPlayerFileTop(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ooo+",
		"+ooo+",
		"+++++",
	})
	p := &Player{X: 1, Y: 2, State: PlayerAlive}
	p.Apply(input.CmdFileTop, 0, stage)
	if p.Y != 1 {
		t.Errorf("want Y=1 (first line), got Y=%d", p.Y)
	}
	// 行の最初のワード文字へジャンプ
	if p.X != 1 {
		t.Errorf("want X=1 (first word), got X=%d", p.X)
	}
}

func TestPlayerFileBottom(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ooo+",
		"+ooo+",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdFileBottom, 0, stage)
	if p.Y != 2 {
		t.Errorf("want Y=2 (last line), got Y=%d", p.Y)
	}
}

func TestPlayerFileTopWithCount(t *testing.T) {
	// 2gg = 2行目へ（1-indexed）
	stage := newStage([]string{
		"+++++",
		"+ooo+",
		"+ooo+",
		"+ooo+",
		"+++++",
	})
	p := &Player{X: 1, Y: 3, State: PlayerAlive}
	p.Apply(input.CmdNum, '2', stage)
	p.Apply(input.CmdFileTop, 0, stage)
	if p.Y != 1 {
		t.Errorf("want Y=1 (line 2), got Y=%d", p.Y)
	}
}

func TestPlayerFileBottomWithCount(t *testing.T) {
	// 3G = 3行目へ（1-indexed）
	stage := newStage([]string{
		"+++++",
		"+ooo+",
		"+ooo+",
		"+ooo+",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdNum, '3', stage)
	p.Apply(input.CmdFileBottom, 0, stage)
	if p.Y != 2 {
		t.Errorf("want Y=2 (line 3), got Y=%d", p.Y)
	}
}
