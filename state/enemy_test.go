package state

import (
	"testing"
)

// --- Strategy ---

func TestAssaultStrategyEval(t *testing.T) {
	s := &AssaultStrategy{}
	p := &Player{X: 5, Y: 5}

	// プレイヤーに近いほど評価値が低い
	closeVal := s.Eval(p, 4, 5) // 距離1
	farVal := s.Eval(p, 0, 0)   // 距離 sqrt(50) ≈ 7.07

	if closeVal >= farVal {
		t.Errorf("close cell should have lower eval: close=%f far=%f", closeVal, farVal)
	}
}

func TestAssaultStrategyEvalAtPlayer(t *testing.T) {
	s := &AssaultStrategy{}
	p := &Player{X: 3, Y: 3}
	if val := s.Eval(p, 3, 3); val != 0 {
		t.Errorf("eval at player position should be 0, got %f", val)
	}
}

func TestTrickyStrategyReturnsValidValue(t *testing.T) {
	s := &TrickyStrategy{}
	p := &Player{X: 5, Y: 5}
	// 決定論的なテストは困難なので戻り値が非負であることを確認
	for i := 0; i < 20; i++ {
		val := s.Eval(p, 4, 5)
		if val < 0 {
			t.Errorf("TrickyStrategy.Eval should return non-negative, got %f", val)
		}
	}
}

// --- Hunter ---

func TestHunterThinkMovesTowardPlayer(t *testing.T) {
	// プレイヤーが右にいる場合、ハンターは右へ移動すべき
	grid := buildGrid([]string{
		"+++++++",
		"+     +",
		"+++++++",
	})
	hunter := NewHunterBuilder().Build(1, 1) // x=1
	p := &Player{X: 5, Y: 1}               // プレイヤーは右

	nx, ny := hunter.Think(p, grid)
	if nx != 2 || ny != 1 {
		t.Errorf("hunter should move right toward player, got (%d,%d)", nx, ny)
	}
}

func TestHunterThinkBlockedByWall(t *testing.T) {
	// 右と下が壁なら上か左へ
	grid := buildGrid([]string{
		"+  +",
		"+  +",
		"+!!+",
		"++++",
	})
	// ハンターは(1,1)、プレイヤーは(2,2)（右下）
	hunter := NewHunterBuilder().Build(1, 1)
	p := &Player{X: 2, Y: 2}

	nx, ny := hunter.Think(p, grid)
	// 右(2,1)はスペース、下(1,2)はスペースなので右か下に動く
	if (nx == 1 && ny == 1) {
		t.Errorf("hunter should move, got same position (%d,%d)", nx, ny)
	}
}

func TestHunterThinkCannotPassWall(t *testing.T) {
	// ハンターは CellWall を通過できない
	grid := buildGrid([]string{
		"++++",
		"+ !+",
		"++++",
	})
	// ハンターは(1,1)、プレイヤーは(2,1)（壁の向こう）
	hunter := NewHunterBuilder().Build(1, 1)
	p := &Player{X: 2, Y: 1}

	nx, ny := hunter.Think(p, grid)
	// (2,1) は CellWall なので選ばれない → 現在地で待機
	if nx != 1 || ny != 1 {
		t.Errorf("hunter should stay when all moves blocked, got (%d,%d)", nx, ny)
	}
}

