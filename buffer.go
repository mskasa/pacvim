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

// Convert characters
// Color characters
// Save firstTargetY, lastTargetY and targetScore
func (b *buffer) plotStageMap() {
	firstFlg := true
	winWidth, _ := termbox.Size()
	textHeight := len(b.lines)
	for y := 0; y < textHeight; y++ {
		for x := b.offset; x < winWidth; x++ {
			if isCharTarget(x, y) {
				if firstFlg {
					b.firstTargetY = y
					firstFlg = false
				}
				b.lastTargetY = y
				targetScore++
			} else if isCharWall(x, y) {
				r := '-'
				if y == 0 || y == len(b.lines) {
					r = '-'
				} else if isCharWall(x, y-1) && isCharWall(x, y+1) {
					r = '|'
				} else if isCharWall(x-1, y) || isCharWall(x+1, y) {
					r = '-'
				} else if isCharWall(x, y-1) || isCharWall(x, y+1) {
					r = '|'
				}
				termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
			} else if isCharPoison(x, y) {
				termbox.SetCell(x, y, chPoison, termbox.ColorMagenta, termbox.ColorBlack)
			}
		}
	}
}

func (b *buffer) plotScore() {
	textHeight := len(b.lines)
	s := []rune("score: " + strconv.Itoa(score) + "/" + strconv.Itoa(targetScore))
	for x, r := range s {
		termbox.SetCell(x, textHeight, r, termbox.ColorGreen, termbox.ColorBlack)
	}
}

func (b *buffer) plotSubInfo(level int, life int) {
	textMap := map[int]string{
		0: "Level: " + strconv.Itoa(level),
		1: "Life : " + strconv.Itoa(life),
		2: "PRESS ENTER TO PLAY!",
		3: "q TO EXIT!"}
	textHeight := len(b.lines) + 1
	for i := 0; i < len(textMap); i++ {
		for x, r := range []rune(textMap[i]) {
			termbox.SetCell(x, textHeight, r, termbox.ColorWhite, termbox.ColorBlack)
		}
		textHeight++
	}
}

func isCharWall(x, y int) bool {
	return isChar(x, y, chWall1) || isChar(x, y, chWall2) || isChar(x, y, chWall3)
}
func isCharGhost(x, y int) bool {
	return isChar(x, y, chGhost)
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
