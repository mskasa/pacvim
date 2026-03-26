package state

import "time"

// Stage は1ステージの設定を保持する。
type Stage struct {
	Level        int
	MapPath      string
	Grid         *Grid
	HunterConfig EnemyConfig
	GhostConfig  *EnemyConfig // ゴーストが存在しないステージでは nil
	GameSpeed    time.Duration
}
