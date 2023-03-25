package main

import (
	"bytes"
	"embed"
	"log"
	"os"
	"path/filepath"
	"time"

	termbox "github.com/nsf/termbox-go"
)

const (
	pose int = iota
	continuing
	quit
	win
	lose

	maxNumOfGhost = 4

	chGhost  = 'G'
	chTarget = 'o'
	chPoison = 'X'
	chWall1  = '#'
	chWall2  = '|'
	chWall3  = '-'

	sceneDir = "files/scene/"
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

	firstTargetY int
	lastTargetY  int
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

	if err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}

	defer termbox.Close()

	stageMaps, err := dirwalk("./files/stage")
	if err != nil {
		return err
	}
	maxLevel := len(stageMaps)

	// スタート画面表示
	if err := switchScene(sceneDir + "start.txt"); err != nil {
		return err
	}

game:
	for {
		if err := stage(stageMaps[level]); err != nil {
			return err
		}
		switch gameState {
		case win:
			if err := switchScene(sceneDir + "youwin.txt"); err != nil {
				return err
			}
			level++
			gameSpeed = time.Duration(1000-(level-1)*50) * time.Millisecond
			if level == maxLevel {
				err = switchScene(sceneDir + "congrats.txt")
				if err != nil {
					return err
				}
				break game
			}
		case lose:
			err = switchScene(sceneDir + "youlose.txt")
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
	err = switchScene(sceneDir + "goodbye.txt")
	if err != nil {
		return err
	}

	return err
}

func stage(stageMap string) error {
	b, w, err := initScene(stageMap)
	if err != nil {
		return err
	}
	if err = w.ShowWithLineNum(b); err != nil {
		return err
	}

	// ゲーム情報の初期化
	gameState = pose
	score = 0
	targetScore = 0
	// マップを整形
	b.checkAllChar()

	// プレイヤー初期化
	p := Initialize(b)

	// ゲーム情報の表示
	b.displayscore()
	b.displayNote()

	// ゴーストを作成
	var gList []*Ghost
	gList, err = b.protGhost()
	if err != nil {
		return err
	}
	// ステージマップを表示
	if err = termbox.Flush(); err != nil {
		return err
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
	_, w, err := initScene(fileName)
	if err != nil {
		return err
	}
	if err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}
	for y, l := range w.lines {
		for x, r := range l.Text {
			termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
		}
	}
	if err = termbox.Flush(); err != nil {
		return err
	}
	time.Sleep(1 * time.Second)
	return err
}

// 前処理
func initScene(fileName string) (*Buffer, *Window, error) {
	f, err := static.ReadFile(fileName)
	if err != nil {
		return nil, nil, err
	}

	// バッファ初期化
	b := createBuffer()
	b.readFileToBuf(bytes.NewReader(f))

	// ウィンドウ初期化
	w := CreateWindow(b)

	return b, w, err
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

func dirwalk(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var paths []string
	for _, file := range files {
		paths = append(paths, filepath.Join(dir, file.Name()))
	}

	return paths, err
}
