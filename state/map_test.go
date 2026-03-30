package state

import (
	"testing"
)

func TestGridAt_OutOfBounds(t *testing.T) {
	grid := &Grid{
		Cells:  [][]Cell{{{Kind: CellApple}}},
		Width:  1,
		Height: 1,
	}
	cases := []struct {
		x, y int
	}{
		{-1, 0},
		{0, -1},
		{1, 0},
		{0, 1},
		{-1, -1},
		{100, 100},
	}
	for _, tc := range cases {
		got := grid.At(tc.x, tc.y)
		if got.Kind != CellWall {
			t.Errorf("At(%d, %d) = %v, want CellWall", tc.x, tc.y, got.Kind)
		}
	}
}

func TestGridAt_InBounds(t *testing.T) {
	grid := &Grid{
		Cells: [][]Cell{
			{{Kind: CellApple}, {Kind: CellSpace}},
			{{Kind: CellWall}, {Kind: CellBoundary}},
		},
		Width:  2,
		Height: 2,
	}
	cases := []struct {
		x, y int
		want CellKind
	}{
		{0, 0, CellApple},
		{1, 0, CellSpace},
		{0, 1, CellWall},
		{1, 1, CellBoundary},
	}
	for _, tc := range cases {
		got := grid.At(tc.x, tc.y)
		if got.Kind != tc.want {
			t.Errorf("At(%d, %d) = %v, want %v", tc.x, tc.y, got.Kind, tc.want)
		}
	}
}

func TestGridIsWalkable(t *testing.T) {
	cases := []struct {
		kind CellKind
		want bool
	}{
		{CellSpace, true},
		{CellApple, true},
		{CellAppleEaten, true},
		{CellPoison, true},
		{CellWall, false},
		{CellBoundary, false},
	}
	for _, tc := range cases {
		grid := &Grid{
			Cells:  [][]Cell{{{Kind: tc.kind}}},
			Width:  1,
			Height: 1,
		}
		got := grid.IsWalkable(0, 0)
		if got != tc.want {
			t.Errorf("IsWalkable with %v = %v, want %v", tc.kind, got, tc.want)
		}
	}
	// 範囲外は歩行不可
	grid := &Grid{Cells: [][]Cell{}, Width: 0, Height: 0}
	if grid.IsWalkable(-1, -1) {
		t.Error("IsWalkable out of bounds should return false")
	}
}

func TestLoadGrid_Map01(t *testing.T) {
	grid, spawns, err := LoadGrid("files/stage/map01.txt")
	if err != nil {
		t.Fatalf("LoadGrid failed: %v", err)
	}

	// リンゴの数を確認
	appleCount := 0
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			if grid.At(x, y).Kind == CellApple {
				appleCount++
			}
		}
	}
	const wantApples = 139
	if appleCount != wantApples {
		t.Errorf("apple count = %d, want %d", appleCount, wantApples)
	}

	// プレイヤー初期座標を確認
	wantPlayer := SpawnPoint{X: 14, Y: 7}
	if spawns.Player != wantPlayer {
		t.Errorf("player spawn = %v, want %v", spawns.Player, wantPlayer)
	}

	// プレイヤー座標のセルが CellSpace になっていることを確認
	if grid.At(spawns.Player.X, spawns.Player.Y).Kind != CellSpace {
		t.Errorf("player spawn cell should be CellSpace")
	}

	// ハンターが2体いることを確認
	if len(spawns.Hunters) != 2 {
		t.Errorf("hunter count = %d, want 2", len(spawns.Hunters))
	}

	// ハンター座標のセルが CellSpace になっていることを確認
	for _, h := range spawns.Hunters {
		if grid.At(h.X, h.Y).Kind != CellSpace {
			t.Errorf("hunter spawn (%d,%d) should be CellSpace", h.X, h.Y)
		}
	}

	// ゴーストなし（map01 にはゴーストがいない）
	if len(spawns.Ghosts) != 0 {
		t.Errorf("ghost count = %d, want 0", len(spawns.Ghosts))
	}
}
