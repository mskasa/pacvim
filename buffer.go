package main

import (
	"bufio"
	"io"
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type buffer struct {
	lines   []*line
	offset  int
	enemies []iEnemy
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
func (b *buffer) plotStageMap(s stage, p *player) {
	firstFlg := true
	width, _ := termbox.Size()
	height := len(b.lines)
	for y := 0; y < height; y++ {
		for x := b.offset; x < width; x++ {
			if isCharTarget(x, y) {
				if firstFlg {
					b.firstTargetY = y
					firstFlg = false
				}
				b.lastTargetY = y
				targetScore++
			} else if isCharPlayer(x, y) {
				p.x = x
				p.y = y
				termbox.SetCell(p.x, p.y, ' ', termbox.ColorWhite, termbox.ColorBlack)
				termbox.SetCursor(p.x, p.y)
			} else if isCharBorder(x, y) {
				termbox.SetCell(x, y, chBorder, termbox.ColorYellow, termbox.ColorBlack)
			} else if isCharObstacle(x, y) {
				var r rune
				if isCharObstacle(x-1, y) || isCharObstacle(x+1, y) {
					r = chObstacle1
				} else if isCharObstacle(x, y-1) || isCharObstacle(x, y+1) {
					r = chObstacle2
				}
				termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
			} else if isCharPoison(x, y) {
				termbox.SetCell(x, y, chPoison, termbox.ColorMagenta, termbox.ColorBlack)
			} else if isCharHunter(x, y) {
				hunter := s.hunter
				hunter.x = x
				hunter.y = y
				termbox.SetCell(hunter.x, hunter.y, hunter.char, hunter.color, termbox.ColorBlack)
				b.enemies = append(b.enemies, &hunter)
			} else if isCharGhost(x, y) {
				ghost := s.ghost
				ghost.x = x
				ghost.y = y
				termbox.SetCell(ghost.x, ghost.y, ghost.char, ghost.color, termbox.ColorBlack)
				b.enemies = append(b.enemies, &ghost)
			}
		}
	}
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
