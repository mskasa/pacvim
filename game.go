package main

import (
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
