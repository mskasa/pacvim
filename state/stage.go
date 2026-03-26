package state

import "time"

// Stage は1ステージの設定と状態を保持する。
type Stage struct {
	Level        int
	MapPath      string
	Grid         *Grid
	HunterConfig EnemyConfig
	GhostConfig  *EnemyConfig // ゴーストが存在しないステージでは nil
	GameSpeed    time.Duration
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
			Level:        1,
			MapPath:      "files/stage/map01.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GameSpeed:    1250 * time.Millisecond,
		},
		{
			Level:        2,
			MapPath:      "files/stage/map02.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder().Strategize(&TrickyStrategy{})},
			GhostConfig:  &EnemyConfig{Builder: NewGhostBuilder()},
			GameSpeed:    1000 * time.Millisecond,
		},
		{
			Level:        3,
			MapPath:      "files/stage/map03.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GhostConfig:  &EnemyConfig{Builder: NewGhostBuilder()},
			GameSpeed:    1000 * time.Millisecond,
		},
		{
			Level:        4,
			MapPath:      "files/stage/map04.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GameSpeed:    750 * time.Millisecond,
		},
		{
			Level:        5,
			MapPath:      "files/stage/map05.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder().Strategize(&TrickyStrategy{})},
			GameSpeed:    750 * time.Millisecond,
		},
	}
}
