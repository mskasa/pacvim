package main

import (
	"errors"
	"math"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type ghost struct {
	x         int
	y         int
	underRune rune
	color     termbox.Attribute
	plotRange []float64
	strategy
}

type strategy interface {
	eval(p *player, x, y int) float64
}
type assault struct{}
type tricky struct{}

func createGhosts(level int, b *buffer) ([]*ghost, error) {
	ghosts := []*ghost{
		{
			strategy:  &assault{},
			plotRange: []float64{0.4, 0.4}, // 2nd quadrant
		},
		{
			strategy:  &assault{},
			plotRange: []float64{0.6, 0.6}, // 4th quadrant
		},
		{
			strategy:  &tricky{},
			plotRange: []float64{0.6, 0.4}, // 1st quadrant
		},
		{
			strategy:  &tricky{},
			plotRange: []float64{0.4, 0.6}, // 3rd quadrant
		},
	}
	maxGhosts := len(ghosts)
	resultList := make([]*ghost, 0, maxGhosts)
	for i := 0; i < func() int {
		n := int(math.Ceil(float64(level)/3.0)) + 1
		if n > maxGhosts {
			n = maxGhosts
		}
		return n
	}(); i++ {
		g := ghosts[i]
		if err := g.initPosition(b); err != nil {
			return nil, err
		}
		resultList = append(resultList, g)
	}
	return resultList, nil
}

func (g *ghost) initPosition(b *buffer) error {
	i := 0
	for {
		yPlotRangeUpperLimit := len(b.lines) - 1
		yPlotRangeBorder := int(float64(yPlotRangeUpperLimit) * g.plotRange[1])
		gY := tempPosition(yPlotRangeBorder, yPlotRangeUpperLimit)
		xPlotRangeUpperLimit := len(b.lines[gY].text) + b.offset
		xPlotRangeBorder := int(float64(xPlotRangeUpperLimit) * g.plotRange[0])
		gX := tempPosition(xPlotRangeBorder, xPlotRangeUpperLimit)

		if isCharTarget(gX, gY) && g.move(gX, gY) {
			break
		}

		i++
		if i == 10000 {
			return errors.New("Play with maps that have enough targets in the ghostplot range!")
		}
	}
	return nil
}
func tempPosition(min, max int) int {
	if max-min > min {
		return random(0, min)
	}
	return random(min, max)
}

// ゴーストを行動させる
func (g *ghost) action(p *player) {

	// 移動のための評価値算出
	up := g.strategy.eval(p, g.x, g.y-1)
	down := g.strategy.eval(p, g.x, g.y+1)
	left := g.strategy.eval(p, g.x-1, g.y)
	right := g.strategy.eval(p, g.x+1, g.y)

	// 移動
	if up <= down && up <= left && up <= right {
		g.move(g.x, g.y-1) // 上
	} else if down <= left && down <= right {
		g.move(g.x, g.y+1) // 下
	} else if left <= right {
		g.move(g.x-1, g.y) // 左
	} else {
		g.move(g.x+1, g.y) // 右
	}

	if g.hasCaptured(p) {
		// ゴーストがプレイヤーを捕まえた場合
		gameState = lose
	}
}

func (s *assault) eval(p *player, x, y int) float64 {
	if isCharWall(x, y) || isCharGhost(x, y) {
		// 移動先が壁もしくはゴーストの場合は移動先から除外（十分に大きな値を返却）
		return 1000
	}
	// X軸の距離とY軸の距離それぞれの二乗の和の平方根
	return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
}

func (s *tricky) eval(p *player, x, y int) float64 {
	if isCharWall(x, y) || isCharGhost(x, y) {
		// 移動先が壁もしくはゴーストの場合は移動先から除外（十分に大きな値を返却）
		return 1000
	}
	if random(0, 4) == 0 {
		return 0
	} else {
		return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
	}
}

// ゴーストを移動させる
func (g *ghost) move(x, y int) bool {
	if !isCharWall(x, y) || !isCharGhost(x, y) {
		// 移動先が壁もしくはゴーストでなければ
		winWidth, _ := termbox.Size()
		// 移動元のセルに元の文字をセット
		termbox.SetCell(g.x, g.y, g.underRune, g.color, termbox.ColorBlack)
		// 移動先のセル情報を保持（次の移動の際に元の文字をセットする必要があるため）
		cell := termbox.CellBuffer()[(winWidth*y)+x]
		g.x = x
		g.y = y
		g.underRune = cell.Ch
		g.color = cell.Fg
		// 移動先のセルにゴーストをセット
		termbox.SetCell(x, y, chGhost, termbox.ColorRed, termbox.ColorBlack)
		return true
	}
	return false
}

// ゴーストがプレイヤーを捕まえたかどうかの判定
func (g *ghost) hasCaptured(p *player) bool {
	if g.x == (p.x) && g.y == p.y {
		// プレイヤーカーソルとゴーストの座標が一致した場合
		return true
	}
	return false
}

func random(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}
