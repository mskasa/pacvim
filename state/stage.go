package state

import "time"

// Stage は1ステージの設定を保持する。
// HunterConfig / GhostConfig は state/enemy.go 実装時に追加する。
type Stage struct {
	Level     int
	MapPath   string
	Grid      *Grid
	GameSpeed time.Duration
}
