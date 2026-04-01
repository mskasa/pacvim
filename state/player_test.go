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
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdMoveLeft, 0, stage, nil)
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
	p.Apply(input.CmdMoveUp, 0, stage, nil)
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
	p.Apply(input.CmdMoveDown, 0, stage, nil)
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
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdNum, '3', stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdNum, '1', stage, nil)
	p.Apply(input.CmdNum, '0', stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdMoveRight, 0, stage, nil)
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
	p.Apply(input.CmdWordForward, 0, stage, nil)
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
	p.Apply(input.CmdWordForward, 0, stage, nil)
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
	p.Apply(input.CmdWordBack, 0, stage, nil)
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
	p.Apply(input.CmdWordEnd, 0, stage, nil)
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
	p.Apply(input.CmdNum, '2', stage, nil)
	p.Apply(input.CmdWordForward, 0, stage, nil)
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
	p.Apply(input.CmdLineBegin, 0, stage, nil)
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
	p.Apply(input.CmdLineEnd, 0, stage, nil)
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
	p.Apply(input.CmdLineFirstWord, 0, stage, nil)
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
	p.Apply(input.CmdLineBegin, 0, stage, nil)
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
	p.Apply(input.CmdFileTop, 0, stage, nil)
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
	p.Apply(input.CmdFileBottom, 0, stage, nil)
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
	p.Apply(input.CmdNum, '2', stage, nil)
	p.Apply(input.CmdFileTop, 0, stage, nil)
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
	p.Apply(input.CmdNum, '3', stage, nil)
	p.Apply(input.CmdFileBottom, 0, stage, nil)
	if p.Y != 2 {
		t.Errorf("want Y=2 (line 3), got Y=%d", p.Y)
	}
}

// --- walk vs jump の動作差異 ---

// TestWalkCollectsApplesAlongPath: walk コマンド（l）は経路上のリンゴを全て収集する。
func TestWalkCollectsApplesAlongPath(t *testing.T) {
	stage := newStage([]string{
		"++++++",
		"+ ooo+",
		"++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdNum, '4', stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
	// x=2,3,4 の 3 個を通過
	if p.Score != 3 {
		t.Errorf("want Score=3 (walk collects all apples), got %d", p.Score)
	}
}

// TestJumpSkipsEnemiesAlongPath: jump コマンド（0）は経路上の敵を無視して目的地のみ判定する。
func TestJumpSkipsEnemiesAlongPath(t *testing.T) {
	stage := newStage([]string{
		"++++++",
		"+    +",
		"++++++",
	})
	// x=2 に敵を配置。プレイヤーは x=4 から 0 で x=1 へジャンプ。
	enemy := NewHunterBuilder().Build(2, 1)
	p := &Player{X: 4, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdLineBegin, 0, stage, []Enemy{enemy})
	if p.X != 1 {
		t.Errorf("want X=1 (jumped to line begin), got X=%d", p.X)
	}
	if p.State != PlayerAlive {
		t.Errorf("jump should skip enemy along path, want PlayerAlive, got %v", p.State)
	}
}

// TestWalkDiesOnPoison: walk コマンド（l）は経路上の毒を踏むと即死する。
func TestWalkDiesOnPoison(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ X +",
		"+++++",
	})
	// x=1 から 3l で毒(x=2)を踏む
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdNum, '3', stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
	if p.State != PlayerDead {
		t.Errorf("walk into poison should kill player, got %v", p.State)
	}
	// 毒で止まっているので x=3 には到達しない
	if p.X != 2 {
		t.Errorf("want X=2 (stopped at poison), got X=%d", p.X)
	}
}

// --- 数値入力（カウントプレフィックス） ---