func TestHunterMoveUpdatesPosition(t *testing.T) {
	grid := buildGrid([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	hunter := NewHunterBuilder().Build(1, 1)
	hunter.Move(2, 1, grid)
	x, y := hunter.Position()
	if x != 2 || y != 1 {
		t.Errorf("want (2,1), got (%d,%d)", x, y)
	}
}

func TestHunterMoveBlockedByBoundary(t *testing.T) {
	grid := buildGrid([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	hunter := NewHunterBuilder().Build(1, 1)
	hunter.Move(0, 1, grid) // (0,1) は CellBoundary
	x, y := hunter.Position()
	if x != 1 || y != 1 {
		t.Errorf("hunter should not move to boundary, got (%d,%d)", x, y)
	}
}

func TestHunterMoveBlockedByWall(t *testing.T) {
	grid := buildGrid([]string{
		"++++",
		"+ !+",
		"++++",
	})
	hunter := NewHunterBuilder().Build(1, 1)
	hunter.Move(2, 1, grid) // (2,1) は CellWall
	x, y := hunter.Position()
	if x != 1 || y != 1 {
		t.Errorf("hunter should not move to wall, got (%d,%d)", x, y)
	}
}

// --- Ghost ---

func TestGhostThinkCanPassWall(t *testing.T) {
	// ゴーストは CellWall を通過できる
	grid := buildGrid([]string{
		"++++",
		"+!!+",
		"+  +",
		"++++",
	})
	// ゴーストは(1,2)、プレイヤーは(1,0)（壁の上）
	ghost := NewGhostBuilder().Build(1, 2)
	p := &Player{X: 1, Y: 0}

	nx, ny := ghost.Think(p, grid)
	// 上(1,1)は CellWall だがゴーストは通過可 → 上へ移動
	if nx != 1 || ny != 1 {
		t.Errorf("ghost should move up through wall, got (%d,%d)", nx, ny)
	}
}

func TestGhostMoveCanPassWall(t *testing.T) {
	grid := buildGrid([]string{
		"++++",
		"+!!+",
		"+  +",
		"++++",
	})
	ghost := NewGhostBuilder().Build(1, 2)
	ghost.Move(1, 1, grid) // (1,1) は CellWall
	x, y := ghost.Position()
	if x != 1 || y != 1 {
		t.Errorf("ghost should move to wall cell, got (%d,%d)", x, y)
	}
}

func TestGhostMoveBlockedByBoundary(t *testing.T) {
	grid := buildGrid([]string{
		"+++++",
		"+   +",
		"+++++",
	})
	ghost := NewGhostBuilder().Build(1, 1)
	ghost.Move(0, 1, grid) // (0,1) は CellBoundary
	x, y := ghost.Position()
	if x != 1 || y != 1 {
		t.Errorf("ghost should not move to boundary, got (%d,%d)", x, y)
	}
}

// --- HasCaptured ---

func TestEnemyHasCapturedWhenSamePosition(t *testing.T) {
	hunter := NewHunterBuilder().Build(3, 3)
	p := &Player{X: 3, Y: 3}
	if !hunter.HasCaptured(p) {
		t.Error("should be captured when positions match")
	}
}

func TestEnemyNotCapturedWhenDifferentPosition(t *testing.T) {
	hunter := NewHunterBuilder().Build(3, 3)
	p := &Player{X: 4, Y: 3}
	if hunter.HasCaptured(p) {
		t.Error("should not be captured when positions differ")
	}
}

// --- Builder ---

func TestEnemyBuilderHunterKind(t *testing.T) {
	enemy := NewHunterBuilder().Build(0, 0)
	if enemy.Kind() != EnemyHunter {
		t.Errorf("want EnemyHunter, got %v", enemy.Kind())
	}
}

func TestEnemyBuilderGhostKind(t *testing.T) {
	enemy := NewGhostBuilder().Build(0, 0)
	if enemy.Kind() != EnemyGhost {
		t.Errorf("want EnemyGhost, got %v", enemy.Kind())
	}
}

func TestEnemyBuilderStrategize(t *testing.T) {
	// Strategize で TrickyStrategy に変更できる
	enemy := NewHunterBuilder().Strategize(&TrickyStrategy{}).Build(1, 1)
	if enemy.Kind() != EnemyHunter {
		t.Errorf("kind should still be EnemyHunter, got %v", enemy.Kind())
	}
}

func TestEnemyBuilderSetsInitialPosition(t *testing.T) {
	enemy := NewHunterBuilder().Build(4, 7)
	x, y := enemy.Position()
	if x != 4 || y != 7 {
		t.Errorf("want (4,7), got (%d,%d)", x, y)
	}
}

func TestEnemyBuilderCanBuildMultiple(t *testing.T) {
	// 同じビルダーで複数のインスタンスを独立して生成できる
	b := NewHunterBuilder()
	e1 := b.Build(1, 1)
	e2 := b.Build(5, 5)

	e1.SetPosition(9, 9)
	x, y := e2.Position()
	if x != 5 || y != 5 {
		t.Errorf("e2 should be independent, got (%d,%d)", x, y)
	}
}

// --- EnemyConfig ---

func TestEnemyConfigBuild(t *testing.T) {
	config := EnemyConfig{Builder: NewHunterBuilder()}
	enemy := config.Build(2, 3)
	x, y := enemy.Position()
	if x != 2 || y != 3 {
		t.Errorf("want (2,3), got (%d,%d)", x, y)
	}
	if enemy.Kind() != EnemyHunter {
		t.Errorf("want EnemyHunter, got %v", enemy.Kind())
	}
}
