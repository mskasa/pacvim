package main

import (
	"bufio"
	"errors"
	"io"
	"math"
	"strconv"

	termbox "github.com/nsf/termbox-go"
)

type Buffer struct {
	Lines  []*Line
	Offset int
}

type Line struct {
	Text []rune
}

func createBuffer() *Buffer {
	b := new(Buffer)
	return b
}

// 対象の行の文字列を取得
func (b *Buffer) getNumOfCharsInTheLine(y int) int {
	return len(b.Lines[y].Text)
}

// 行数の取得
func (b *Buffer) getTotalNumOfLines() int {
	return len(b.Lines)
}

// ファイルをバッファに読み込む
func (b *Buffer) readFileToBuf(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := new(Line)
		l.Text = []rune(scanner.Text())
		b.Lines = append(b.Lines, l)
	}
}

/*
 * マップ読み込み時の処理
 * 1. viコマンド処理に必要な値を保持（firstTargetY、lastTargetY）
 * 2. 目標スコアの設定：game.AddtargetScore()
 * 3. 壁の文字変換：convertWallChar
 * 4. 毒に色を付与：colorThePoison
 */
func (b *Buffer) checkAllChar() {
	firstFlg := true
	winWidth, _ := termbox.Size()
	textHeight := b.getTotalNumOfLines()
	for y := 0; y < textHeight; y++ {
		for x := b.Offset; x < winWidth; x++ {
			if isCharTarget(x, y) {
				if firstFlg {
					firstTargetY = y
					firstFlg = false
				}
				lastTargetY = y
				targetScore++
			}
			b.convertWallChar(x, y)
			colorThePoison(x, y)
		}
	}
}

// 壁を変換（# → | or -）
func (b *Buffer) convertWallChar(x, y int) {
	if isCharWall(x, y) {
		r := '-'
		if y == 0 || y == b.getTotalNumOfLines() {
			r = '-'
		} else if isCharWall(x, y-1) && isCharWall(x, y+1) {
			r = '|'
		} else if isCharWall(x-1, y) || isCharWall(x+1, y) {
			r = '-'
		} else if isCharWall(x, y-1) || isCharWall(x, y+1) {
			r = '|'
		}
		termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
	}
}

// 毒に色を付与
func colorThePoison(x, y int) {
	if isCharPoison(x, y) {
		termbox.SetCell(x, y, chPoison, termbox.ColorMagenta, termbox.ColorBlack)
	}
}

// スコアの表示
func (b *Buffer) displayscore() {
	textHeight := b.getTotalNumOfLines()
	score := []rune("score:" + strconv.Itoa(score) + "/" + strconv.Itoa(targetScore))
	for x, r := range score {
		termbox.SetCell(x, textHeight, r, termbox.ColorRed, termbox.ColorBlack)
	}
}

// 補足情報の表示
func (b *Buffer) displayNote() {
	textMap := map[int]string{
		0: "Level:" + strconv.Itoa(level),
		1: "Life:" + strconv.Itoa(life),
		2: "PRESS ENTER TO PLAY!",
		3: "q TO EXIT!"}
	textHeight := b.getTotalNumOfLines() + 1
	for i := 0; i < len(textMap); i++ {
		for x, r := range []rune(textMap[i]) {
			termbox.SetCell(x, textHeight, r, termbox.ColorWhite, termbox.ColorBlack)
		}
		textHeight++
	}
}

// 文字判定
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

// ゴーストを生み出す
func (b *Buffer) protGhost() ([]*Ghost, error) {
	var err error
	var gList []*Ghost
	// 一体目：第二象限 二体目：第四象限 三体目：第一象限 四体目：第三象限
	var gPlotRangeList = [][]float64{{0.4, 0.4}, {0.6, 0.6}, {0.6, 0.4}, {0.4, 0.6}}

	// ゲームレベルに応じて最大4体のゴーストを生み出す
	for i := 0; i < getNumOfGhost(); i++ {
		g := new(Ghost)
		g.tactics = i/2 + 1

		j := 0
		for {
			// y座標の仮決定（可読性のため敢えて本ブロック内に一連の処理をまとめて記述）
			yPlotRangeUpperLimit := b.getTotalNumOfLines() - 1
			yPlotRangeBorder := int(float64(yPlotRangeUpperLimit) * gPlotRangeList[i][1])
			gY := decidePlotPosition(yPlotRangeBorder, yPlotRangeUpperLimit)
			// x座標の仮決定
			xPlotRangeUpperLimit := b.getNumOfCharsInTheLine(gY) + b.Offset
			xPlotRangeBorder := int(float64(xPlotRangeUpperLimit) * gPlotRangeList[i][0])
			gX := decidePlotPosition(xPlotRangeBorder, xPlotRangeUpperLimit)
			// 仮決定した座標がドットであれば確定
			if isCharTarget(gX, gY) && g.move(gX, gY) {
				gList = append(gList, g)
				break
			}
			// 10000回回してプロット位置が決まらなかった場合
			j++
			if j == 10000 {
				return nil, errors.New("ゴーストプロット範囲にターゲットが十分あるマップで遊んでください")
			}
		}
	}
	return gList, err
}
func decidePlotPosition(min, max int) int {
	if max-min > min {
		return random(0, min)
	}
	return random(min, max)
}
func getNumOfGhost() int {
	numOfGhost := int(math.Ceil(float64(level)/3.0)) + 1
	if numOfGhost > maxNumOfGhost {
		numOfGhost = maxNumOfGhost
	}
	return numOfGhost
}
