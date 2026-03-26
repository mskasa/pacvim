// Package renderer は Ebitengine を使ってゲーム状態を描画する。
// state パッケージを読み取るだけで変更しない。
package renderer

import (
	"fmt"
	"image/color"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	textv2 "github.com/hajimehoshi/ebiten/v2/text/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"

	"github.com/masahiro-kasatani/pacvim/state"
)

const (
	TileSize     = 16
	OffsetX      = 40 // 行番号エリアの幅（ピクセル）
	OffsetY      = 16 // 上部マージン
	ScreenWidth  = 800
	ScreenHeight = 600

	statusBarY = ScreenHeight - 24
)

// 色定義
var (
	colorBackground = color.RGBA{10, 10, 20, 255}
	colorBoundary   = color.RGBA{40, 40, 80, 255}
	colorWall       = color.RGBA{90, 90, 110, 255}
	colorApple      = color.RGBA{80, 200, 80, 255}
	colorAppleEaten = color.RGBA{30, 60, 30, 255}
	colorPoison     = color.RGBA{200, 50, 50, 255}
	colorPlayer     = color.RGBA{255, 255, 100, 255}
	colorHunter     = color.RGBA{255, 140, 0, 255}
	colorGhost      = color.RGBA{180, 100, 255, 255}
	colorLineNum    = color.RGBA{100, 100, 140, 255}
	colorStatusBg   = color.RGBA{30, 30, 50, 255}
	colorStatusText = color.RGBA{200, 200, 200, 255}
	colorOverlay    = color.RGBA{0, 0, 0, 180}
	colorTitle      = color.RGBA{255, 255, 100, 255}
	colorHint       = color.RGBA{150, 150, 200, 255}
)

// Renderer は Draw() の呼び出しごとにゲーム状態を画面に描画する。
type Renderer struct {
	face *textv2.GoXFace
}

// New は Renderer を生成する。
func New() *Renderer {
	return &Renderer{
		face: textv2.NewGoXFace(basicfont.Face7x13),
	}
}

// Draw はゲーム全体を描画する。state を変更しない。
func (r *Renderer) Draw(screen *ebiten.Image, gs *state.GameState) {
	screen.Fill(colorBackground)

	switch gs.Phase {
	case state.PhaseOpening:
		r.drawGame(screen, gs)
		r.drawOpeningOverlay(screen)
	case state.PhasePlaying:
		r.drawGame(screen, gs)
	case state.PhaseGameOver:
		r.drawGame(screen, gs)
		r.drawMessageOverlay(screen, "GAME OVER", "press any key to quit")
	case state.PhaseGameClear:
		r.drawGame(screen, gs)
		r.drawMessageOverlay(screen, "YOU WIN!", "press any key to quit")
	}
}

// --- ゲーム画面 ---

func (r *Renderer) drawGame(screen *ebiten.Image, gs *state.GameState) {
	stage := gs.Stage()
	if stage.Grid == nil {
		return
	}
	r.drawGrid(screen, stage.Grid)
	r.drawLineNumbers(screen, stage.Grid)
	r.drawEnemies(screen, gs.Enemies)
	r.drawPlayer(screen, gs.Player)
	r.drawStatusBar(screen, gs)
}

func (r *Renderer) drawGrid(screen *ebiten.Image, grid *state.Grid) {
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			sx, sy := gridToScreen(x, y)
			cell := grid.At(x, y)
			switch cell.Kind {
			case state.CellBoundary:
				vector.FillRect(screen, float32(sx), float32(sy), TileSize, TileSize, colorBoundary, false)
				r.drawChar(screen, "+", sx+4, sy+11, colorWall)
			case state.CellWall:
				vector.FillRect(screen, float32(sx), float32(sy), TileSize, TileSize, colorWall, false)
			case state.CellApple:
				r.drawChar(screen, "o", sx+4, sy+11, colorApple)
			case state.CellAppleEaten:
				r.drawChar(screen, "·", sx+4, sy+11, colorAppleEaten)
			case state.CellPoison:
				r.drawChar(screen, "X", sx+3, sy+11, colorPoison)
			}
		}
	}
}

