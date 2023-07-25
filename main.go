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

	termbox "github.com/nsf/termbox-go"
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
	chBoundary  = '+'
	chObstacle1 = '-'
	chObstacle2 = '|'
	chObstacle3 = '!'

	sceneStart   = "files/scene/start.txt"
	sceneYouwin  = "files/scene/youwin.txt"
	sceneYoulose = "files/scene/youlose.txt"
	sceneGoodbye = "files/scene/goodbye.txt"
)

//go:embed files
var static embed.FS

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

	level := flag.Int("level", stages[0].level, "Level at the start of the game.")
	life := flag.Int("life", 2, "Remaining lives.")
	flag.Parse()

	stages = splitStages(stages, level)

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

	i := 0
game:
	for i < len(stages) && *life >= 0 {
		p := new(player)
		if err := stages[i].init(p, *life); err != nil {
			return err
		}

		standBy(p)

		if err := stages[i].start(p); err != nil {
			return err
		}

		switch p.state {
		case win:
			if err := switchScene(sceneYouwin); err != nil {
				return err
			}
			i++
		case lose:
			if err := switchScene(sceneYoulose); err != nil {
				return err
			}
			*life--
		case quit:
			break game
		}
	}

	if err := switchScene(sceneGoodbye); err != nil {
		return err
	}

	return nil
}

func splitStages(stages []stage, level *int) []stage {
	for i, stage := range stages {
		if *level == stage.level {
			return stages[i:]
		}
	}
	return stages
}

func standBy(p *player) {
	p.state = pose
	for {
		ev := termbox.PollEvent()
		if ev.Key == termbox.KeyEnter {
			p.state = continuing
			break
		}
		if ev.Ch == 'q' {
			p.state = quit
			break
		}
	}
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
		b, err := static.ReadFile(s.mapPath)
		if err != nil {
			return err
		}
		if err := validateMimeType(b, s.mapPath); err != nil {
			return err
		}
		if err := validateStageMap(b, s.mapPath); err != nil {
			return err
		}
	}
	return nil
}

func validateMimeType(b []byte, filePath string) error {
	mimeType := http.DetectContentType(b)
	if mimeType != stageMapMimeType {
		err := errors.New(filePath + "; Invalid mime type: " + mimeType + ";")
		return fmt.Errorf("%w: %+v", mimeTypeValidationError, err)
	}
	return nil
}

func validateStageMap(b []byte, filePath string) error {
	lines, err := validateStageMapSize(b, filePath)
	if err != nil {
		return err
	}
	if err := validateStageMapShape(filePath, lines); err != nil {
		return err
	}
	if err := validateStageMapBoundary(filePath, lines); err != nil {
		return err
	}
	return nil
}

func validateStageMapSize(b []byte, filePath string) ([]string, error) {
	f := bytes.NewReader(b)
	scanner := bufio.NewScanner(f)
	lines := []string{}
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if len(lines[0]) > maxStageMapWidth {
		err := errors.New(filePath + "; Please keep the stage within 50 columns;")
		return nil, fmt.Errorf("%w: %+v", stageMapValidationError, err)
	}
	if len(lines) > maxStageMapHeight {
		err := errors.New(filePath + "; Please keep the stage within 20 lines;")
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

func validateStageMapBoundary(filePath string, lines []string) error {
	width := len(lines[0])
	height := len(lines)
	lineNo := 1
	errLineNo := []string{}
	for _, s := range lines {
		if lineNo == 1 || lineNo == height {
			if s != strings.Repeat(string(chBoundary), width) {
				errLineNo = append(errLineNo, strconv.Itoa(lineNo))
			}
		} else {
			if !strings.HasPrefix(s, string(chBoundary)) || !strings.HasSuffix(s, string(chBoundary)) {
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
