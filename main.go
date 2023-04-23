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

	sceneStart    = "files/scene/start.txt"
	sceneYouwin   = "files/scene/youwin.txt"
	sceneYoulose  = "files/scene/youlose.txt"
	sceneCongrats = "files/scene/congrats.txt"
	sceneGoodbye  = "files/scene/goodbye.txt"
)

var (
	gameState = 0

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

	level := flag.Int("level", 1, "Level at the start of the game. (1-"+strconv.Itoa(maxLevel)+")")
	life := flag.Int("life", 3, "Remaining lives. (0-5)")
	if err := validateArgs(level, life, maxLevel); err != nil {
		return err
	}

	if err := termbox.Init(); err != nil {
		return err
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		return err
	}
	defer termbox.Close()

	if err := switchScene(sceneStart); err != nil {
		return err
	}

game:
	for {
		stage, err := getStage(stages, *level)
		if err != nil {
			return err
		}
		p := new(player)
		if err = stage.init(p, *life); err != nil {
			return err
		}

		standBy()

		eg := new(errgroup.Group)

		eg.Go(func() error {
			if err := p.action(stage); err != nil {
				return err
			}
			return nil
		})

		eg.Go(func() error {
			for gameState == continuing {
				for _, e := range stage.enemies {
					e.move(e.think(p))
					e.hasCaptured(p)
				}
				if err = termbox.Flush(); err != nil {
					return err
				}
				time.Sleep(stage.gameSpeed)
			}
			return nil
		})

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

func standBy() {
	gameState = pose
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

func validateArgs(level *int, life *int, maxLevel int) error {
	flag.Parse()
	if *level > maxLevel || *level < 1 {
		return errors.New("Validation Error: level must be (1-" + strconv.Itoa(maxLevel) + ").")
	}
	if *life > 5 || *life < 0 {
		return errors.New("Validation Error: life must be (0-5).")
	}
	return nil
}

var (
	mimeTypeValidationError = errors.New("MIME Type Validation Error")
	fileSizeValidationError = errors.New("File Size Validation Error")
	stageMapValidationError = errors.New("Stage Map Validation Error")
)

const (
	stageMapMimeType  = "text/plain; charset=utf-8"
	maxFileSize       = 1024
	maxStageMapWidth  = 50
	maxStageMapHeight = 20
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
	b, err := static.ReadFile(filePath)
	if err != nil {
		return err
	}
	mimeType := http.DetectContentType(b)
	if mimeType != stageMapMimeType {
		err = errors.New(filePath + "; Invalid mime type: " + mimeType + ";")
		return fmt.Errorf("%w: %+v", mimeTypeValidationError, err)
	}
	return nil
}

func validateFileSize(filePath string) (err error) {
	f, err := static.Open(filePath)
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
		err = errors.New(filePath + "; File size exceeded: " + strconv.Itoa(int(fi.Size())) + " (Max file size is " + strconv.Itoa(maxFileSize) + ");")
		return fmt.Errorf("%w: %+v", fileSizeValidationError, err)
	}
	return nil
}

func validateStageMap(filePath string) error {
	lines, err := validateStageMapSize(filePath)
	if err != nil {
		return err
	}
	if err := validateStageMapShape(filePath, lines); err != nil {
		return err
	}
	if err := validateStageMapBorder(filePath, lines); err != nil {
		return err
	}
	return nil
}

func validateStageMapSize(filePath string) ([]string, error) {
	b, err := static.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	f := bytes.NewReader(b)
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines[0]) > maxStageMapWidth {
		err = errors.New(filePath + "; Please keep the stage within 50 columns;")
		return nil, fmt.Errorf("%w: %+v", stageMapValidationError, err)
	}
	if len(lines) > maxStageMapHeight {
		err = errors.New(filePath + "; Please keep the stage within 20 lines;")
		return nil, fmt.Errorf("%w: %+v", stageMapValidationError, err)
	}
	return lines, nil
}

func validateStageMapShape(filePath string, lines []string) error {
	width := len(lines[0])
	lineNo := 1
	errLineNo := []string{}
	for _, line := range lines {
		if len(line) != width {
			errLineNo = append(errLineNo, strconv.Itoa(lineNo))
		}
		lineNo++
	}
	if len(errLineNo) > 0 {
		err := errors.New(filePath + "; Make the width of the stage map uniform (line " + strings.Join(errLineNo, ",") + ");")
		return fmt.Errorf("%w: %+v", stageMapValidationError, err)
	}
	return nil
}

func validateStageMapBorder(filePath string, lines []string) error {
	width := len(lines[0])
	height := len(lines)
	lineNo := 1
	errLineNo := []string{}
	for _, s := range lines {
		if lineNo == 1 || lineNo == height {
			if s != strings.Repeat(string(chBorder), width) {
				errLineNo = append(errLineNo, strconv.Itoa(lineNo))
			}
		} else {
			if !strings.HasPrefix(s, string(chBorder)) || !strings.HasSuffix(s, string(chBorder)) {
				errLineNo = append(errLineNo, strconv.Itoa(lineNo))
			}
		}
		lineNo++
	}
	if len(errLineNo) > 0 {
		err := errors.New(filePath + "; Create a boundary for the stage map with '+' (line " + strings.Join(errLineNo, ",") + ");")
		return fmt.Errorf("%w: %+v", stageMapValidationError, err)
	}
	return nil
}
