package main

import (
	"bytes"
	"testing"

	termbox "github.com/nsf/termbox-go"
)

func TestEnemyControl(t *testing.T) {
	cases := map[string]struct {
		mapPath       string
		expectedMoves int
		enemyBuilder  iEnemyBuilder
	}{
		"hunter_with_obstacle": {
			mapPath:       "files/test/enemy/control/hunter_with_obstacle.txt",
			expectedMoves: 5,
			enemyBuilder:  newEnemyBuilder().defaultHunter(),
		},
		"ghost_with_obstacle": {
			mapPath:       "files/test/enemy/control/ghost_with_obstacle.txt",
			expectedMoves: 6,
			enemyBuilder:  newEnemyBuilder().defaultGhost(),
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p, stage, err := enemyActionTestInit(t, tt.mapPath, tt.enemyBuilder)
			if err != nil {
				t.Error(err)
			}
			p.state = continuing
			count := 0
			for p.state == continuing {
				if err := control(stage, p); err != nil {
					t.Error(err)
				}
				count++
			}
			if count != tt.expectedMoves {
				t.Errorf("expected %d but %d", tt.expectedMoves, count)
			}
		})
	}
}

func TestEnemyMove(t *testing.T) {
	cases := map[string]struct {
		mapPath       string
		expectedMoves int
		enemyBuilder  iEnemyBuilder
	}{
		"hunter_with_other_enemies": {
			mapPath:       "files/test/enemy/move/hunter_with_other_enemies.txt",
			expectedMoves: 5,
			enemyBuilder:  newEnemyBuilder().defaultHunter(),
		},
		"ghost_with_other_enemies": {
			mapPath:       "files/test/enemy/move/ghost_with_other_enemies.txt",
			expectedMoves: 10,
			enemyBuilder:  newEnemyBuilder().defaultGhost(),
		},
		"tricky": {
			mapPath:       "files/test/enemy/move/tricky.txt",
			expectedMoves: 5,
			enemyBuilder:  newEnemyBuilder().defaultHunter().strategize(&tricky{}),
		},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p, stage, err := enemyActionTestInit(t, tt.mapPath, tt.enemyBuilder)
			if err != nil {
				t.Error(err)
			}
			e := stage.enemies[0]
			p.state = continuing
			count := 0
			for p.state == continuing {
				e.move(e.think(p))
				e.hasCaptured(p)
				count++
			}
			if count != tt.expectedMoves && name != "tricky" {
				t.Errorf("expected %d but %d", tt.expectedMoves, count)
			}
		})
	}
}

func TestEnemyThink(t *testing.T) {
	cases := map[string]struct {
		playerX   int
		playerY   int
		expectedX int
		expectedY int
	}{
		"up": {
			playerX:   5,
			playerY:   1,
			expectedX: 5,
			expectedY: 2,
		},
		"down": {
			playerX:   5,
			playerY:   5,
			expectedX: 5,
			expectedY: 4,
		},
		"left": {
			playerX:   3,
			playerY:   3,
			expectedX: 4,
			expectedY: 3,
		},
		"right": {
			playerX:   7,
			playerY:   3,
			expectedX: 6,
			expectedY: 3,
		},
	}
	p, stage, err := enemyActionTestInit(t, "files/test/enemy/think/think.txt", newEnemyBuilder().defaultHunter())
	if err != nil {
		t.Error()
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p.x, p.y = tt.playerX, tt.playerY
			x, y := stage.enemies[0].think(p)
			if x != tt.expectedX || y != tt.expectedY {
				t.Errorf("expected %d %d but %d %d", tt.expectedX, tt.expectedY, x, y)
			}
		})
	}
}

func TestRandom(t *testing.T) {
	const min, max = 1, 5
	expected := make(map[int]int, max)
	for i := min; i <= max; i++ {
		expected[i] = i
	}
	result := make(map[int]int, max)
	for len(expected) > 0 {
		key := random(min, max)
		if v, ok := expected[key]; ok {
			result[key] = v
			delete(expected, key)
		} else if _, ok := result[key]; !ok {
			t.Errorf("%d is not the expected value", key)
			break
		}
	}
}

func enemyActionTestInit(t *testing.T, mapPath string, enemyBuilder iEnemyBuilder) (*player, stage, error) {
	t.Helper()
	if err := termbox.Init(); err != nil {
		t.Error(err)
	}
	if err := termbox.Clear(termbox.ColorWhite, termbox.ColorBlack); err != nil {
		t.Error(err)
	}
	t.Cleanup(func() {
		termbox.Close()
	})
	stage := stage{
		mapPath:       mapPath,
		hunterBuilder: enemyBuilder,
		ghostBuilder:  enemyBuilder,
	}
	f, err := static.ReadFile(stage.mapPath)
	if err != nil {
		return nil, stage, err
	}
	b := createBuffer(bytes.NewReader(f))
	w := createWindow(b)
	if err = w.show(b); err != nil {
		return nil, stage, err
	}
	p := new(player)
	stage.plot(b, p)
	return p, stage, nil
}
