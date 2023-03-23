package main

import (
	"bufio"
	"io"

	termbox "github.com/nsf/termbox-go"
)

type Buffer struct {
	Lines  []*Line
	Offset int
}

type Line struct {
	Text []rune
}

const (
	ChGhost  = 'G'
	ChTarget = 'o'
	ChPoison = 'X'
	ChWall1  = '#'
	ChWall2  = '|'
	ChWall3  = '-'
)

var (
	FirstTargetY int
	LastTargetY  int
)

func CreateBuffer() *Buffer {
	b := new(Buffer)
	return b
}

// 対象の行の文字列を取得
func (b *Buffer) GetTextOnLine(y int) []rune {
	return b.Lines[y].Text
}

// 行数の取得
func (b *Buffer) NumOfLines() int {
	return len(b.Lines)
}

// ファイルをバッファに読み込む
func (b *Buffer) ReadFileToBuf(reader io.Reader) {
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		l := new(Line)
		l.Text = []rune(scanner.Text())
		b.Lines = append(b.Lines, l)
	}
}

/*
 * マップ読み込み時の処理
 * 1. viコマンド処理に必要な値を保持（FirstTargetY、LastTargetY）
 * 2. 目標スコアの設定：game.AddTargetScore()
 * 3. 壁の文字変換：convertWallChar
 * 4. 毒に色を付与：colorThePoison
 */
func (b *Buffer) CheckAllChar() {
	firstFlg := true
	winWidth, _ := termbox.Size()
	textHeight := b.NumOfLines()
	for y := 0; y < textHeight; y++ {
		for x := b.Offset; x < winWidth; x++ {
			if IsTarget(x, y) {
				if firstFlg {
					FirstTargetY = y
					firstFlg = false
				}
				LastTargetY = y
				AddTargetScore()
			}
			b.convertWallChar(x, y)
			colorThePoison(x, y)
		}
	}
}

// 壁を変換（# → | or -）
func (b *Buffer) convertWallChar(x, y int) {
	if IsWall(x, y) {
		r := '-'
		if y == 0 || y == b.NumOfLines() {
			r = '-'
		} else if IsWall(x, y-1) && IsWall(x, y+1) {
			r = '|'
		} else if IsWall(x-1, y) || IsWall(x+1, y) {
			r = '-'
		} else if IsWall(x, y-1) || IsWall(x, y+1) {
			r = '|'
		}
		termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
	}
}

// 毒に色を付与
func colorThePoison(x, y int) {
	if IsPoison(x, y) {
		termbox.SetCell(x, y, ChPoison, termbox.ColorMagenta, termbox.ColorBlack)
	}
}

// スコアの表示
func (b *Buffer) DisplayScore() {
	textHeight := b.NumOfLines()
	score := []rune("Score:" + GetScore() + "/" + GetTargetScore())
	for x, r := range score {
		termbox.SetCell(x, textHeight, r, termbox.ColorRed, termbox.ColorBlack)
	}
}

// 補足情報の表示
func (b *Buffer) DisplayNote() {
	textMap := map[int]string{
		0: "Level:" + GetLevel(),
		1: "Life:" + GetLife(),
		2: "PRESS ENTER TO PLAY!",
		3: "q TO EXIT!"}
	textHeight := b.NumOfLines() + 1
	for i := 0; i < len(textMap); i++ {
		for x, r := range []rune(textMap[i]) {
			termbox.SetCell(x, textHeight, r, termbox.ColorWhite, termbox.ColorBlack)
		}
		textHeight++
	}
}

// 文字判定
func IsWall(x, y int) bool {
	return IsChar(x, y, ChWall1) || IsChar(x, y, ChWall2) || IsChar(x, y, ChWall3)
}
func IsGhost(x, y int) bool {
	return IsChar(x, y, ChGhost)
}
func IsSpace(x, y int) bool {
	return IsChar(x, y, ' ')
}
func IsTarget(x, y int) bool {
	return IsChar(x, y, ChTarget)
}
func IsPoison(x, y int) bool {
	return IsChar(x, y, ChPoison)
}
func IsChar(x, y int, r rune) bool {
	winWidth, _ := termbox.Size()
	cell := termbox.CellBuffer()[(winWidth*y)+x]
	return r == cell.Ch
}
func IsColorWhite(x, y int) bool {
	winWidth, _ := termbox.Size()
	cell := termbox.CellBuffer()[(winWidth*y)+x]
	return cell.Fg == termbox.ColorWhite
}
