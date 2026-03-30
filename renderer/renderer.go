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
	TileSize     = 32
	OffsetX      = 48 // 行番号エリアの幅（ピクセル）
	OffsetY      = 16 // 上部マージン
	ScreenWidth  = OffsetX + 29*TileSize // 48 + 928 = 976
	ScreenHeight = OffsetY + 15*TileSize + 28 // 16 + 480 + 28 = 524

	statusBarY = ScreenHeight - 28
)

// 色定義
var (
	colorBackground = color.RGBA{10, 10, 20, 255}
	colorBoundary   = color.RGBA{40, 40, 80, 255}
	colorWall       = color.RGBA{90, 90, 110, 255}
	colorLineNum    = color.RGBA{100, 100, 140, 255}
	colorStatusBg   = color.RGBA{30, 30, 50, 255}
	colorStatusText = color.RGBA{200, 200, 200, 255}
	colorOverlay    = color.RGBA{0, 0, 0, 180}
	colorTitle      = color.RGBA{255, 255, 100, 255}
	colorHint       = color.RGBA{150, 150, 200, 255}
)

// Renderer は Draw() の呼び出しごとにゲーム状態を画面に描画する。
type Renderer struct {
	face   *textv2.GoXFace
	sheets *spriteSheet // スプライトシート
	tick   int          // アニメーション用フレームカウンタ
}

// New は Renderer を生成する。スプライトシートのロードに失敗した場合は error を返す。
func New() (*Renderer, error) {
	ss, err := newSpriteSheet()
	if err != nil {
		return nil, err
	}
	return &Renderer{
		face:   textv2.NewGoXFace(basicfont.Face7x13),
		sheets: ss,
	}, nil
}

// Draw はゲーム全体を描画する。state を変更しない。
func (r *Renderer) Draw(screen *ebiten.Image, gs *state.GameState) {
	r.tick++
	screen.Fill(colorBackground)

	switch gs.Phase {
	case state.PhaseOpening:
		r.drawGame(screen, gs)
		r.drawOpeningOverlay(screen)
	case state.PhaseReady:
		r.drawGame(screen, gs)
		r.drawReadyOverlay(screen, gs)
	case state.PhasePlaying:
		r.drawGame(screen, gs)
	case state.PhaseDead:
		r.drawGame(screen, gs)
		r.drawDeadOverlay(screen)
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
				r.drawSprite16(screen, r.sheets.wall, x, y)
			case state.CellApple:
				r.drawSprite16(screen, r.sheets.apple, x, y)
			case state.CellAppleEaten:
				r.drawSprite16(screen, r.sheets.appleEaten, x, y)
			case state.CellPoison:
				r.drawSprite16(screen, r.sheets.poison, x, y)
			}
		}
	}
}

func (r *Renderer) drawLineNumbers(screen *ebiten.Image, grid *state.Grid) {
	for y := 0; y < grid.Height; y++ {
		_, sy := gridToScreen(0, y)
		num := fmt.Sprintf("%2d", y+1)
		r.drawChar(screen, num, 2, sy+TileSize/2+4, colorLineNum)
	}
}

func (r *Renderer) drawEnemies(screen *ebiten.Image, enemies []state.Enemy) {
	frame := (r.tick / 20) % 4
	for _, e := range enemies {
		x, y := e.Position()
		switch e.Kind() {
		case state.EnemyHunter:
			r.drawSprite32(screen, r.sheets.hunterFrames[frame], x, y)
		case state.EnemyGhost:
			r.drawSprite32(screen, r.sheets.ghostFrames[frame], x, y)
		}
	}
}

func (r *Renderer) drawPlayer(screen *ebiten.Image, p *state.Player) {
	if p == nil {
		return
	}
	frame := (r.tick / 20) % 4
	r.drawSprite32(screen, r.sheets.gopherFrames[frame], p.X, p.Y)
}

// drawSprite32 は 32×32 スプライトをグリッド座標 (gx, gy) のタイルにぴったり描画する。
// TileSize=32 なのでスケール不要。
func (r *Renderer) drawSprite32(screen *ebiten.Image, sprite *ebiten.Image, gx, gy int) {
	sx, sy := gridToScreen(gx, gy)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(sx), float64(sy))
	screen.DrawImage(sprite, op)
}

// drawSprite16 は 16×16 スプライトを 2 倍に拡大して TileSize(32) に合わせて描画する。
func (r *Renderer) drawSprite16(screen *ebiten.Image, sprite *ebiten.Image, gx, gy int) {
	sx, sy := gridToScreen(gx, gy)
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(2, 2)
	op.GeoM.Translate(float64(sx), float64(sy))
	screen.DrawImage(sprite, op)
}

func (r *Renderer) drawStatusBar(screen *ebiten.Image, gs *state.GameState) {
	vector.FillRect(screen, 0, float32(statusBarY), ScreenWidth, 28, colorStatusBg, false)
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

func (r *Renderer) drawReadyOverlay(screen *ebiten.Image, gs *state.GameState) {
	vector.FillRect(screen, 0, 0, ScreenWidth, ScreenHeight, colorOverlay, false)
	cx := ScreenWidth / 2
	cy := ScreenHeight / 2
	stage := gs.Stage()
	r.drawCharCentered(screen, fmt.Sprintf("Level %d", stage.Level), cx, cy-30, colorTitle)
	r.drawCharCentered(screen, fmt.Sprintf("Life: %d", gs.Life), cx, cy-10, colorHint)
	r.drawCharCentered(screen, "Press any key to start", cx, cy+16, colorStatusText)
	r.drawCharCentered(screen, "q: quit", cx, cy+36, colorHint)
}

func (r *Renderer) drawDeadOverlay(screen *ebiten.Image) {
	vector.FillRect(screen, 0, 0, ScreenWidth, ScreenHeight, color.RGBA{120, 0, 0, 80}, false)
	cx := ScreenWidth / 2
	cy := ScreenHeight / 2
	r.drawCharCentered(screen, "MISS...", cx, cy, colorTitle)
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
