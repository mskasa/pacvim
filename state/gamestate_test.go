package state

import (
	"errors"
	"testing"
	"time"

	"github.com/masahiro-kasatani/pacvim/input"
)

// --- NewGameState ---

func TestNewGameState(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	if gs.Life != 3 {
		t.Errorf("want Life=3, got %d", gs.Life)
	}
	if gs.Phase != PhaseOpening {
		t.Errorf("want PhaseOpening, got %v", gs.Phase)
	}
	if gs.StageIdx != 0 {
		t.Errorf("want StageIdx=0, got %d", gs.StageIdx)
	}
	if gs.Player == nil {
		t.Fatal("Player should not be nil")
	}
	if gs.Stage().Grid == nil {
		t.Fatal("Stage Grid should be loaded")
	}
}

func TestNewGameStatePlayerTargetScore(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	if gs.Player.TargetScore <= 0 {
		t.Errorf("TargetScore should be positive, got %d", gs.Player.TargetScore)
	}
}

func TestNewGameStateEnemiesSpawned(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	// map01.txt には2体のハンターがいる
	if len(gs.Enemies) != 2 {
		t.Errorf("want 2 enemies, got %d", len(gs.Enemies))
	}
}

// --- countApples ---

func TestCountApples(t *testing.T) {
	grid := buildGrid([]string{
		"+++++",
		"+ooo+",
		"+++++",
	})
	if n := countApples(grid); n != 3 {
		t.Errorf("want 3 apples, got %d", n)
	}
}

