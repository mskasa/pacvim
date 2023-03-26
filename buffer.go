package main

import (
	"bufio"
	"errors"
	"io"
	"math"
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type buffer struct {
	lines  []*line
	offset int
}

type line struct {
	text []rune
}

func (b *buffer) save(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := new(line)
		l.text = []rune(scanner.Text())
		b.lines = append(b.lines, l)
	}
}

func (b *buffer) plotStageMap() {
	firstFlg := true
	winWidth, _ := termbox.Size()
	textHeight := len(b.lines)
	for y := 0; y < textHeight; y++ {
		for x := b.offset; x < winWidth; x++ {
			if isCharTarget(x, y) {
				if firstFlg {
					firstTargetY = y
					firstFlg = false
				}
				lastTargetY = y
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
	score := []rune("score:" + strconv.Itoa(score) + "/" + strconv.Itoa(targetScore))
	for x, r := range score {
		termbox.SetCell(x, textHeight, r, termbox.ColorRed, termbox.ColorBlack)
	}
}

func (b *buffer) plotSubInfo() {
	textMap := map[int]string{
		0: "Level:" + strconv.Itoa(level),
		1: "Life:" + strconv.Itoa(life),
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

// the 1st one:	2nd quadrant, tactics1
// the 2nd one:	4th quadrant, tactics1
// the 3rd one:	1st quadrant, tactics2
// the 4th one:	3rd quadrant, tactics2
func (b *buffer) plotGhost() ([]*ghost, error) {
	var err error
	var gList []*ghost
	var gPlotRangeList = [][]float64{{0.4, 0.4}, {0.6, 0.6}, {0.6, 0.4}, {0.4, 0.6}}

	for i := 0; i < numOfGhosts(); i++ {
		g := new(ghost)
		g.tactics = i/2 + 1

		j := 0
		for {
			yPlotRangeUpperLimit := len(b.lines) - 1
			yPlotRangeBorder := int(float64(yPlotRangeUpperLimit) * gPlotRangeList[i][1])
			gY := ghostPlotPosition(yPlotRangeBorder, yPlotRangeUpperLimit)
			xPlotRangeUpperLimit := len(b.lines[gY].text) + b.offset
			xPlotRangeBorder := int(float64(xPlotRangeUpperLimit) * gPlotRangeList[i][0])
			gX := ghostPlotPosition(xPlotRangeBorder, xPlotRangeUpperLimit)

			if isCharTarget(gX, gY) && g.move(gX, gY) {
				gList = append(gList, g)
				break
			}

			j++
			if j == 10000 {
				return nil, errors.New("Play with maps that have enough targets in the ghostplot range!")
			}
		}
	}
	return gList, err
}
func ghostPlotPosition(min, max int) int {
	if max-min > min {
		return random(0, min)
	}
	return random(min, max)
}
func numOfGhosts() int {
	numOfGhost := int(math.Ceil(float64(level)/3.0)) + 1
	if numOfGhost > maxNumOfGhost {
		numOfGhost = maxNumOfGhost
	}
	return numOfGhost
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
