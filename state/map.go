package state

import (
	"bufio"
	"embed"
	"fmt"
)

//go:embed files/stage
var stageFS embed.FS

// CellKind はマップセルの種類を表す。
type CellKind int

const (
	CellSpace       CellKind = iota // 通路（歩行可能）
	CellWall                        // 障害物（歩行不可）
	CellBoundary                    // 境界線（歩行不可）
	CellApple                       // リンゴ（収集対象・歩行可能）
	CellAppleEaten                  // 収集済みリンゴ（歩行可能）
	CellPoison                      // 毒（歩行可能だがダメージ）
)

// Cell はグリッド上の1マスを表す。
type Cell struct {
	Kind CellKind
}

// Grid はゲームマップを二次元配列で保持する。
// Cells[y][x] の順でアクセスする。
type Grid struct {
	Cells  [][]Cell
	Width  int
	Height int
}

// At は (x, y) のセルを返す。範囲外の場合は CellWall を返す。
func (g *Grid) At(x, y int) Cell {
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return Cell{Kind: CellWall}
	}
	return g.Cells[y][x]
}

// Set は (x, y) のセルの種類を設定する。範囲外の場合は何もしない。
func (g *Grid) Set(x, y int, kind CellKind) {
	if x < 0 || y < 0 || x >= g.Width || y >= g.Height {
		return
	}
	g.Cells[y][x].Kind = kind
}

// IsWalkable はプレイヤーが (x, y) に移動できるかを返す。
func (g *Grid) IsWalkable(x, y int) bool {
	k := g.At(x, y).Kind
	return k == CellSpace || k == CellApple || k == CellAppleEaten || k == CellPoison
}

// SpawnPoint は座標を表す。
type SpawnPoint struct {
	X, Y int
}

// SpawnPoints はマップ上の各エンティティの初期位置を保持する。
type SpawnPoints struct {
	Player  SpawnPoint
	Hunters []SpawnPoint
	Ghosts  []SpawnPoint
}

// マップファイルで使用する文字。
// 実際のマップファイル（files/stage/*.txt）の文字規則に従う。
const (
	charBoundary = '+'
	charWall1    = '!'
	charWall2    = '-'
	charWall3    = '|'
	charApple    = 'o'
	charPoison   = 'X'
	charPlayer   = 'P'
	charHunter   = 'H'
	charGhost    = 'G'
)

// LoadGrid は埋め込みファイルシステムからマップファイルを読み込み、
// Grid と SpawnPoints を返す。path は "files/stage/map01.txt" の形式で指定する。
func LoadGrid(path string) (*Grid, *SpawnPoints, error) {
	f, err := stageFS.Open(path)
	if err != nil {
		return nil, nil, fmt.Errorf("open %s: %w", path, err)
	}
	defer f.Close()

	var rows [][]rune
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		rows = append(rows, []rune(scanner.Text()))
	}
	if err := scanner.Err(); err != nil {
		return nil, nil, fmt.Errorf("scan %s: %w", path, err)
	}
	if len(rows) == 0 {
		return nil, nil, fmt.Errorf("%s: empty file", path)
	}

	height := len(rows)
	width := 0
	for _, row := range rows {
		if len(row) > width {
			width = len(row)
		}
	}

	cells := make([][]Cell, height)
	for i := range cells {
		cells[i] = make([]Cell, width)
	}
	grid := &Grid{Cells: cells, Width: width, Height: height}
	spawns := &SpawnPoints{}

	for y, row := range rows {
		for x, ch := range row {
			switch ch {
			case charBoundary:
				grid.Cells[y][x] = Cell{Kind: CellBoundary}
			case charWall1, charWall2, charWall3:
				grid.Cells[y][x] = Cell{Kind: CellWall}
			case charApple:
				grid.Cells[y][x] = Cell{Kind: CellApple}
			case charPoison:
				grid.Cells[y][x] = Cell{Kind: CellPoison}
			case charPlayer:
				grid.Cells[y][x] = Cell{Kind: CellSpace}
				spawns.Player = SpawnPoint{X: x, Y: y}
			case charHunter:
				grid.Cells[y][x] = Cell{Kind: CellSpace}
				spawns.Hunters = append(spawns.Hunters, SpawnPoint{X: x, Y: y})
			case charGhost:
				grid.Cells[y][x] = Cell{Kind: CellSpace}
				spawns.Ghosts = append(spawns.Ghosts, SpawnPoint{X: x, Y: y})
			default:
				grid.Cells[y][x] = Cell{Kind: CellSpace}
			}
		}
	}

	return grid, spawns, nil
}
