package state

import (
	"math"
	"math/rand"
)

// EnemyKind は敵の種類を表す。
type EnemyKind int

const (
	EnemyHunter EnemyKind = iota
	EnemyGhost
)

// Enemy は敵の振る舞いを定義するインターフェース。
type Enemy interface {
	Position() (x, y int)
	SetPosition(x, y int)
	Think(p *Player, grid *Grid) (x, y int)
	Move(x, y int, grid *Grid)
	HasCaptured(p *Player) bool
	Kind() EnemyKind
}

// Strategy は敵の移動評価戦略を定義するインターフェース。
// Eval は移動候補 (x, y) の評価値を返す。値が小さいほど優先される。
type Strategy interface {
	Eval(p *Player, x, y int) float64
}

// AssaultStrategy はプレイヤーへの直線距離を最小化する戦略（ハンター型）。
type AssaultStrategy struct{}

func (s *AssaultStrategy) Eval(p *Player, x, y int) float64 {
	dx := float64(p.X - x)
	dy := float64(p.Y - y)
	return math.Sqrt(dx*dx + dy*dy)
}

// TrickyStrategy は20%の確率でランダムな評価値を返す変則戦略。
type TrickyStrategy struct{}

func (s *TrickyStrategy) Eval(p *Player, x, y int) float64 {
	if rand.Intn(5) == 0 {
		return float64(rand.Intn(30))
	}
	dx := float64(p.X - x)
	dy := float64(p.Y - y)
	return math.Sqrt(dx*dx + dy*dy)
}

// baseEnemy は敵の共通フィールドと実装を持つ。
type baseEnemy struct {
	x, y     int
	kind     EnemyKind
	strategy Strategy
}

func (e *baseEnemy) Position() (int, int) { return e.x, e.y }
func (e *baseEnemy) SetPosition(x, y int) { e.x, e.y = x, y }
func (e *baseEnemy) Kind() EnemyKind      { return e.kind }

func (e *baseEnemy) HasCaptured(p *Player) bool {
	return e.x == p.X && e.y == p.Y
}

// thinkNext は canMove で通過可否を判断しながら最良の隣接セルを返す。
// 移動優先順位: 上 > 下 > 左 > 右（同評価値のとき）。
// 全方向がブロックされている場合は現在地を返す。
func (e *baseEnemy) thinkNext(p *Player, canMove func(int, int) bool) (int, int) {
	x, y := e.x, e.y
	const blocked = 1000.0

	eval := func(nx, ny int) float64 {
		if !canMove(nx, ny) {
			return blocked
		}
		return e.strategy.Eval(p, nx, ny)
	}

	up := eval(x, y-1)
	down := eval(x, y+1)
	left := eval(x-1, y)
	right := eval(x+1, y)

	// 全方向ブロック → 現在地で待機
	if up >= blocked && down >= blocked && left >= blocked && right >= blocked {
		return x, y
	}

	if up <= down && up <= left && up <= right {
		return x, y - 1
	} else if down <= left && down <= right {
		return x, y + 1
	} else if left <= right {
		return x - 1, y
	} else {
		return x + 1, y
	}
}

// --- Hunter ---

// hunterEnemy は壁・境界を通過できない通常の敵。
type hunterEnemy struct {
	baseEnemy
}

func (e *hunterEnemy) Think(p *Player, grid *Grid) (int, int) {
	return e.thinkNext(p, func(x, y int) bool {
		k := grid.At(x, y).Kind
		return k != CellWall && k != CellBoundary
	})
}

func (e *hunterEnemy) Move(x, y int, grid *Grid) {
	k := grid.At(x, y).Kind
	if k != CellWall && k != CellBoundary {
		e.x, e.y = x, y
	}
}

// --- Ghost ---

// ghostEnemy は壁を通過できるが境界は通過できない幽霊型の敵。
type ghostEnemy struct {
	baseEnemy
}

func (e *ghostEnemy) Think(p *Player, grid *Grid) (int, int) {
	return e.thinkNext(p, func(x, y int) bool {
		return grid.At(x, y).Kind != CellBoundary
	})
}

func (e *ghostEnemy) Move(x, y int, grid *Grid) {
	if grid.At(x, y).Kind != CellBoundary {
		e.x, e.y = x, y
	}
}

// --- Builder ---

// EnemyBuilder は Enemy を段階的に構築する。
type EnemyBuilder struct {
	kind     EnemyKind
	strategy Strategy
}

// NewHunterBuilder はデフォルト設定のハンタービルダーを返す。
func NewHunterBuilder() *EnemyBuilder {
	return &EnemyBuilder{kind: EnemyHunter, strategy: &AssaultStrategy{}}
}

// NewGhostBuilder はデフォルト設定のゴーストビルダーを返す。
func NewGhostBuilder() *EnemyBuilder {
	return &EnemyBuilder{kind: EnemyGhost, strategy: &AssaultStrategy{}}
}

// Strategize は戦略を変更する。
func (b *EnemyBuilder) Strategize(s Strategy) *EnemyBuilder {
	b.strategy = s
	return b
}

// Build は指定座標に Enemy インスタンスを生成する。
func (b *EnemyBuilder) Build(x, y int) Enemy {
	base := baseEnemy{x: x, y: y, kind: b.kind, strategy: b.strategy}
	switch b.kind {
	case EnemyHunter:
		return &hunterEnemy{base}
	default:
		return &ghostEnemy{base}
	}
}

// --- EnemyConfig ---

// EnemyConfig はステージ上の敵の生成設定を保持する。
type EnemyConfig struct {
	Builder *EnemyBuilder
}

// Build は指定座標に Enemy インスタンスを生成する。
func (c *EnemyConfig) Build(x, y int) Enemy {
	return c.Builder.Build(x, y)
}
