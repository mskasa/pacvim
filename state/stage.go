package state

import (
	"time"

	"github.com/masahiro-kasatani/pacvim/input"
)

// Stage は1ステージの設定と状態を保持する。
type Stage struct {
	Level          int
	MapPath        string
	Grid           *Grid
	HunterConfig   EnemyConfig
	GhostConfig    *EnemyConfig    // ゴーストが存在しないステージでは nil
	GameSpeed      time.Duration
	Theme          string          // 学習テーマ（例: "Basic Movement"）
	Commands       []string        // このステージで練習するコマンドの説明（表示用）
	RestrictedCmds []input.Command // このステージで使用できないコマンド
}

// Load はマップファイルを読み込み、Grid を設定して SpawnPoints を返す。
// GameState がステージを開始するときに呼ぶ。
func (s *Stage) Load() (*SpawnPoints, error) {
	grid, spawns, err := LoadGrid(s.MapPath)
	if err != nil {
		return nil, err
	}
	s.Grid = grid
	return spawns, nil
}

// InitStages は全ステージの設定を返す。Grid は Load() を呼ぶまで nil。
func InitStages() []Stage {
	return []Stage{
		{
			Level:          1,
			MapPath:        "files/stage/map01.txt",
			HunterConfig:   EnemyConfig{Builder: NewHunterBuilder()},
			GameSpeed:      1250 * time.Millisecond,
			Theme:          "Basic Movement",
			Commands: []string{
				"h  move left",
				"l  move right",
				"j  move down",
				"k  move up",
			},
			RestrictedCmds: []input.Command{},
		},
		{
			Level:        2,
			MapPath:      "files/stage/map02.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder().Strategize(&TrickyStrategy{})},
			GhostConfig:  &EnemyConfig{Builder: NewGhostBuilder()},
			GameSpeed:    1000 * time.Millisecond,
			Theme:        "Word Motion",
			Commands: []string{
				"w  move to next word",
				"e  move to end of word",
				"b  move to previous word",
			},
			RestrictedCmds: []input.Command{
				input.CmdMoveLeft,
				input.CmdMoveRight,
			},
		},
		{
			Level:        3,
			MapPath:      "files/stage/map03.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GhostConfig:  &EnemyConfig{Builder: NewGhostBuilder()},
			GameSpeed:    1000 * time.Millisecond,
			Theme:        "Line & File Motion",
			Commands: []string{
				"0   move to start of line",
				"$   move to end of line",
				"^   move to first word on line",
				"gg  move to first line",
				"G   move to last line  (NG: line N)",
			},
			RestrictedCmds: []input.Command{},
		},
		{
			Level:        4,
			MapPath:      "files/stage/map04.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GameSpeed:    750 * time.Millisecond,
			Theme:        "Find Character",
			Commands: []string{
				"f{c}  jump to character c",
				"t{c}  jump before character c",
				";     repeat last find",
				",     repeat last find (reverse)",
			},
			RestrictedCmds: []input.Command{
				input.CmdWordForward,
				input.CmdWordEnd,
				input.CmdWordBack,
			},
		},
		{
			Level:        5,
			MapPath:      "files/stage/map05.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder().Strategize(&TrickyStrategy{})},
			GameSpeed:    750 * time.Millisecond,
			Theme:        "All Commands",
			Commands: []string{
				"Use everything you've learned!",
				"h/j/k/l  w/e/b  0/$/^",
				"gg/G  f/t/;/,",
			},
			RestrictedCmds: []input.Command{},
		},
	}
}
