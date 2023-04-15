package main

import (
	"bytes"
	"embed"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"strconv"
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

	chPlayer    = 'P'
	chHunter    = 'H'
	chGhost     = 'G'
	chTarget    = 'o'
	chPoison    = 'X'
	chBorder    = '+'
	chObstacle1 = '-'
	chObstacle2 = '|'
	chObstacle3 = '!'

	defaultLevel     = 1
	defaultLife      = 3
	defaultGameSpeed = 3
	upperLimitLife   = 5

	sceneStart    = "files/scene/start.txt"
	sceneYouwin   = "files/scene/youwin.txt"
	sceneYoulose  = "files/scene/youlose.txt"
	sceneCongrats = "files/scene/congrats.txt"
	sceneGoodbye  = "files/scene/goodbye.txt"

	stageDir         = "files/stage/"
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

type stage struct {
	level         int
	mapPath       string
	hunterBuilder iEnemyBuilder
	ghostBuilder  iEnemyBuilder
	enemies       []iEnemy
}

func initStages() []stage {
	defaultHunterBuilder := newEnemyBuilder().defaultHunter()
	defaultGhostBuilder := newEnemyBuilder().defaultGhost()
	return []stage{
		{
			level:         1,
			mapPath:       "files/stage/map01.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
		},
		{
			level:         2,
			mapPath:       "files/stage/map02.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
		},
	}
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	stages := initStages()
	if err := validateFiles(stages); err != nil {
		return err
	}
	maxLevel := len(stages)

	// Read command line arguments
	level := flag.Int("level", defaultLevel, "Level at the start of the game. (1-"+strconv.Itoa(maxLevel)+")")
	life := flag.Int("life", defaultLife, "Remaining lives. (0-"+strconv.Itoa(upperLimitLife)+")")
	gameSpeed := flag.Int("speed", defaultGameSpeed, "Game speed. Bigger is faster. (1-"+strconv.Itoa(len(gameSpeedMap))+")")

	// Validate command line arguments
	if err := validateArgs(level, life, gameSpeed, maxLevel); err != nil {
		return err
	}

	// Initialize termbox
	if err := termbox.Init(); err != nil {
		return err
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}
	defer termbox.Close()

	// Display the start screen
	if err := switchScene(sceneStart); err != nil {
		return err
	}

game:
	for {
		stage, err := getStage(stages, *level)
		if err != nil {
			return err
		}
		f, err := static.ReadFile(stage.mapPath)
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

		p := new(player)

		stage.plot(b, p)
		b.plotScore()
		b.plotSubInfo(*level, *life)

		if err = termbox.Flush(); err != nil {
			return err
		}

		// stand-by
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

		eg := new(errgroup.Group)

		// Starts a new goroutine that runs for player actions
		eg.Go(func() error {
			return p.action(b)
		})

		// Starts a new goroutine that runs for ghost control
		eg.Go(func() error {
			for gameState == continuing {
				for _, e := range stage.enemies {
					e.move(e.think(p))
					e.hasCaptured(p)
				}
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
			if err := switchScene(sceneYouwin); err != nil {
				return err
			}
			*level++
			if *level > maxLevel {
				if err = switchScene(sceneCongrats); err != nil {
					return err
				}
				break game
			}
		case lose:
			if err = switchScene(sceneYoulose); err != nil {
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
	if err := switchScene(sceneGoodbye); err != nil {
		return err
	}

	return nil
}

func getStage(stages []stage, level int) (stage, error) {
	for _, stage := range stages {
		if level == stage.level {
			return stage, nil
		}
	}
	return stage{}, errors.New("File does not exist: " + stages[level].mapPath)
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
	time.Sleep(750 * time.Millisecond)
	return err
}

func validateFiles(stages []stage) error {
	for i := 0; i < len(stages); i++ {
		if err := validateMimeType(stages[i].mapPath); err != nil {
			return err
		}
		if err := validateFileSize(stages[i].mapPath); err != nil {
			return err
		}
	}
	return nil
}
func validateMimeType(stageMap string) error {
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
func validateFileSize(dir string) (err error) {
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
