package main

import (
	"bytes"
	"embed"
	"flag"
	"log"
	"time"

	"pacvim/buffer"
	"pacvim/game"
	"pacvim/ghost"
	"pacvim/player"
	"pacvim/window"

	termbox "github.com/nsf/termbox-go"
)

//go:embed files
var static embed.FS
var level = flag.Int("level", 1, "e.g. -level 2")

func main() {
	var err error
	// termbox初期化
	initTermbox()
	defer termbox.Close()
	// スタート画面表示
	switchScene("files/scene/start.txt")
	// コマンドライン引数でレベルを設定
	flag.Parse()
	game.SetLevel(*level)

	for {
		err = stage()
		if game.HasStageWon() {
			// ステージクリア
			switchScene("files/scene/youwin.txt")
		} else if game.HasStageLost() {
			// ステージ失敗
			switchScene("files/scene/youlose.txt")
		}
		if game.IsFinished() || err != nil {
			// 終了判定およびエラー判定
			break
		}
	}
	if err != nil {
		// エラー発生の場合
		log.Fatal(err)
	} else {
		if game.HasGameCleared() {
			// ゲームクリア
			switchScene("files/scene/congrats.txt")
		}
		switchScene("files/scene/goodbye.txt")
	}
}

func stage() error {
	// ステージマップ読み込み
	fileName := "files/stage/map" + game.GetLevel() + ".txt"
	b, w := new(fileName)
	err := w.ShowWithLineNum(b)
	if err != nil {
		return err
	}

	// ゲーム情報の初期化
	game.Reset()
	// マップを整形
	b.CheckAllChar()

	// プレイヤー初期化
	p := player.Initialize(b)

	// ゲーム情報の表示
	b.DisplayScore()
	b.DisplayNote()

	// ゴーストを作成
	var gList []*ghost.Ghost
	gList, err = ghost.Spawn(b)
	if err != nil {
		return err
	}
	// ステージマップを表示
	termbox.Flush()

	// ゲーム開始待ち状態
	standBy()

	// プレイヤーゴルーチン開始
	ch1 := make(chan bool)
	go p.Control(ch1, b, w)

	// ゴーストゴルーチン開始
	ch2 := make(chan bool)
	go ghost.Control(ch2, p, gList)

	// プレイヤーとゴーストの同期を取る
	<-ch1
	<-ch2

	return err
}

// 画面の切り替え
func switchScene(fileName string) {
	termbox.HideCursor()
	b, w := new(fileName)
	err := w.Show(b)
	if err != nil {
		log.Fatal(err)
	}
	termbox.Flush()
	time.Sleep(1 * time.Second)
}

// 前処理
func new(fileName string) (*buffer.Buffer, *window.Window) {
	f, err := static.ReadFile(fileName)
	if err != nil {
		panic(err)
	}

	// バッファ初期化
	b := buffer.New()
	b.ReadFileToBuf(bytes.NewReader(f))

	// ウィンドウ初期化
	w := window.New(b)

	return b, w
}
func initTermbox() {
	err := termbox.Init()
	if err != nil {
		log.Fatal(err)
	}
	err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack)
	if err != nil {
		log.Fatal(err)
	}
}

// 待機状態
func standBy() {
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEnter {
			game.Start()
			break
		}
		if ev.Ch == 'q' {
			game.Quit()
			break
		}
	}
}
