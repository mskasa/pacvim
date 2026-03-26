package main

import (
	"errors"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/masahiro-kasatani/pacvim/input"
	"github.com/masahiro-kasatani/pacvim/renderer"
	"github.com/masahiro-kasatani/pacvim/state"
)

// Game は ebiten.Game インターフェースを実装する。
type Game struct {
	gs  *state.GameState
	in  *input.Handler
	ren *renderer.Renderer
}

// NewGame はゲームを初期化して返す。
func NewGame() (*Game, error) {
	gs, err := state.NewGameState(3)
	if err != nil {
		return nil, err
	}
	return &Game{
		gs:  gs,
		in:  &input.Handler{},
		ren: renderer.New(),
	}, nil
}

// Update は毎フレーム呼ばれ、入力を受け取ってゲーム状態を更新する。
func (g *Game) Update() error {
	cmd, ch := g.in.Read()
	if err := g.gs.Update(cmd, ch); err != nil {
		if errors.Is(err, state.ErrQuit) {
			return ebiten.Termination
		}
		return err
	}

	// GameOver/GameClear 後は任意キーで終了
	if g.gs.Phase == state.PhaseGameOver || g.gs.Phase == state.PhaseGameClear {
		if cmd != input.CmdNone {
			return ebiten.Termination
		}
	}

	return nil
}

// Draw は毎フレーム呼ばれ、現在のゲーム状態を描画する。state を変更しない。
func (g *Game) Draw(screen *ebiten.Image) {
	g.ren.Draw(screen, g.gs)
}

// Layout はゲームの論理画面サイズを返す。
func (g *Game) Layout(_, _ int) (int, int) {
	return renderer.ScreenWidth, renderer.ScreenHeight
}
