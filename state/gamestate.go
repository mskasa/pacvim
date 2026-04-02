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
	PhaseOpening     GamePhase = iota // 開幕画面（最初のキー入力待ち）
	PhaseStageSelect                  // ステージセレクト画面
	PhaseReady                        // 準備画面（ステージ開始前・ミス後）
	PhasePlaying                      // プレイ中
	PhaseDead                         // 死亡演出中（deadDuration フレーム後に PhaseReady へ遷移）
	PhaseGameOver                     // ゲームオーバー
	PhaseGameClear                    // 全ステージクリア
)

// deadDuration は PhaseDead の持続フレーム数（約1.5秒 @ TPS=60）。
const deadDuration = 90

// GameState はゲーム全体の状態を保持する。
type GameState struct {
	Stages          []Stage
	StageIdx        int
	StageSelectIdx  int  // ステージセレクト画面でのカーソル位置
	Player          *Player
	Enemies         []Enemy
	Life            int
	EnemyTick       int
	Phase           GamePhase
	deadTimer       int  // PhaseDead の経過フレーム数
	RestrictedInput bool // 直前のフレームで封印コマンドが入力されたか
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
		gs.updateWaitKey(cmd, PhaseStageSelect)
	case PhaseStageSelect:
		gs.updateStageSelect(cmd)
	case PhaseReady:
		gs.updateWaitKey(cmd, PhasePlaying)
	case PhasePlaying:
		gs.updatePlaying(cmd, ch)
	case PhaseDead:
		gs.updateDead()
	}
	return nil
}

// updateStageSelect はステージセレクト画面の入力を処理する。
// j/k でカーソル移動、l で選択して PhaseReady へ遷移する。
func (gs *GameState) updateStageSelect(cmd input.Command) {
	n := len(gs.Stages)
	switch cmd {
	case input.CmdMoveDown:
		gs.StageSelectIdx = (gs.StageSelectIdx + 1) % n
	case input.CmdMoveUp:
		gs.StageSelectIdx = (gs.StageSelectIdx + n - 1) % n
	case input.CmdMoveRight:
		gs.StageIdx = gs.StageSelectIdx
		_ = gs.loadStage()
		gs.Phase = PhaseReady
	}
}

// updateDead は PhaseDead のフレームを進め、deadDuration 経過後に PhaseReady へ遷移する。
func (gs *GameState) updateDead() {
	gs.deadTimer++
	if gs.deadTimer >= deadDuration {
		_ = gs.loadStage()
		gs.Phase = PhaseReady
	}
}

// updateWaitKey は CmdNone 以外のキー入力で next フェーズへ遷移する。
// キーはゲームコマンドとして処理しない（消費するだけ）。
func (gs *GameState) updateWaitKey(cmd input.Command, next GamePhase) {
	if cmd != input.CmdNone {
		gs.Phase = next
	}
}

// isRestricted は cmd が現在ステージの封印コマンドに含まれているか返す。
func (gs *GameState) isRestricted(cmd input.Command) bool {
	for _, c := range gs.Stage().RestrictedCmds {
		if c == cmd {
			return true
		}
	}
	return false
}

// updatePlaying はプレイ中の更新処理。
func (gs *GameState) updatePlaying(cmd input.Command, ch rune) {
	gs.RestrictedInput = false
	if cmd != input.CmdNone && gs.isRestricted(cmd) {
		gs.RestrictedInput = true
		return
	}
	gs.Player.Apply(cmd, ch, gs.Stage(), gs.Enemies)

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
// 移動先が他の敵の現在地または他の敵の移動先と重複する場合は移動しない。
// インデックスが小さい敵が優先される。
func (gs *GameState) moveEnemies() {
	grid := gs.Stage().Grid

	// 全員の移動先を先に決定する
	targets := make([][2]int, len(gs.Enemies))
	for i, e := range gs.Enemies {
		nx, ny := e.Think(gs.Player, grid)
		targets[i] = [2]int{nx, ny}
	}

	for i, e := range gs.Enemies {
		nx, ny := targets[i][0], targets[i][1]
		conflict := false
		for j, other := range gs.Enemies {
			if i == j {
				continue
			}
			// 他の敵の現在地と重複
			ox, oy := other.Position()
			if nx == ox && ny == oy {
				conflict = true
				break
			}
			// 他の敵の移動先と重複（インデックスが小さい方を優先）
			if nx == targets[j][0] && ny == targets[j][1] && i > j {
				conflict = true
				break
			}
		}
		if !conflict {
			e.Move(nx, ny, grid)
		}
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
		gs.Phase = PhaseReady
	} else {
		gs.Phase = PhaseGameClear
	}
}

// handlePlayerDeath はプレイヤー死亡時の処理を行う。
// ライフが残っている場合は PhaseDead へ遷移し、loadStage は deadDuration 後まで遅らせる。
func (gs *GameState) handlePlayerDeath() {
	gs.Life--
	if gs.Life <= 0 {
		gs.Phase = PhaseGameOver
	} else {
		gs.deadTimer = 0
		gs.Phase = PhaseDead
	}
}
