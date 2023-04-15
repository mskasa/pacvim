package main

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type window struct {
	lines []*line
}

func createWindow(b *buffer) *window {
	w := new(window)
	w.lines = []*line{}
	winWidth, winHeight := termbox.Size()
	for i := 0; i < len(b.lines); i++ {
		if i > winHeight-1 {
			break
		}
		w.lines = append(w.lines, &line{})
		for j := 0; j < len(b.lines[i].text); j++ {
			if j+getDigit(len(b.lines))+1 > winWidth-1 {
				break
			}
			w.lines[i].text = append(w.lines[i].text, b.lines[i].text[j])
		}
	}
	return w
}

func (w *window) show(b *buffer) error {
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	offset := getDigit(len(b.lines))
	b.offset = offset + 1
	for y, l := range w.lines {
		linenums := makeLineNum(y+1, offset)
		t := append(linenums, l.text...)
		for x, r := range t {
			termbox.SetCell(x, y, r, termbox.ColorWhite, termbox.ColorBlack)
		}
	}
	return err
}
func makeLineNum(num int, offset int) []rune {
	numstr := strconv.Itoa(num)
	lineNum := make([]rune, offset+1)
	for i := 0; i < len(lineNum); i++ {
		lineNum[i] = ' '
	}
	digit := getDigit(num)
	for i, v := range numstr {
		lineNum[i+(offset-digit)] = v
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
