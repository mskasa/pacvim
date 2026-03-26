package state

import (
	"errors"
	"time"

	"github.com/masahiro-kasatani/pacvim/input"
)

// ErrQuit はプレイヤーが q でゲームを終了したときに返す番兵エラー。
var ErrQuit = errors.New("quit")

// GamePhase はゲームの進行フェーズを表す。
type GamePhase int

const (
	PhaseOpening   GamePhase = iota // 開幕画面（最初のキー入力待ち）
	PhasePlaying                    // プレイ中
	PhaseStageClear                 // ステージクリア
	PhaseGameOver                   // ゲームオーバー
	PhaseGameClear                  // 全ステージクリア
)

// GameState はゲーム全体の状態を保持する。
type GameState struct {
	Stages    []Stage
	StageIdx  int
	Player    *Player
	Enemies   []Enemy
	Life      int
	EnemyTick int
	Phase     GamePhase
}

// NewGameState は初期状態の GameState を生成する。
// life には初期ライフ数を指定する。
func NewGameState(life int) (*GameState, error) {
	stages := InitStages()
	gs := &GameState{
		Stages:   stages,
		Life:     life,
		Phase:    PhaseOpening,
	}
	if err := gs.loadStage(); err != nil {
		return nil, err
	}
	return gs, nil
}

// Stage は現在のステージを返す。
func (gs *GameState) Stage() *Stage {
	return &gs.Stages[gs.StageIdx]
}

// loadStage は現在のステージを読み込み、Player と Enemies を初期化する。
func (gs *GameState) loadStage() error {
	stage := gs.Stage()
	spawns, err := stage.Load()
	if err != nil {
		return err
	}

	apples := countApples(stage.Grid)
	gs.Player = &Player{
		X:           spawns.Player.X,
		Y:           spawns.Player.Y,
		State:       PlayerAlive,
		TargetScore: apples,
	}

	enemies := make([]Enemy, 0, len(spawns.Hunters)+len(spawns.Ghosts))
	for _, sp := range spawns.Hunters {
		enemies = append(enemies, stage.HunterConfig.Build(sp.X, sp.Y))
	}
	if stage.GhostConfig != nil {
		for _, sp := range spawns.Ghosts {
			enemies = append(enemies, stage.GhostConfig.Build(sp.X, sp.Y))
		}
	}
	gs.Enemies = enemies
	gs.EnemyTick = 0
	return nil
}

// countApples はグリッド上のリンゴの数を数える。
func countApples(grid *Grid) int {
	count := 0
	for y := 0; y < grid.Height; y++ {
		for x := 0; x < grid.Width; x++ {
			if grid.At(x, y).Kind == CellApple {
				count++
			}
		}
	}
	return count
}

// Update は1フレーム分のゲーム状態を更新する。
// q キーで ErrQuit を返す。
func (gs *GameState) Update(cmd input.Command, ch rune) error {
	if cmd == input.CmdQuit {
		return ErrQuit
	}

	switch gs.Phase {
	case PhaseOpening:
		gs.updateOpening(cmd, ch)
	case PhasePlaying:
		gs.updatePlaying(cmd, ch)
	}
	return nil
}

// updateOpening はオープニング画面の更新処理。
// CmdNone 以外のコマンドが来たらゲームを開始する。
func (gs *GameState) updateOpening(cmd input.Command, ch rune) {
	if cmd != input.CmdNone {
		gs.Phase = PhasePlaying
		gs.updatePlaying(cmd, ch)
	}
}

// updatePlaying はプレイ中の更新処理。
func (gs *GameState) updatePlaying(cmd input.Command, ch rune) {
	gs.Player.Apply(cmd, ch, gs.Stage())

	gs.checkEnemyCapture()
	if gs.Phase != PhasePlaying {
		return
	}

	gs.advanceEnemyTick()
	gs.checkEnemyCapture()
	if gs.Phase != PhasePlaying {
		return
	}

	switch gs.Player.State {
	case PlayerWon:
		gs.handleStageClear()
	case PlayerDead:
		gs.handlePlayerDeath()
	}
}

// advanceEnemyTick は EnemyTick を進め、GameSpeed に達したら敵を移動させる。
func (gs *GameState) advanceEnemyTick() {
	gs.EnemyTick++
	if time.Duration(gs.EnemyTick)*(time.Second/60) >= gs.Stage().GameSpeed {
		gs.moveEnemies()
		gs.EnemyTick = 0
	}
}

// moveEnemies は全敵を1ステップ移動させる。
func (gs *GameState) moveEnemies() {
	for _, e := range gs.Enemies {
		nx, ny := e.Think(gs.Player, gs.Stage().Grid)
		e.Move(nx, ny, gs.Stage().Grid)
	}
}

// checkEnemyCapture はいずれかの敵がプレイヤーを捕獲しているか確認する。
// 捕獲されていた場合は handlePlayerDeath を呼ぶ。
func (gs *GameState) checkEnemyCapture() {
	for _, e := range gs.Enemies {
		if e.HasCaptured(gs.Player) {
			gs.Player.State = PlayerDead
			gs.handlePlayerDeath()
			return
		}
	}
}

// handleStageClear はステージクリア時の処理を行う。
func (gs *GameState) handleStageClear() {
	if gs.StageIdx+1 < len(gs.Stages) {
		gs.StageIdx++
		_ = gs.loadStage()
		gs.Phase = PhasePlaying
	} else {
		gs.Phase = PhaseGameClear
	}
}

// handlePlayerDeath はプレイヤー死亡時の処理を行う。
func (gs *GameState) handlePlayerDeath() {
	gs.Life--
	if gs.Life <= 0 {
		gs.Phase = PhaseGameOver
	} else {
		_ = gs.loadStage()
		gs.Phase = PhasePlaying
	}
}
