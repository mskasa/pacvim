package main

import (
	"bufio"
	"io"
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type buffer struct {
	lines  []*line
	offset int
	// For command gg or G
	firstTargetY int
	lastTargetY  int
}

type line struct {
	text []rune
}

// Save .txt to buffer.
func createBuffer(reader io.Reader) *buffer {
	b := new(buffer)
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := new(line)
		l.text = []rune(scanner.Text())
		b.lines = append(b.lines, l)
	}
	return b
}

func (b *buffer) plotScore() {
	position := len(b.lines)
	s := []rune("score: " + strconv.Itoa(score) + "/" + strconv.Itoa(targetScore))
	for x, r := range s {
		termbox.SetCell(x, position, r, termbox.ColorGreen, termbox.ColorBlack)
	}
}

func (b *buffer) plotSubInfo(level int, life int) {
	textMap := map[int]string{
		0: "Level: " + strconv.Itoa(level),
		1: "Life : " + strconv.Itoa(life),
		2: "PRESS ENTER TO PLAY!",
		3: "q TO EXIT!"}
	position := len(b.lines) + 1
	for i := 0; i < len(textMap); i++ {
		for x, r := range []rune(textMap[i]) {
			termbox.SetCell(x, position, r, termbox.ColorWhite, termbox.ColorBlack)
		}
		position++
	}
}

func isCharBorder(x, y int) bool {
	return isChar(x, y, chBorder)
}
func isCharObstacle(x, y int) bool {
	return isChar(x, y, chObstacle1) || isChar(x, y, chObstacle2) || isChar(x, y, chObstacle3)
}
func isCharWall(x, y int) bool {
	return isCharObstacle(x, y) || isCharBorder(x, y)
}
func isCharPlayer(x, y int) bool {
	return isChar(x, y, chPlayer)
}
func isCharHunter(x, y int) bool {
	return isChar(x, y, chHunter)
}
func isCharGhost(x, y int) bool {
	return isChar(x, y, chGhost)
}
func isCharEnemy(x, y int) bool {
	return isCharHunter(x, y) || isCharGhost(x, y)
}
func isCharSpace(x, y int) bool {
	return isChar(x, y, ' ')
}
func isCharTarget(x, y int) bool {
	return isChar(x, y, chTarget)
}
func isCharPoison(x, y int) bool {
	return isChar(x, y, chPoison)
}
func isChar(x, y int, r rune) bool {
	winWidth, _ := termbox.Size()
	cell := termbox.CellBuffer()[(winWidth*y)+x]
	return r == cell.Ch
}
func isColorWhite(x, y int) bool {
	winWidth, _ := termbox.Size()
	cell := termbox.CellBuffer()[(winWidth*y)+x]
	return cell.Fg == termbox.ColorWhite
}
