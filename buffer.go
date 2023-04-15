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

type window struct {
	lines []*line
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

func createWindow(b *buffer) *window {
	w := new(window)
	for i := 0; i < len(b.lines); i++ {
		w.lines = append(w.lines, &line{})
		for j := 0; j < len(b.lines[i].text); j++ {
			w.lines[i].text = append(w.lines[i].text, b.lines[i].text[j])
		}
	}
	return w
}

func (w *window) show(b *buffer) error {
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	maxDigit := getDigit(len(b.lines))
	b.offset = getOffset(len(b.lines))
	for y, l := range w.lines {
		linenums := makeLineNum(y+1, maxDigit, b.offset)
		t := append(linenums, l.text...)
		for x, r := range t {
			termbox.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	return err
}
func makeLineNum(num int, maxDigit int, maxOffset int) []rune {
	lineNum := make([]rune, maxOffset)
	for i := 0; i < len(lineNum); i++ {
		// fill with spaces
		lineNum[i] = ' '
	}
	numstr := strconv.Itoa(num)
	currentDigit := getDigit(num)
	for i, v := range numstr {
		// align right
		lineNum[i+(maxDigit-currentDigit)] = v
	}
	return lineNum
}
func getDigit(linenum int) int {
	d := 0
	for linenum != 0 {
		linenum = linenum / 10
		d++
	}
	return d
}
func getOffset(linenum int) int {
	// 1 is a half-width space to the right of the line number
	return getDigit(linenum) + 1
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
