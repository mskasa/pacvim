package main

import (
	"math"
	"strconv"
	"time"
)

const (
	Pose int = iota
	Continuing
	Quit
	StageWin
	StageLose

	maxLevel      = 10
	maxNumOfGhost = 4
)

var (
	gameState           = 0
	targetScore         = 0
	score               = 0
	hogeLevel           = 1
	life                = 3
	inputNum            = 0
	isLowercaseGEntered = false
	gameSpeed           time.Duration
)

// スコアを初期化する
func Reset() {
	gameState = Pose
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

// ゲームレベルを取得する
func GetLevel() string {
	return strconv.Itoa(hogeLevel)
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
