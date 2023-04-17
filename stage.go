package main

import (
	"errors"
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
	b.firstLine = 0
	rightEnd := len(b.lines[0].text) + b.offset
	lowerEnd := len(b.lines)
	for y := 0; y < lowerEnd; y++ {
		for x := b.offset; x < rightEnd; x++ {
			if isCharTarget(x, y) || isCharSpace(x, y) {
				if b.firstLine == 0 {
					b.firstLine = y
				}
				b.endLine = y
				if isCharTarget(x, y) {
					targetScore++
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
