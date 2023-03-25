package main

import (
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type Window struct {
	lines []*line
}

func (w *Window) copyBufToWindow(b *buffer) {
	w.lines = []*line{}
	winWidth, winHeight := termbox.Size()
	for i := 0; i < len(b.lines); i++ {
		if i > winHeight-1 {
			break
		}
		w.lines = append(w.lines, &line{})
		for j := 0; j < len(b.lines[i].text); j++ {
			if j+getDigit(b.getTotalNumOfLines())+1 > winWidth-1 {
				break
			}
			w.lines[i].text = append(w.lines[i].text, b.lines[i].text[j])
		}
	}
}

// 行番号なしで表示
func (w *Window) Show(b *buffer) error {
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		return err
	}
	for y, l := range w.lines {
		for x, r := range l.text {
			termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
		}
	}
	return err
}

// 行番号ありで表示
func (w *Window) ShowWithLineNum(b *buffer) error {
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	offset := getDigit(b.getTotalNumOfLines())
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
func makeLineNum(num int, digit int) []rune {
	numstr := strconv.Itoa(num)
	lineNum := make([]rune, digit+1)
	for i := 0; i < len(lineNum); i++ {
		lineNum[i] = ' '
	}
	cdigit := getDigit(num)
	for i, c := range numstr {
		lineNum[i+(digit-cdigit)] = c
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
