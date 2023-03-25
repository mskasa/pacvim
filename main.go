package main

import (
	"bytes"
	"embed"
	"fmt"
	"log"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	pose int = iota
	continuing
	quit
	win
	lose

	maxLevel      = 10
	maxNumOfGhost = 4
)

var (
	gameState           = 0
	targetScore         = 0
	score               = 0
	level               = 1
	life                = 3
	inputNum            = 0
	isLowercaseGEntered = false
	gameSpeed           = time.Second
)

//go:embed files
var static embed.FS

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	err := termbox.Init()
	if err != nil {
		return err
	}

	err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		return err
	}

	defer termbox.Close()

	// スタート画面表示
	if err := switchScene("files/scene/start.txt"); err != nil {
		return err
	}

game:
	for {
		if err := stage(); err != nil {
			return err
		}
		switch gameState {
		case win:
			if err := switchScene("files/scene/youwin.txt"); err != nil {
				return err
			}
			level++
			gameSpeed = time.Duration(1000-(level-1)*50) * time.Millisecond
			if level == maxLevel {
				err = switchScene("files/scene/congrats.txt")
				if err != nil {
					return err
				}
				break game
			}
		case lose:
			err = switchScene("files/scene/youlose.txt")
			if err != nil {
				return err
			}
			life--
			if life < 0 {
				break game
			}
		case quit:
			break game
		}
	}
	err = switchScene("files/scene/goodbye.txt")
	if err != nil {
		return err
	}

	return err
}

func stage() error {
	// ステージマップ読み込み
	fileName := "files/stage/map" + strconv.Itoa(level) + ".txt"
	b, w := initScene(fileName)
	err := w.ShowWithLineNum(b)
	if err != nil {
		return err
	}

	// ゲーム情報の初期化
	gameState = pose
	score = 0
	targetScore = 0
	// マップを整形
	b.CheckAllChar()

	// プレイヤー初期化
	p := Initialize(b)

	// ゲーム情報の表示
	b.Displayscore()
	b.DisplayNote()

	// ゴーストを作成
	var gList []*Ghost
	gList, err = b.protGhost()
	if err != nil {
		return err
	}
	// ステージマップを表示
	err = termbox.Flush()
	if err != nil {
		return fmt.Errorf("termbox.Flush(): %s, %v", fileName, err)
	}

	// ゲーム開始待ち状態
	standBy()

	// プレイヤーゴルーチン開始
	ch1 := make(chan bool)
	go p.Control(ch1, b, w)

	// ゴーストゴルーチン開始
	ch2 := make(chan bool)
	go Control(ch2, p, gList)

	// プレイヤーとゴーストの同期を取る
	<-ch1
	<-ch2

	return err
}

// 画面の切り替え
func switchScene(fileName string) error {
	termbox.HideCursor()
	_, w := initScene(fileName)
	err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		return fmt.Errorf("termbox.Clear(): %s, %v", fileName, err)
	}
	for y, l := range w.lines {
		for x, r := range l.Text {
			termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
		}
	}
	err = termbox.Flush()
	if err != nil {
		return fmt.Errorf("termbox.Flush(): %s, %v", fileName, err)
	}
	time.Sleep(1 * time.Second)
	return err
}

// 前処理
func initScene(fileName string) (*Buffer, *Window) {
	f, err := static.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// バッファ初期化
	b := CreateBuffer()
	b.ReadFileToBuf(bytes.NewReader(f))

	// ウィンドウ初期化
	w := CreateWindow(b)

	return b, w
}

// 待機状態
func standBy() {
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEnter {
			gameState = continuing
			break
		}
		if ev.Ch == 'q' {
			gameState = quit
			break
		}
	}
}
