package state

import (
	"testing"
	"time"
)

func TestInitStagesCount(t *testing.T) {
	stages := InitStages()
	if len(stages) != 5 {
		t.Errorf("want 5 stages, got %d", len(stages))
	}
}

func TestInitStagesLevels(t *testing.T) {
	stages := InitStages()
	for i, s := range stages {
		if s.Level != i+1 {
			t.Errorf("stage[%d].Level = %d, want %d", i, s.Level, i+1)
		}
	}
}

func TestInitStagesGameSpeed(t *testing.T) {
	stages := InitStages()
	cases := []time.Duration{
		1250 * time.Millisecond,
		1000 * time.Millisecond,
		1000 * time.Millisecond,
		750 * time.Millisecond,
		750 * time.Millisecond,
	}
	for i, want := range cases {
		if stages[i].GameSpeed != want {
			t.Errorf("stage[%d].GameSpeed = %v, want %v", i, stages[i].GameSpeed, want)
		}
	}
}

func TestInitStagesGhostConfig(t *testing.T) {
	stages := InitStages()
	// ゴーストなし: level 1, 4, 5
	noGhost := []int{0, 3, 4}
	for _, i := range noGhost {
		if stages[i].GhostConfig != nil {
			t.Errorf("stage[%d] should have no ghost, but GhostConfig is set", i)
		}
	}
	// ゴーストあり: level 2, 3
	withGhost := []int{1, 2}
	for _, i := range withGhost {
		if stages[i].GhostConfig == nil {
			t.Errorf("stage[%d] should have ghost, but GhostConfig is nil", i)
		}
	}
}

func TestInitStagesHunterConfigNotNil(t *testing.T) {
	stages := InitStages()
	for i, s := range stages {
		if s.HunterConfig.Builder == nil {
			t.Errorf("stage[%d].HunterConfig.Builder should not be nil", i)
		}
	}
}

func TestInitStagesGridIsNilBeforeLoad(t *testing.T) {
	stages := InitStages()
	for i, s := range stages {
		if s.Grid != nil {
			t.Errorf("stage[%d].Grid should be nil before Load()", i)
		}
	}
}

func TestStageLoad(t *testing.T) {
	stages := InitStages()
	s := &stages[0] // map01.txt

	spawns, err := s.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// Grid が設定されていること
	if s.Grid == nil {
		t.Fatal("Grid should be set after Load()")
	}
	if s.Grid.Width == 0 || s.Grid.Height == 0 {
		t.Errorf("Grid dimensions should be non-zero, got %dx%d", s.Grid.Width, s.Grid.Height)
	}

	// SpawnPoints が返ること
	if spawns == nil {
		t.Fatal("SpawnPoints should not be nil")
	}

	// map01.txt のプレイヤー初期位置
	if spawns.Player.X != 14 || spawns.Player.Y != 7 {
		t.Errorf("player spawn = (%d,%d), want (14,7)", spawns.Player.X, spawns.Player.Y)
	}

	// map01.txt のハンター数
	if len(spawns.Hunters) != 2 {
		t.Errorf("want 2 hunters, got %d", len(spawns.Hunters))
	}
}

func TestStageLoadSetsGridWalkable(t *testing.T) {
	stages := InitStages()
	s := &stages[0]
	if _, err := s.Load(); err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// プレイヤー初期位置は CellSpace になっている
	cell := s.Grid.At(14, 7)
	if cell.Kind != CellSpace {
		t.Errorf("player spawn cell should be CellSpace, got %v", cell.Kind)
	}
}

func TestStageLoadBuildsHunterEnemies(t *testing.T) {
	stages := InitStages()
	s := &stages[0]
	spawns, err := s.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	// HunterConfig でハンターを生成できること
	for _, sp := range spawns.Hunters {
		e := s.HunterConfig.Build(sp.X, sp.Y)
		if e.Kind() != EnemyHunter {
			t.Errorf("want EnemyHunter, got %v", e.Kind())
		}
		x, y := e.Position()
		if x != sp.X || y != sp.Y {
			t.Errorf("enemy position = (%d,%d), want (%d,%d)", x, y, sp.X, sp.Y)
		}
	}
}

func TestStageLoadWithGhost(t *testing.T) {
	stages := InitStages()
	s := &stages[1] // map02.txt（ゴーストあり）
	spawns, err := s.Load()
	if err != nil {
		t.Fatalf("Load() failed: %v", err)
	}

	if len(spawns.Ghosts) == 0 {
		t.Error("map02.txt should have ghost spawn points")
	}
	if s.GhostConfig == nil {
		t.Fatal("GhostConfig should not be nil for level 2")
	}

	for _, sp := range spawns.Ghosts {
		e := s.GhostConfig.Build(sp.X, sp.Y)
		if e.Kind() != EnemyGhost {
			t.Errorf("want EnemyGhost, got %v", e.Kind())
		}
	}
}
