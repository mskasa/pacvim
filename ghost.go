package main

import (
	"math"
	"math/rand"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type Ghost struct {
	x         int
	y         int
	underRune rune
	color     termbox.Attribute
	tactics   int
}

// ゴーストの制御
func Control(ch chan bool, p *Player, gList []*Ghost) {
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

// TODO
func random(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}

// ゴーストを行動させる
func (g *Ghost) action(wg *sync.WaitGroup, p *Player) {
	defer wg.Done()

	var (
		up    float64
		down  float64
		left  float64
		right float64
	)
	// 移動のための評価値算出
	if g.tactics == 1 {
		up = eval(p, g.x, g.y-1)
		down = eval(p, g.x, g.y+1)
		left = eval(p, g.x-1, g.y)
		right = eval(p, g.x+1, g.y)
	} else {
		up = eval2(p, g.x, g.y-1)
		down = eval2(p, g.x, g.y+1)
		left = eval2(p, g.x-1, g.y)
		right = eval2(p, g.x+1, g.y)
	}

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

// プレイヤー追跡タイプ
func eval(p *Player, x, y int) float64 {
	if isCharWall(x, y) || isCharGhost(x, y) {
		// 移動先が壁もしくはゴーストの場合は移動先から除外（十分に大きな値を返却）
		return 1000
	}
	// X軸の距離とY軸の距離それぞれの二乗の和の平方根
	return math.Sqrt(math.Pow(float64(p.Y-y), 2) + math.Pow(float64(p.X-x), 2))
}

// 待ち伏せ徘徊タイプ
func eval2(p *Player, x, y int) float64 {
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
func (g *Ghost) move(x, y int) bool {
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
func (g *Ghost) hasCaptured(p *Player) bool {
	if g.x == (p.X) && g.y == p.Y {
		// プレイヤーカーソルとゴーストの座標が一致した場合
		return true
	}
	return false
}