// TestCountMoveRight: 3l で右に3マス移動する。
func TestCountMoveRight(t *testing.T) {
	stage := newStage([]string{
		"++++++++",
		"+      +",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdNum, '3', stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
	if p.X != 4 {
		t.Errorf("3l: want X=4, got X=%d", p.X)
	}
}

// TestCount10MoveDown: 1 → 0 → j で10マス下移動する。
func TestCount10MoveDown(t *testing.T) {
	stage := newStage([]string{
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
		"+ +",
	})
	p := &Player{X: 1, Y: 0, State: PlayerAlive}
	p.Apply(input.CmdNum, '1', stage, nil)
	p.Apply(input.CmdNum, '0', stage, nil)
	p.Apply(input.CmdMoveDown, 0, stage, nil)
	if p.Y != 10 {
		t.Errorf("10j: want Y=10, got Y=%d", p.Y)
	}
}

// TestCountGotoLine: 5G で5行目（1-indexed）に移動する。
func TestCountGotoLine(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+ooo+",
		"+ooo+",
		"+ooo+",
		"+ooo+",
		"+ooo+",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdNum, '5', stage, nil)
	p.Apply(input.CmdFileBottom, 0, stage, nil)
	// gotoLine(5) は 1-indexed → y=4（0-indexed）
	if p.Y != 4 {
		t.Errorf("5G: want Y=4 (line 5, 1-indexed), got Y=%d", p.Y)
	}
}

// TestCountNotResetByCmdNone: CmdNone が来てもカウントはリセットされない。
// 毎フレーム Apply が呼ばれる実ゲームで 3 と l の間に CmdNone フレームが挟まっても正しく動作する。
func TestCountNotResetByCmdNone(t *testing.T) {
	stage := newStage([]string{
		"++++++++",
		"+      +",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdNum, '3', stage, nil)
	// CmdNone が複数フレーム挟まってもカウントは保持される
	p.Apply(input.CmdNone, 0, stage, nil)
	p.Apply(input.CmdNone, 0, stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
	if p.X != 4 {
		t.Errorf("3 + (CmdNone×2) + l: want X=4, got X=%d", p.X)
	}
}

// --- f F t T ; ,（文字検索） ---

// TestFindForward: fo で右方向の次の o（リンゴ）まで walk 移動し、リンゴを取得する。
// "+  o o  +" → x=3 に最初の o がある
func TestFindForward(t *testing.T) {
	stage := newStage([]string{
		"++++++++++",
		"+  o o   +",
		"++++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	// 最初の o は x=3
	if p.X != 3 {
		t.Errorf("fo: want X=3, got X=%d", p.X)
	}
	// x=3 の o を取得
	if p.Score != 1 {
		t.Errorf("fo: want Score=1, got %d", p.Score)
	}
}

// TestFindForwardCollectsApplesAlongPath: fo で経路上のリンゴも収集する。
func TestFindForwardCollectsApplesAlongPath(t *testing.T) {
	stage := newStage([]string{
		"+++++++",
		"+ oo o+",
		"+++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	// fo: 最初の o は x=2。経路上のリンゴ（x=2）を取得
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	if p.X != 2 {
		t.Errorf("fo: want X=2, got X=%d", p.X)
	}
	// 2回目の fo: 次の o は x=3。x=3 のリンゴを取得
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	if p.X != 3 {
		t.Errorf("fo: want X=3, got X=%d", p.X)
	}
	if p.Score != 2 {
		t.Errorf("fo: want Score=2, got %d", p.Score)
	}
}

// TestFindBack: Fo で左方向の次の o まで移動する。
// "+  o o  +" → x=5 に右側の o がある
func TestFindBack(t *testing.T) {
	stage := newStage([]string{
		"++++++++++",
		"+  o o   +",
		"++++++++++",
	})
	p := &Player{X: 8, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdFindBack, 'o', stage, nil)
	// 右から見て最初の o は x=5
	if p.X != 5 {
		t.Errorf("Fo: want X=5, got X=%d", p.X)
	}
}

// TestTillForward: to で右方向の次の o の1つ手前まで移動する。
func TestTillForward(t *testing.T) {
	stage := newStage([]string{
		"+++++++",
		"+   o +",
		"+++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdTillForward, 'o', stage, nil)
	// o は x=4、1つ手前は x=3
	if p.X != 3 {
		t.Errorf("to: want X=3, got X=%d", p.X)
	}
}

// TestTillBack: To で左方向の次の o の1つ右まで移動する。
func TestTillBack(t *testing.T) {
	stage := newStage([]string{
		"+++++++",
		"+ o   +",
		"+++++++",
	})
	p := &Player{X: 5, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdTillBack, 'o', stage, nil)
	// o は x=2、1つ右は x=3
	if p.X != 3 {
		t.Errorf("To: want X=3, got X=%d", p.X)
	}
}

// TestFindNotFound: 対象文字が存在しない場合は移動しない。
func TestFindNotFound(t *testing.T) {
	stage := newStage([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	if p.X != 1 {
		t.Errorf("fo (not found): should not move, got X=%d", p.X)
	}
}

// TestRepeatFind: fo の後に ; で同じ方向に繰り返し移動する。
// "+  o   o   o  +" → x=3, x=7, x=11 に o がある
func TestRepeatFind(t *testing.T) {
	stage := newStage([]string{
		"+++++++++++++++",
		"+  o   o   o  +",
		"+++++++++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	// fo: x=3 へ
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	if p.X != 3 {
		t.Fatalf("fo: want X=3, got X=%d", p.X)
	}
	// ; で繰り返し: x=7 へ
	p.Apply(input.CmdRepeatFind, 0, stage, nil)
	if p.X != 7 {
		t.Errorf(";: want X=7, got X=%d", p.X)
	}
	// ; で繰り返し: x=11 へ
	p.Apply(input.CmdRepeatFind, 0, stage, nil)
	if p.X != 11 {
		t.Errorf(";: want X=11, got X=%d", p.X)
	}
}

// TestRepeatFindReverse: fo の後に , で逆方向に移動する。
// プレイヤーを中間点（x=5）に置き、右の o（x=8）へ fo で移動後、
// 左の o（x=2）へ , で戻る。fo で x=2 は通らないので食べられていない。
func TestRepeatFindReverse(t *testing.T) {
	stage := newStage([]string{
		"+++++++++++",
		"+ o     o +",
		"+++++++++++",
	})
	// x=2: o, x=8: o。プレイヤーは x=5 からスタート
	p := &Player{X: 5, Y: 1, State: PlayerAlive, TargetScore: 99}
	// fo: x=8 へ（x=2 は通らないので食べられない）
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	if p.X != 8 {
		t.Fatalf("fo: want X=8, got X=%d", p.X)
	}
	// lastFind = {forward: true, ch: 'o'}
	// , で逆方向（左）: x=7 から検索 → x=2 の o へ
	p.Apply(input.CmdRepeatFindRev, 0, stage, nil)
	if p.X != 2 {
		t.Errorf(",: want X=2, got X=%d", p.X)
	}
}

// TestRepeatFindNilLastFind: lastFind が nil の場合は何もしない。
func TestRepeatFindNilLastFind(t *testing.T) {
	stage := newStage([]string{
		"+++++++",
		"+ o o +",
		"+++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	p.Apply(input.CmdRepeatFind, 0, stage, nil)
	if p.X != 1 {
		t.Errorf("; with no lastFind: should not move, got X=%d", p.X)
	}
}

// TestFindWalkDiesOnPoison: fo は walk なので経路上に毒があれば死亡する。
func TestFindWalkDiesOnPoison(t *testing.T) {
	stage := newStage([]string{
		"+++++++",
		"+ X o +",
		"+++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive, TargetScore: 99}
	p.Apply(input.CmdFindForward, 'o', stage, nil)
	// 経路上の x=2 に毒があるので死亡
	if p.State != PlayerDead {
		t.Errorf("fo through poison: want PlayerDead, got %v", p.State)
	}
	if p.X != 2 {
		t.Errorf("fo through poison: want stopped at X=2, got X=%d", p.X)
	}
}

// TestCountResetAfterCommand: コマンド実行後にカウントがリセットされ、次のコマンドに引き継がれない。
func TestCountResetAfterCommand(t *testing.T) {
	stage := newStage([]string{
		"++++++++",
		"+      +",
		"++++++++",
	})
	p := &Player{X: 1, Y: 1, State: PlayerAlive}
	// 3l でX=4へ
	p.Apply(input.CmdNum, '3', stage, nil)
	p.Apply(input.CmdMoveRight, 0, stage, nil)
	if p.X != 4 {
		t.Errorf("3l: want X=4, got X=%d", p.X)
	}
	// カウントリセット後の l は1マスだけ移動する
	p.Apply(input.CmdMoveRight, 0, stage, nil)
	if p.X != 5 {
		t.Errorf("after reset, l: want X=5, got X=%d", p.X)
	}
}
