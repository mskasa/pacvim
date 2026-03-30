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
	Theme        string   // 学習テーマ（例: "基本移動"）
	Commands     []string // このステージで練習するコマンドの説明（表示用）
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
			Theme:        "基本移動",
			Commands: []string{
				"h  左へ移動",
				"l  右へ移動",
				"j  下へ移動",
				"k  上へ移動",
			},
		},
		{
			Level:        2,
			MapPath:      "files/stage/map02.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder().Strategize(&TrickyStrategy{})},
			GhostConfig:  &EnemyConfig{Builder: NewGhostBuilder()},
			GameSpeed:    1000 * time.Millisecond,
			Theme:        "単語移動",
			Commands: []string{
				"w  次の単語の先頭へ",
				"e  単語の末尾へ",
				"b  前の単語の先頭へ",
			},
		},
		{
			Level:        3,
			MapPath:      "files/stage/map03.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GhostConfig:  &EnemyConfig{Builder: NewGhostBuilder()},
			GameSpeed:    1000 * time.Millisecond,
			Theme:        "行・ファイル移動",
			Commands: []string{
				"0   行の先頭へ",
				"$   行の末尾へ",
				"^   行の最初の単語へ",
				"gg  ファイルの先頭へ",
				"G   ファイルの末尾へ（N行目へ）",
			},
		},
		{
			Level:        4,
			MapPath:      "files/stage/map04.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder()},
			GameSpeed:    750 * time.Millisecond,
			Theme:        "文字検索",
			Commands: []string{
				"f{c}  右方向で文字 c へジャンプ",
				"t{c}  右方向で文字 c の手前へ",
				";     直前の検索を繰り返す",
				",     直前の検索を逆方向に繰り返す",
			},
		},
		{
			Level:        5,
			MapPath:      "files/stage/map05.txt",
			HunterConfig: EnemyConfig{Builder: NewHunterBuilder().Strategize(&TrickyStrategy{})},
			GameSpeed:    750 * time.Millisecond,
			Theme:        "総合",
			Commands: []string{
				"全コマンドを駆使してクリアしよう！",
				"h/j/k/l  w/e/b  0/$/^",
				"gg/G  f/t/;/,",
			},
		},
	}
}
