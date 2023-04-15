package main

import (
	"errors"

	termbox "github.com/nsf/termbox-go"
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

func getStage(stages []stage, level int) (stage, error) {
	for _, stage := range stages {
		if level == stage.level {
			return stage, nil
		}
	}
	return stage{}, errors.New("File does not exist: " + stages[level].mapPath)
}

// Convert characters
// Color characters
// Save firstTargetY, lastTargetY and targetScore
func (s *stage) plot(b *buffer, p *player) {
	firstFlg := true
	width, height := termbox.Size()
	for y := 0; y < height; y++ {
		for x := b.offset; x < width; x++ {
			if isCharTarget(x, y) {
				if firstFlg {
					b.firstTargetY = y
					firstFlg = false
				}
				b.lastTargetY = y
				targetScore++
			} else if isCharPlayer(x, y) {
				p.x = x
				p.y = y
				termbox.SetCell(p.x, p.y, ' ', termbox.ColorWhite, termbox.ColorBlack)
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
				termbox.SetCell(x, y, char, color, termbox.ColorBlack)
				s.enemies = append(s.enemies, h)
			} else if isCharGhost(x, y) {
				g := s.ghostBuilder.build()
				g.setPosition(x, y)
				char, color := g.getDisplayFormat()
				termbox.SetCell(x, y, char, color, termbox.ColorBlack)
				s.enemies = append(s.enemies, g)
			}
		}
	}
}
