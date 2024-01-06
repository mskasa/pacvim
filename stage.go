package main

import (
	"bytes"
	"strconv"
	"time"

	termbox "github.com/nsf/termbox-go"
	"golang.org/x/sync/errgroup"
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
}

func initStages() []stage {
	return []stage{
		{
			level:         1,
			mapPath:       "files/stage/map01.txt",
			hunterBuilder: newEnemyBuilder().defaultHunter(),
			gameSpeed:     1250 * time.Millisecond,
		},
		{
			level:         2,
			mapPath:       "files/stage/map02.txt",
			hunterBuilder: newEnemyBuilder().defaultHunter().strategize(&tricky{}),
			ghostBuilder:  newEnemyBuilder().defaultGhost(),
			gameSpeed:     1000 * time.Millisecond,
		},
		{
			level:         3,
			mapPath:       "files/stage/map03.txt",
			hunterBuilder: newEnemyBuilder().defaultHunter(),
			ghostBuilder:  newEnemyBuilder().defaultGhost(),
			gameSpeed:     1000 * time.Millisecond,
		},
		{
			level:         4,
			mapPath:       "files/stage/map04.txt",
			hunterBuilder: newEnemyBuilder().defaultHunter(),
			gameSpeed:     750 * time.Millisecond,
		},
		{
			level:         5,
			mapPath:       "files/stage/map05.txt",
			hunterBuilder: newEnemyBuilder().defaultHunter().strategize(&tricky{}),
			gameSpeed:     750 * time.Millisecond,
		},
	}
}

func (s *stage) init(p *player, life int) error {
	f, err := static.ReadFile(s.mapPath)
	if err != nil {
		return err
	}

	b := createBuffer(bytes.NewReader(f))
	w := createWindow(b)
	if err = w.show(b); err != nil {
		return err
	}

	s.plot(b, p)
	p.plotScore(*s)
	s.plotSubInfo(life)

	if err = termbox.Flush(); err != nil {
		return err
	}
	return nil
}

func (s *stage) plot(b *buffer, p *player) {
	s.enemies = nil
	s.width = len(b.lines[0].text) + b.offset
	s.height = len(b.lines)
	for y := 0; y < s.height; y++ {
		for x := b.offset; x < s.width; x++ {
			if isCharApple(x, y) {
				p.targetScore++
			} else if isCharPlayer(x, y) {
				p.x, p.y = x, y
				termbox.SetCell(p.x, p.y, chSpace, termbox.ColorWhite, termbox.ColorBlack)
				termbox.SetCursor(p.x, p.y)
			} else if isCharBoundary(x, y) {
				termbox.SetCell(x, y, chBoundary, termbox.ColorYellow, termbox.ColorBlack)
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

func (s stage) plotSubInfo(life int) {
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

func (s stage) start(p *player) error {
	eg := new(errgroup.Group)

	eg.Go(func() error {
		for p.state == continuing {
			if err := p.control(s); err != nil {
				return err
			}
		}
		return nil
	})

	eg.Go(func() error {
		for p.state == continuing {
			if err := s.control(p); err != nil {
				return err
			}
		}
		return nil
	})

	if err := eg.Wait(); err != nil {
		return err
	}
	return nil
}

func (s stage) control(p *player) error {
	// Implemented as sequential execution for the following reasons:
	// - The processing content is light.
	// - Considering the overlap of enemies makes the implementation complex.
	for _, e := range s.enemies {
		e.move(e.think(p))
		e.hasCaptured(p)
	}
	if err := termbox.Flush(); err != nil {
		return err
	}
	time.Sleep(s.gameSpeed)
	return nil
}
