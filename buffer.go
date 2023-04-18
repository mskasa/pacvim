package main

import (
	"bufio"
	"bytes"
	"io"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type buffer struct {
	lines  []*line
	offset int
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
		lineNum[i] = chSpace
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

func switchScene(fileName string) error {
	termbox.HideCursor()
	f, err := static.ReadFile(fileName)
	if err != nil {
		return err
	}

	b := createBuffer(bytes.NewReader(f))
	w := createWindow(b)
	if err != nil {
		return err
	}

	if err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}
	for y, l := range w.lines {
		for x, r := range l.text {
			termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
		}
	}
	if err = termbox.Flush(); err != nil {
		return err
	}
	time.Sleep(750 * time.Millisecond)
	return err
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
	return isChar(x, y, chSpace)
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
