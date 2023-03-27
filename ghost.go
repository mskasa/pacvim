package main

import (
	"errors"
	"math"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type ghost struct {
	x         int
	y         int
	underRune rune
	color     termbox.Attribute
	strategy  strategy
	plotRange []float64
}

type strategy interface {
	eval(p *player, x, y int) float64
}
type strategyA struct{}
type strategyB struct{}

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
func numOfGhosts() int {
	numOfGhost := int(math.Ceil(float64(level)/3.0)) + 1
	if numOfGhost > maxNumOfGhosts {
		numOfGhost = maxNumOfGhosts
	}
	return numOfGhost
}

// ゴーストの制御
func control(ch chan bool, p *player, gList []*ghost) {
	var wg sync.WaitGroup

	for gameState == continuing {
		// ゲームが継続している限り
		wg.Add(len(gList))
		// ゴーストを並行に行動させる
		for _, g := range gList {
			go g.action(&wg, p)
		}
		// ゴーストの行動の同期
		wg.Wait()
		// ゴーストを表示
		termbox.Flush()
		// ゴーストの行動間隔
		time.Sleep(gameSpeed)
	}
	ch <- true
}

// ゴーストを行動させる
func (g *ghost) action(wg *sync.WaitGroup, p *player) {
	defer wg.Done()

	// 移動のための評価値算出
	up := g.strategy.eval(p, g.x, g.y-1)
	down := g.strategy.eval(p, g.x, g.y+1)
	left := g.strategy.eval(p, g.x-1, g.y)
	right := g.strategy.eval(p, g.x+1, g.y)

	// 移動
	if up <= down && up <= left && up <= right {
		g.move(g.x, g.y-1) // 上
	} else if down <= left && down <= right && down <= up {
		g.move(g.x, g.y+1) // 下
	} else if left <= right && left <= up && left <= down {
		g.move(g.x-1, g.y) // 左
	} else if right <= up && right <= down && right <= left {
		g.move(g.x+1, g.y) // 右
	}

	if g.hasCaptured(p) {
		// ゴーストがプレイヤーを捕まえた場合
		gameState = lose
	}
}

func newStrategy(i int) strategy {
	if i/2+1 == 1 {
		return &strategyA{}
	}
	return &strategyB{}
}

// プレイヤー追跡タイプ
func (s *strategyA) eval(p *player, x, y int) float64 {
	if isCharWall(x, y) || isCharGhost(x, y) {
		// 移動先が壁もしくはゴーストの場合は移動先から除外（十分に大きな値を返却）
		return 1000
	}
	// X軸の距離とY軸の距離それぞれの二乗の和の平方根
	return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
}

// 待ち伏せ徘徊タイプ
func (s *strategyB) eval(p *player, x, y int) float64 {
	if isCharWall(x, y) || isCharGhost(x, y) {
		return 30
	}
	if !isThereTargetsAround(x, y) {
		return 20
	}
	return float64(random(1, 10))
}
func isThereTargetsAround(x, y int) bool {
	k := x + 2
	l := y + 2
	for i := lowerLimitZero(x); i <= k; i++ {
		for j := lowerLimitZero(y); j <= l; j++ {
			if isCharTarget(i, j) && isColorWhite(i, j) {
				return true
			}
		}
	}
	return false
}
func lowerLimitZero(i int) int {
	ret := i - 2
	if ret < 0 {
		return 0
	}
	return ret
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
		g.underRune = rune(cell.Ch)
		g.color = cell.Fg
		// 移動先のセルにゴーストをセット
		termbox.SetCell(x, y, rune(chGhost), termbox.ColorRed, termbox.ColorBlack)
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