func TestCountApplesEmpty(t *testing.T) {
	grid := buildGrid([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	if n := countApples(grid); n != 0 {
		t.Errorf("want 0 apples, got %d", n)
	}
}

// --- ErrQuit ---

func TestUpdateQuit(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	err = gs.Update(input.CmdQuit, 0)
	if !errors.Is(err, ErrQuit) {
		t.Errorf("want ErrQuit, got %v", err)
	}
}

// --- PhaseOpening ---

func TestUpdateOpeningIgnoresCmdNone(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	_ = gs.Update(input.CmdNone, 0)
	if gs.Phase != PhaseOpening {
		t.Errorf("CmdNone should keep PhaseOpening, got %v", gs.Phase)
	}
}

func TestUpdateOpeningStartsOnCommand(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	_ = gs.Update(input.CmdMoveRight, 0)
	if gs.Phase != PhasePlaying {
		t.Errorf("command should transition to PhasePlaying, got %v", gs.Phase)
	}
}

// --- EnemyTick ---

func TestEnemyTickAdvances(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	gs.Phase = PhasePlaying
	_ = gs.Update(input.CmdNone, 0)
	if gs.EnemyTick != 1 {
		t.Errorf("want EnemyTick=1, got %d", gs.EnemyTick)
	}
}

func TestEnemyTickResetsAfterGameSpeed(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	gs.Phase = PhasePlaying
	// GameSpeed を 1 フレーム相当に設定して強制リセット
	gs.Stages[0].GameSpeed = time.Second / 60
	_ = gs.Update(input.CmdNone, 0)
	if gs.EnemyTick != 0 {
		t.Errorf("EnemyTick should reset after GameSpeed reached, got %d", gs.EnemyTick)
	}
}

// --- 敵キャプチャ → 死亡・ライフ減少 ---

func TestEnemyCaptureReducesLife(t *testing.T) {
	gs := newGameStateWithGrid([]string{
		"+++++",
		"+   +",
		"+++++",
	}, 3)
	// 敵をプレイヤーと同じ位置に配置
	gs.Enemies = []Enemy{NewHunterBuilder().Build(gs.Player.X, gs.Player.Y)}
	gs.Phase = PhasePlaying

	_ = gs.Update(input.CmdNone, 0)
	if gs.Life != 2 {
		t.Errorf("want Life=2 after capture, got %d", gs.Life)
	}
}

func TestEnemyCaptureResetsStage(t *testing.T) {
	gs := newGameStateWithGrid([]string{
		"+++++",
		"+ooo+",
		"+++++",
	}, 3)
	gs.Enemies = []Enemy{NewHunterBuilder().Build(gs.Player.X, gs.Player.Y)}
	gs.Phase = PhasePlaying
	initialScore := gs.Player.Score

	_ = gs.Update(input.CmdNone, 0)
	// ステージリセット後はスコアが 0 に戻っている
	if gs.Player.Score != initialScore {
		t.Errorf("score should reset after stage reload, got %d", gs.Player.Score)
	}
}

func TestEnemyCaptureGameOverWhenNoLife(t *testing.T) {
	gs := newGameStateWithGrid([]string{
		"+++++",
		"+   +",
		"+++++",
	}, 1)
	gs.Enemies = []Enemy{NewHunterBuilder().Build(gs.Player.X, gs.Player.Y)}
	gs.Phase = PhasePlaying

	_ = gs.Update(input.CmdNone, 0)
	if gs.Phase != PhaseGameOver {
		t.Errorf("want PhaseGameOver, got %v", gs.Phase)
	}
}

// --- 毒 → 死亡 ---

func TestPoisonKillsPlayer(t *testing.T) {
	gs := newGameStateWithGrid([]string{
		"++++++",
		"+  X +",
		"++++++",
	}, 3)
	gs.Phase = PhasePlaying
	gs.Player.X = 2
	gs.Player.Y = 1

	_ = gs.Update(input.CmdMoveRight, 0)
	// 毒(3,1)に踏み込んで死亡 → ライフ減少
	if gs.Life != 2 {
		t.Errorf("want Life=2 after poison, got %d", gs.Life)
	}
}

// --- ステージクリア ---

func TestStageClearAdvancesStage(t *testing.T) {
	gs := newGameStateWithGrid([]string{
		"+++++",
		"+  o+",
		"+++++",
	}, 3)
	gs.Player.X = 2
	gs.Player.Y = 1
	gs.Player.TargetScore = 1
	gs.Phase = PhasePlaying

	// リンゴ(3,1)へ移動してステージクリア
	_ = gs.Update(input.CmdMoveRight, 0)
	if gs.StageIdx != 1 {
		t.Errorf("want StageIdx=1 after clear, got %d", gs.StageIdx)
	}
	if gs.Phase != PhasePlaying {
		t.Errorf("want PhasePlaying after stage clear, got %v", gs.Phase)
	}
}

func TestGameClearOnLastStage(t *testing.T) {
	grid := buildGrid([]string{
		"+++++",
		"+  o+",
		"+++++",
	})
	apples := countApples(grid)
	stages := InitStages()
	lastIdx := len(stages) - 1
	// 最終ステージのグリッドをテスト用に差し替え
	stages[lastIdx].Grid = grid
	gs := &GameState{
		Stages:   stages,
		StageIdx: lastIdx,
		Player: &Player{
			X:           2,
			Y:           1,
			State:       PlayerAlive,
			TargetScore: apples,
		},
		Enemies: []Enemy{},
		Life:    3,
		Phase:   PhasePlaying,
	}

	_ = gs.Update(input.CmdMoveRight, 0)
	if gs.Phase != PhaseGameClear {
		t.Errorf("want PhaseGameClear after last stage, got %v", gs.Phase)
	}
}

// --- 敵の重複移動防止 ---

func TestEnemiesDoNotOverlapSameTarget(t *testing.T) {
	// 2体の敵が同じセルへ移動しようとする状況
	// グリッド: 横一列の通路、プレイヤーが右端
	// E1(1,1) E2(2,1) → 両者とも右へ進もうとするが、E2 の移動先(3,1)に E1 も行こうとする
	gs := newGameStateWithGrid([]string{
		"++++++",
		"+    +",
		"++++++",
	}, 3)
	gs.Player.X = 4
	gs.Player.Y = 1
	gs.Enemies = []Enemy{
		NewHunterBuilder().Build(1, 1),
		NewHunterBuilder().Build(2, 1),
	}
	gs.Phase = PhasePlaying
	gs.Stages[0].GameSpeed = time.Second / 60 // 即座に移動

	_ = gs.Update(input.CmdNone, 0)

	x0, y0 := gs.Enemies[0].Position()
	x1, y1 := gs.Enemies[1].Position()
	if x0 == x1 && y0 == y1 {
		t.Errorf("enemies should not occupy the same cell, both at (%d,%d)", x0, y0)
	}
}

func TestEnemiesDoNotMoveToOccupiedCell(t *testing.T) {
	// E1 の現在地に E2 が移動しようとする場合、E2 は移動しない
	gs := newGameStateWithGrid([]string{
		"++++++",
		"+    +",
		"++++++",
	}, 3)
	gs.Player.X = 4
	gs.Player.Y = 1
	// E1(3,1) E2(1,1): E2 は右へ進もうとするが、E1 がいるので最終的に重なってはいけない
	gs.Enemies = []Enemy{
		NewHunterBuilder().Build(3, 1),
		NewHunterBuilder().Build(1, 1),
	}
	gs.Phase = PhasePlaying
	gs.Stages[0].GameSpeed = time.Second / 60

	// 複数フレーム動かしても重複しない
	for range 10 {
		_ = gs.Update(input.CmdNone, 0)
		x0, y0 := gs.Enemies[0].Position()
		x1, y1 := gs.Enemies[1].Position()
		if x0 == x1 && y0 == y1 {
			t.Errorf("enemies overlapped at (%d,%d)", x0, y0)
			return
		}
	}
}

// --- Stage() アクセサ ---

func TestStageAccessor(t *testing.T) {
	gs, err := NewGameState(3)
	if err != nil {
		t.Fatalf("NewGameState failed: %v", err)
	}
	if gs.Stage() != &gs.Stages[0] {
		t.Error("Stage() should return pointer to current stage")
	}
}

// --- ヘルパー ---

// newGameStateWithGrid はテスト用のカスタムグリッドを持つ GameState を構築する。
// マップファイルは使わず Stage のグリッドを直接差し替える。
func newGameStateWithGrid(rows []string, life int) *GameState {
	grid := buildGrid(rows)
	apples := countApples(grid)
	player := &Player{
		X:           1,
		Y:           1,
		State:       PlayerAlive,
		TargetScore: apples,
	}
	stages := InitStages()
	// ステージ0のグリッドを差し替え（ファイル読み込みなし）
	stages[0].Grid = grid
	return &GameState{
		Stages:   stages,
		StageIdx: 0,
		Player:   player,
		Enemies:  []Enemy{},
		Life:     life,
		Phase:    PhasePlaying,
	}
}
