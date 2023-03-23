package main

import (
	"math"
	"strconv"
	"time"
)

const (
	Pose int = iota
	Continuing
	GameClear
	GameOver
	StageWin
	StageLose

	maxLevel      = 10
	maxNumOfGhost = 4
)

var (
	GameState       = 0
	targetScore     = 0
	score           = 0
	hogeLevel       = 1
	life            = 3
	inputNum        = 0
	gHasBeenEntered = false
	gameSpeed       time.Duration
)

// ゲームの終了判定（true:ゲーム終了）
func IsFinished() bool {
	// ステージ終了判定
	switch GameState {
	case StageWin:
		// ステージクリアの場合、ゲームレベルを上げる
		hogeLevel++
		if hogeLevel == maxLevel {
			// レベルが最大レベルに達した場合、ゲームクリア
			GameState = GameClear
		}
	case StageLose:
		// ステージ失敗の場合、残機を減らす
		life--
		if life < 0 {
			return true
		}
	case GameClear, GameOver:
		return true
	}
	return false
}

// スコアを初期化する
func Reset() {
	GameState = Pose
	score = 0
	targetScore = 0
}

// ゲームレベルから決定したゴーストの数を返却する
func GetNumOfGhost() int {
	numOfGhost := int(math.Ceil(float64(hogeLevel)/3.0)) + 1
	if numOfGhost > maxNumOfGhost {
		numOfGhost = maxNumOfGhost
	}
	return numOfGhost
}

// マップ読み込み時に目標スコアを算出する
func AddTargetScore() {
	targetScore++
}

// 現在スコアを加算する
func AddScore() {
	score++
	if score == targetScore {
		// 目標スコアに達した場合、ステージクリア
		GameState = StageWin
	}
}

// 現在スコアを取得する
func GetScore() string {
	return strconv.Itoa(score)
}

// 目標スコアを取得する
func GetTargetScore() string {
	return strconv.Itoa(targetScore)
}

// ゲームレベルを取得する
func GetLevel() string {
	return strconv.Itoa(hogeLevel)
}

// 残機を取得する
func GetLife() string {
	return strconv.Itoa(life)
}

// ゲームスピードを取得する
func GetGameSpeed() time.Duration {
	return gameSpeed
}

// ゲームレベルを設定する
func SetLevel(l int) {
	hogeLevel = l
	gameSpeed = time.Duration(1000-(hogeLevel-1)*100/2) * time.Millisecond
}

// 入力値が数字かどうか判定する
func IsInputNum(r rune) (string, bool) {
	s := string(r)
	i, err := strconv.Atoi(s)
	if err == nil && (i != 0 || (i == 0 && inputNum != 0)) {
		// 数値変換成功かつ入力数値が「0」でない場合
		return s, true
	}
	return s, false
}

// インプット数値を設定する
func SetInputNum(s string) {
	if inputNum != 0 {
		// 既に入力数値がある場合は文字列として数値を足算する（例：1+2=12）
		s = strconv.Itoa(inputNum) + s
	}
	inputNum, _ = strconv.Atoi(s)
}

// インプット数値を取得する
func GetInputNum() int {
	return inputNum
}

// インプット数値を初期化する
func InitInputNum() {
	inputNum = 0
}

// 「g」入力の初期化
func InitInput_g() {
	gHasBeenEntered = false
}

// 「g」入力済み状態にする
func FirstInput_g() {
	gHasBeenEntered = true
}

// 「g」入力済みかどうか判定する
func IsFirstInput_g() bool {
	return gHasBeenEntered
}
