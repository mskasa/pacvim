package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"sync"
	"time"

	termbox "github.com/nsf/termbox-go"
	"go.uber.org/multierr"
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

	defaultLevel     = 1
	defaultLife      = 3
	defaultGameSpeed = 3
	upperLimitLife   = 5

	stageDir         = "files/stage/"
	sceneDir         = "files/scene/"
	stageMapMimeType = "text/plain; charset=utf-8"
	maxFileSize      = 512
)

var (
	gameState   = 0
	targetScore = 0
	score       = 0

	gameSpeedMap = map[int]time.Duration{
		1: 1500 * time.Millisecond,
		2: 1250 * time.Millisecond,
		3: 1000 * time.Millisecond,
		4: 750 * time.Millisecond,
		5: 500 * time.Millisecond,
	}

	//go:embed files
	static embed.FS
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	// Get path of text files
	stageMaps, err := dirwalk(stageDir)
	if err != nil {
		return err
	}

	if err = validFiles(stageMaps); err != nil {
		return err
	}

	// Maximum level = Number of stage maps
	maxLevel := len(stageMaps)

	// Read command line arguments
	level := flag.Int("level", defaultLevel, "Level at the start of the game. (1-"+strconv.Itoa(maxLevel)+")")
	life := flag.Int("life", defaultLife, "Remaining lives. (0-"+strconv.Itoa(upperLimitLife)+")")
	gameSpeed := flag.Int("speed", defaultGameSpeed, "Game speed. Bigger is faster. (1-"+strconv.Itoa(len(gameSpeedMap))+")")

	// Validate command line arguments
	if err = validateArgs(level, life, gameSpeed, maxLevel); err != nil {
		return err
	}

	// Initialize termbox
	if err = termbox.Init(); err != nil {
		return err
	}
	if err = termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}
	defer termbox.Close()

	// Display the start screen
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

		gameState = pose
		score = 0
		targetScore = 0
		b.plotStageMap()

		p := new(player)
		p.initPosition(b)

		b.plotScore()
		b.plotSubInfo(*level, *life)

		ghostList, err := createGhosts(*level, b)
		if err != nil {
			return err
		}

		if err = termbox.Flush(); err != nil {
			return err
		}

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
				time.Sleep(gameSpeedMap[*gameSpeed])
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
			if *level == maxLevel {
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

func validateArgs(level *int, life *int, gameSpeed *int, maxLevel int) error {
	flag.Parse()
	if *level > maxLevel || *level < 1 {
		return errors.New("Validation Error: level must be (1-" + strconv.Itoa(maxLevel) + ").")
	}
	if *life > upperLimitLife || *life < 0 {
		return errors.New("Validation Error: life must be (0-" + strconv.Itoa(upperLimitLife) + ").")
	}
	if *gameSpeed > len(gameSpeedMap) || *gameSpeed < 1 {
		return errors.New("Validation Error: speed must be (1-" + strconv.Itoa(len(gameSpeedMap)) + ").")
	}
	return nil
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

func dirwalk(dir string) (map[int]string, error) {
	pathMap := make(map[int]string, 10)
	r1 := regexp.MustCompile(`^map`)
	r2 := regexp.MustCompile(`.txt$`)

	i := 1
	if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && r1.MatchString(info.Name()) && r2.MatchString(info.Name()) {
			pathMap[i] = filepath.Join(dir, info.Name())
			i++
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return pathMap, nil
}

func validFiles(stageMaps map[int]string) error {
	for i := 1; i <= len(stageMaps); i++ {
		if err := validMimeType(stageMaps[i]); err != nil {
			return err
		}
		if err := validFileSize(stageMaps[i]); err != nil {
			return err
		}
	}
	return nil
}
func validMimeType(stageMap string) error {
	bytes, err := os.ReadFile(stageMap)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(bytes)
	if mimeType != stageMapMimeType {
		return errors.New("Invalid mime type: " + mimeType)
	}
	return nil
}
func validFileSize(dir string) (err error) {
	f, err := os.Open(dir)
	if err != nil {
		return err
	}

	defer func() {
		err = multierr.Append(err, f.Close())
	}()

	fi, err := f.Stat()
	if err != nil {
		return err
	}

	if fi.Size() > maxFileSize {
		return errors.New("File size exceeded:" + strconv.Itoa(int(fi.Size())) + " (Max file size is " + strconv.Itoa(maxFileSize) + ")")
	}
	return nil
}
