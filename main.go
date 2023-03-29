package main

import (
	"bytes"
	"embed"
	"flag"
	"log"
	"math"
	"os"
	"path/filepath"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
	"golang.org/x/sync/errgroup"
)

const (
	pose int = iota
	continuing
	quit
	win
	lose

	chGhost  = 'G'
	chTarget = 'o'
	chPoison = 'X'
	chWall1  = '#'
	chWall2  = '|'
	chWall3  = '-'

	sceneDir = "files/scene/"
)

var (
	gameState   = 0
	targetScore = 0           // main, player, buffer
	score       = 0           // main, player, buffer
	gameSpeed   = time.Second // mainのrun, stage

	//go:embed files
	static embed.FS
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	level := flag.Int("level", 1, "level at the start of the game")
	life := flag.Int("life", 3, "remaining lives")
	maxNumOfGhosts := flag.Int("maxNumOfGhosts", 4, "maximum number of ghosts")
	flag.Parse()

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

	// スタート画面表示
	if err := switchScene(sceneDir + "start.txt"); err != nil {
		return err
	}

game:
	for {
		f, err := static.ReadFile(stageMaps[*level])
		if err != nil {
			return err
		}

		b := createBuffer(bytes.NewReader(f))
		w := createWindow(b)
		if err = w.show(b); err != nil {
			return err
		}

		// ゲーム情報の初期化
		gameState = pose
		score = 0
		targetScore = 0
		b.plotStageMap()

		// プレイヤー初期化
		p := new(player)
		p.initPosition(b)

		// ゲーム情報の表示
		b.plotScore()
		b.plotSubInfo(*level, *life)

		ghostList := make([]*ghost, 0, *maxNumOfGhosts)
		ghostPlotRangeList := [][]float64{
			{0.4, 0.4}, // The 1st one:	2nd quadrant, strategyA
			{0.6, 0.6}, // The 2nd one:	4th quadrant, strategyA
			{0.6, 0.4}, // The 3rd one:	1st quadrant, strategyB
			{0.4, 0.6}, // The 4th one:	3rd quadrant, strategyB
		}
		for i := 0; i < numOfGhosts(*level, *maxNumOfGhosts); i++ {
			g := &ghost{
				strategy:  newStrategy(i),
				plotRange: ghostPlotRangeList[i],
			}
			if err = g.initPosition(b); err != nil {
				return err
			}
			ghostList = append(ghostList, g)
		}

		// ステージマップを表示
		if err = termbox.Flush(); err != nil {
			return err
		}

		// ゲーム開始待ち状態
		standBy()

		eg := new(errgroup.Group)

		// Starts a new goroutine that runs for player actions
		eg.Go(func() error {
			return p.action(b)
		})

		// Starts a new goroutine that runs for ghost control
		eg.Go(func() error {
			var wg sync.WaitGroup

			for gameState == continuing {
				wg.Add(len(ghostList))
				// Starts new goroutines that runs for ghosts actions
				for _, g := range ghostList {
					go g.action(&wg, p)
				}
				// Synchronization(waiting for ghosts goroutines to finish)
				wg.Wait()
				if err = termbox.Flush(); err != nil {
					return err
				}
				time.Sleep(gameSpeed)
			}
			return nil
		})

		// Synchronization(waiting for player action goroutine and ghost control goroutine to finish)
		if err := eg.Wait(); err != nil {
			return err
		}

		switch gameState {
		case win:
			if err := switchScene(sceneDir + "youwin.txt"); err != nil {
				return err
			}
			*level++
			gameSpeed = time.Duration(1000-(*level-1)*50) * time.Millisecond
			if *level == len(stageMaps) {
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
			*life--
			if *life < 0 {
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
	time.Sleep(1 * time.Second)
	return err
}

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

func numOfGhosts(level int, maxNumOfGhosts int) int {
	numOfGhost := int(math.Ceil(float64(level)/3.0)) + 1
	if numOfGhost > maxNumOfGhosts {
		numOfGhost = maxNumOfGhosts
	}
	return numOfGhost
}
