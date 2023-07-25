package main

import (
	"bytes"
	"testing"

	termbox "github.com/nsf/termbox-go"
)

const enemyTestMapPath = "files/test/enemy/"

func TestEnemyControl(t *testing.T) {
	cases := map[string]struct {
		enemyBuilder  iEnemyBuilder
		expectedMoves int
		mapFileName   string
	}{
		"hunter with obstacle": {newEnemyBuilder().defaultHunter(), 5, "hunter_with_obstacle.txt"},
		"ghost with obstacle":  {newEnemyBuilder().defaultGhost(), 6, "ghost_with_obstacle.txt"},
		"tricky":               {newEnemyBuilder().defaultHunter().strategize(&tricky{}), 5, "hunter_with_obstacle.txt"},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p, stage, err := enemyActionTestInit(t, enemyTestMapPath+tt.mapFileName, tt.enemyBuilder)
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
			if name != "tricky" && count != tt.expectedMoves {
				t.Errorf("expected %d but %d", tt.expectedMoves, count)
			}
		})
	}
}

// Test enemies do not overlap each other.
func TestEnemiesOverlap(t *testing.T) {
	cases := map[string]struct {
		enemyBuilder  iEnemyBuilder
		expectedMoves int
		mapFileName   string
	}{
		"hunter with enemies": {newEnemyBuilder().defaultHunter(), 5, "hunter_with_enemies.txt"},
		"ghost with enemies":  {newEnemyBuilder().defaultGhost(), 10, "ghost_with_enemies.txt"},
	}
	for name, tt := range cases {
		tt := tt
		t.Run(name, func(t *testing.T) {
			p, stage, err := enemyActionTestInit(t, enemyTestMapPath+tt.mapFileName, tt.enemyBuilder)
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
			if count != tt.expectedMoves {
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
		"up":    {5, 1, 5, 2},
		"down":  {5, 5, 5, 4},
		"left":  {3, 3, 4, 3},
		"right": {7, 3, 6, 3},
	}
	p, stage, err := enemyActionTestInit(t, enemyTestMapPath+"hunter.txt", newEnemyBuilder().defaultHunter())
	if err != nil {
		t.Error(err)
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
