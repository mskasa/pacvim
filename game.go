package main

import (
	"strconv"
	"time"
)

const (
	pose int = iota
	continuing
	quit
	win
	lose

	maxLevel      = 10
	maxNumOfGhost = 4
)

var (
	gameState           = 0
	targetScore         = 0
	score               = 0
	level               = 1
	life                = 3
	inputNum            = 0
	isLowercaseGEntered = false
	gameSpeed           = time.Second
)

// スコアを初期化する
func Reset() {
	gameState = pose
	score = 0
	targetScore = 0
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
