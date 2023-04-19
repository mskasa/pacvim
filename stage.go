package main

import (
	"errors"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type stage struct {
	level         int
	mapPath       string
	hunterBuilder iEnemyBuilder
	ghostBuilder  iEnemyBuilder
	enemies       []iEnemy
	gameSpeed     time.Duration
	width         int
	height        int
	firstLine     int
	endLine       int
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
			gameSpeed:     1250 * time.Millisecond,
		},
		{
			level:         2,
			mapPath:       "files/stage/map02.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     1000 * time.Millisecond,
		},
		{
			level:         3,
			mapPath:       "files/stage/map03.txt",
			hunterBuilder: defaultHunterBuilder.strategize(&tricky{}),
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     750 * time.Millisecond,
		},
		{
			level:         4,
			mapPath:       "files/stage/map04.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     1000 * time.Millisecond,
		},
		{
			level:         5,
			mapPath:       "files/stage/map05.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     1000 * time.Millisecond,
		},
		{
			level:         6,
			mapPath:       "files/stage/map06.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     750 * time.Millisecond,
		},
		{
			level:         7,
			mapPath:       "files/stage/map07.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder.speed(1).strategize(&tricky{}),
			gameSpeed:     750 * time.Millisecond,
		},
		{
			level:         8,
			mapPath:       "files/stage/map08.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     750 * time.Millisecond,
		},
		{
			level:         9,
			mapPath:       "files/stage/map09.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     500 * time.Millisecond,
		},
		{
			level:         10,
			mapPath:       "files/stage/map10.txt",
			hunterBuilder: defaultHunterBuilder,
			ghostBuilder:  defaultGhostBuilder,
			gameSpeed:     500 * time.Millisecond,
		},
	}
}

func getStage(stages []stage, level int) (stage, error) {
	for _, stage := range stages {
		if level == stage.level {
			return stage, nil
		}
	}
	return stage{}, errors.New("File does not exist: " + stages[level].mapPath)
}

func (s *stage) plot(b *buffer, p *player) {
	s.firstLine = 0
	s.width = len(b.lines[0].text) + b.offset
	s.height = len(b.lines)
	for y := 0; y < s.height; y++ {
		for x := b.offset; x < s.width; x++ {
			if isCharTarget(x, y) || isCharSpace(x, y) {
				if s.firstLine == 0 {
					s.firstLine = y
				}
				s.endLine = y
				if isCharTarget(x, y) {
					p.targetScore++
				}
			} else if isCharPlayer(x, y) {
				p.x = x
				p.y = y
				termbox.SetCell(p.x, p.y, chSpace, termbox.ColorWhite, termbox.ColorBlack)
				termbox.SetCursor(p.x, p.y)
			} else if isCharBorder(x, y) {
				termbox.SetCell(x, y, chBorder, termbox.ColorYellow, termbox.ColorBlack)
			} else if isCharObstacle(x, y) {
				var r rune
				if isCharObstacle(x-1, y) || isCharObstacle(x+1, y) {
					r = chObstacle1
				} else if isCharObstacle(x, y-1) || isCharObstacle(x, y+1) {
					r = chObstacle2
				}
				termbox.SetCell(x, y, r, termbox.ColorYellow, termbox.ColorBlack)
			} else if isCharPoison(x, y) {
				termbox.SetCell(x, y, chPoison, termbox.ColorMagenta, termbox.ColorBlack)
			} else if isCharHunter(x, y) {
				h := s.hunterBuilder.build()
				h.setPosition(x, y)
				char, color := h.getDisplayFormat()
				termbox.SetCell(x, y, char, color, color)
				s.enemies = append(s.enemies, h)
			} else if isCharGhost(x, y) {
				g := s.ghostBuilder.build()
				g.setPosition(x, y)
				char, color := g.getDisplayFormat()
				termbox.SetCell(x, y, char, color, color)
				s.enemies = append(s.enemies, g)
			}
		}
	}
}

func (s *stage) plotSubInfo(life int) {
	textMap := map[int]string{
		0: "Level: " + strconv.Itoa(s.level),
		1: "Life : " + strconv.Itoa(life),
		2: "PRESS ENTER TO PLAY!",
		3: "q TO EXIT!"}
	position := s.height + 1
	for i := 0; i < len(textMap); i++ {
		for x, r := range []rune(textMap[i]) {
			termbox.SetCell(x, position, r, termbox.ColorWhite, termbox.ColorBlack)
		}
		position++
	}
}
