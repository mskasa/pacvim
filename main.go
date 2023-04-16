package main

import (
	"bufio"
	"bytes"
	"embed"
	"errors"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
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
	chSpace     = ' '
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

var (
	mimeTypeValidationError = errors.New("MIME Type Validation Error")
	fileSizeValidationError = errors.New("File Size Validation Error")
	stageMapValidationError = errors.New("Stage Map Validation Error")
)

const (
	stageMapMimeType = "text/plain; charset=utf-8"
	maxFileSize      = 1024
)

func validateFiles(stages []stage) error {
	for _, s := range stages {
		if err := validateMimeType(s.mapPath); err != nil {
			return err
		}
		if err := validateFileSize(s.mapPath); err != nil {
			return err
		}
		if err := validateStageMap(s.mapPath); err != nil {
			return err
		}
	}
	return nil
}

func validateMimeType(filePath string) error {
	bytes, err := os.ReadFile(filePath)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(bytes)
	if mimeType != stageMapMimeType {
		err = errors.New(filePath + "; Invalid mime type: " + mimeType + ";")
		return fmt.Errorf("%w: %+v", mimeTypeValidationError, err)
	}
	return nil
}

func validateFileSize(filePath string) (err error) {
	f, err := os.Open(filePath)
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
		err = errors.New(filePath + ": File size exceeded:" + strconv.Itoa(int(fi.Size())) + " (Max file size is " + strconv.Itoa(maxFileSize) + ")")
		return fmt.Errorf("%+v: %w", err, fileSizeValidationError)
	}
	return nil
}

func validateStageMap(filePath string) error {
	var err error
	f, _ := os.Open(filePath)
	defer func() {
		err = multierr.Append(err, f.Close())
	}()
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	height := len(lines)
	if height < 10 || height > 20 {
		err = errors.New(filePath + ": Make the stage map 10 to 20 lines")
		return fmt.Errorf("%+v: %w", err, stageMapValidationError)
	}
	width := len(lines[0])
	if width < 20 || width > 50 {
		err = errors.New(filePath + ": Make the stage map 20 to 50 columns")
		return fmt.Errorf("%+v: %w", err, stageMapValidationError)
	}
	i := 1
	var errNoBorder1 []string
	var errNoBorder2 []string
	countPlayer := 0
	countEnemy := 0
	countTarget := 0
	for _, s := range lines {
		if len(s) != width {
			errNoBorder1 = append(errNoBorder1, strconv.Itoa(i))
		}
		if i == 1 || i == height {
			if s != strings.Repeat(string(chBorder), width) {
				errNoBorder2 = append(errNoBorder2, strconv.Itoa(i))
			}
		} else {
			if !strings.HasPrefix(s, string(chBorder)) || !strings.HasSuffix(s, string(chBorder)) {
				errNoBorder2 = append(errNoBorder2, strconv.Itoa(i))
			}
			countPlayer += strings.Count(s, string(chPlayer))
			if countPlayer > 1 {
				err = multierr.Append(err, errors.New(filePath+": Please set only one P"))
				return fmt.Errorf("%+v: %w", err, stageMapValidationError)
			}
			countEnemy += strings.Count(s, string(chHunter))
			countEnemy += strings.Count(s, string(chGhost))
			countTarget += strings.Count(s, string(chTarget))
		}
		i++
	}
	if len(errNoBorder1) > 0 {
		err = multierr.Append(err, errors.New(filePath+": Make the width of the stage map uniform (line "+strings.Join(errNoBorder1, ",")+")"))
	}
	if len(errNoBorder2) > 0 {
		err = multierr.Append(err, errors.New(filePath+": Create a boundary for the stage map with '+' (line "+strings.Join(errNoBorder2, ",")+")"))
	}
	if countPlayer == 0 {
		err = multierr.Append(err, errors.New(filePath+": Please set one P"))
	}
	if countEnemy == 0 {
		err = multierr.Append(err, errors.New(filePath+": Please set one or more enemies"))
	}
	if countTarget == 0 {
		err = multierr.Append(err, errors.New(filePath+": Please set one or more targets"))
	}
	if err != nil {
		return fmt.Errorf("%+v: %w", err, stageMapValidationError)
	}
	return nil
}
