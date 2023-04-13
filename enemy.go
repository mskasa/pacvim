package main

import (
	"math"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type iEnemy interface {
	getPosition() (x, y int)
	think(p *player) (x, y int)
	move(x, y int)
	hasCaptured(p *player)
	eval(p *player, x, y int) float64
}

type enemy struct {
	x           int
	y           int
	char        rune
	underRune   rune
	waitingTime int
	color       termbox.Attribute
	strategy
}

type hunter struct {
	enemy
}

type ghost struct {
	enemy
}

type strategy interface {
	eval(p *player, x, y int) float64
}
type assault struct{}
type tricky struct{}

func (h *hunter) getPosition() (x, y int) {
	return h.x, h.y
}
func (g *ghost) getPosition() (x, y int) {
	return g.x, g.y
}

func (h *hunter) think(p *player) (x, y int) {
	return fuga(h, p)
}
func (g *ghost) think(p *player) (x, y int) {
	return fuga(g, p)
}
func fuga(e iEnemy, p *player) (int, int) {
	x, y := e.getPosition()
	// 移動のための評価値算出
	up := e.eval(p, x, y-1)
	down := e.eval(p, x, y+1)
	left := e.eval(p, x-1, y)
	right := e.eval(p, x+1, y)

	// 移動
	if up <= down && up <= left && up <= right {
		return x, y - 1 // 上
	} else if down <= left && down <= right {
		return x, y + 1 // 下
	} else if left <= right {
		return x - 1, y // 左
	} else {
		return x + 1, y // 右
	}
}

func (e *enemy) move(x, y int) {
	winWidth, _ := termbox.Size()
	// 移動元のセルに元の文字をセット
	termbox.SetCell(e.x, e.y, e.underRune, e.color, termbox.ColorBlack)
	// 移動先のセル情報を保持（次の移動の際に元の文字をセットする必要があるため）
	cell := termbox.CellBuffer()[(winWidth*y)+x]
	e.x = x
	e.y = y
	e.underRune = cell.Ch
	e.color = cell.Fg
	// 移動先のセルにゴーストをセット
	termbox.SetCell(x, y, e.char, termbox.ColorRed, termbox.ColorBlack)
}
func (h *hunter) move(x, y int) {
	if !isCharWall(x, y) || !isCharEnemy(x, y) {
		h.enemy.move(x, y)
	}
}
func (g *ghost) move(x, y int) {
	if !isCharBorder(x, y) || !isCharEnemy(x, y) {
		g.enemy.move(x, y)
	}
}

func (h *hunter) hasCaptured(p *player) {
	if h.x == (p.x) && h.y == p.y {
		// プレイヤーカーソルとゴーストの座標が一致した場合
		gameState = lose
	}
}

// ゴーストがプレイヤーを捕まえたかどうかの判定
func (g *ghost) hasCaptured(p *player) {
	if g.x == (p.x) && g.y == p.y {
		// プレイヤーカーソルとゴーストの座標が一致した場合
		gameState = lose
	}
}

func (h *hunter) eval(p *player, x, y int) float64 {
	if isCharWall(x, y) || isCharEnemy(x, y) {
		// 移動先が壁もしくはゴーストの場合は移動先から除外（十分に大きな値を返却）
		return 1000
	}
	return h.strategy.eval(p, x, y)
}

func (g *ghost) eval(p *player, x, y int) float64 {
	if isCharBorder(x, y) || isCharEnemy(x, y) {
		// 移動先が壁もしくはゴーストの場合は移動先から除外（十分に大きな値を返却）
		return 1000
	}
	return g.strategy.eval(p, x, y)
}

func (s *assault) eval(p *player, x, y int) float64 {
	// X軸の距離とY軸の距離それぞれの二乗の和の平方根
	return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
}

func (s *tricky) eval(p *player, x, y int) float64 {
	if random(0, 4) == 0 {
		return 0
	} else {
		return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
	}
}

func random(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}