func (r *Renderer) drawLineNumbers(screen *ebiten.Image, grid *state.Grid) {
	for y := 0; y < grid.Height; y++ {
		_, sy := gridToScreen(0, y)
		num := fmt.Sprintf("%2d", y+1)
		r.drawChar(screen, num, 2, sy+11, colorLineNum)
	}
}

func (r *Renderer) drawEnemies(screen *ebiten.Image, enemies []state.Enemy) {
	for _, e := range enemies {
		x, y := e.Position()
		sx, sy := gridToScreen(x, y)
		switch e.Kind() {
		case state.EnemyHunter:
			vector.FillRect(screen, float32(sx)+2, float32(sy)+2, TileSize-4, TileSize-4, colorHunter, false)
			r.drawChar(screen, "H", sx+4, sy+11, colorBackground)
		case state.EnemyGhost:
			vector.FillRect(screen, float32(sx)+2, float32(sy)+2, TileSize-4, TileSize-4, colorGhost, false)
			r.drawChar(screen, "G", sx+4, sy+11, colorBackground)
		}
	}
}

func (r *Renderer) drawPlayer(screen *ebiten.Image, p *state.Player) {
	if p == nil {
		return
	}
	sx, sy := gridToScreen(p.X, p.Y)
	vector.FillRect(screen, float32(sx)+1, float32(sy)+1, TileSize-2, TileSize-2, colorPlayer, false)
	r.drawChar(screen, "P", sx+4, sy+11, colorBackground)
}

func (r *Renderer) drawStatusBar(screen *ebiten.Image, gs *state.GameState) {
	vector.FillRect(screen, 0, float32(statusBarY), ScreenWidth, 24, colorStatusBg, false)
	stage := gs.Stage()
	text := fmt.Sprintf("  Level: %d    Score: %d/%d    Life: %d",
		stage.Level, gs.Player.Score, gs.Player.TargetScore, gs.Life)
	r.drawChar(screen, text, 4, statusBarY+15, colorStatusText)
}

// --- オーバーレイ ---

func (r *Renderer) drawOpeningOverlay(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, ScreenWidth, ScreenHeight, colorOverlay, false)
	cx := ScreenWidth / 2
	cy := ScreenHeight / 2
	r.drawCharCentered(screen, "PacVim", cx, cy-30, colorTitle)
	r.drawCharCentered(screen, "Learn Vim by playing Pac-Man!", cx, cy-10, colorHint)
	r.drawCharCentered(screen, "Press any key to start", cx, cy+16, colorStatusText)
	r.drawCharCentered(screen, "q: quit", cx, cy+36, colorHint)
}

func (r *Renderer) drawMessageOverlay(screen *ebiten.Image, title, hint string) {
	vector.FillRect(screen, 0, 0, ScreenWidth, ScreenHeight, colorOverlay, false)
	cx := ScreenWidth / 2
	cy := ScreenHeight / 2
	r.drawCharCentered(screen, title, cx, cy-16, colorTitle)
	r.drawCharCentered(screen, hint, cx, cy+10, colorHint)
}

// --- テキスト描画ヘルパー ---

func (r *Renderer) drawChar(screen *ebiten.Image, s string, x, y int, clr color.Color) {
	op := &textv2.DrawOptions{}
	op.GeoM.Translate(float64(x), float64(y))
	op.ColorScale.ScaleWithColor(clr)
	textv2.Draw(screen, s, r.face, op)
}

func (r *Renderer) drawCharCentered(screen *ebiten.Image, s string, cx, y int, clr color.Color) {
	w, _ := textv2.Measure(s, r.face, 0)
	r.drawChar(screen, s, cx-int(w/2), y, clr)
}

// --- 座標変換 ---

// gridToScreen はグリッド座標をスクリーン座標に変換する。
func gridToScreen(gx, gy int) (int, int) {
	return OffsetX + gx*TileSize, OffsetY + gy*TileSize
}

// DebugPrint はデバッグ情報を画面右下に表示する。
func DebugPrint(screen *ebiten.Image) {
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("TPS: %.0f", ebiten.ActualTPS()), ScreenWidth-80, ScreenHeight-20)
}
