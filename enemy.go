package main

import (
	"math"
	"math/rand"
	"time"

	termbox "github.com/nsf/termbox-go"
)

type iEnemy interface {
	getPosition() (x, y int)
	setPosition(x, y int)
	getDisplayFormat() (rune, termbox.Attribute)
	think(p *player) (int, int)
	move(x, y int)
	hasCaptured(p *player)
	eval(p *player, x, y int) float64
}
type enemy struct {
	x            int
	y            int
	char         rune
	color        termbox.Attribute
	waitingTime  int
	oneActionInN int
	canMove      func(int, int) bool
	strategy
	underRune
}
type hunter struct {
	enemy
}
type ghost struct {
	enemy
}

type strategy interface {
	eval(p *player, x, y int) float64
}
type assault struct{}
type tricky struct{}

type underRune struct {
	char  rune
	color termbox.Attribute
}

type iEnemyBuilder interface {
	position(int, int) iEnemyBuilder
	displayFormat(rune, string) iEnemyBuilder
	speed(int) iEnemyBuilder
	strategize(strategy) iEnemyBuilder
	movable(func(int, int) bool) iEnemyBuilder
	defaultHunter() iEnemyBuilder
	defaultGhost() iEnemyBuilder
	build() iEnemy
}
type enemyBuilder struct {
	x            int
	y            int
	char         rune
	color        termbox.Attribute
	waitingTime  int
	oneActionInN int
	canMove      func(int, int) bool
	strategy     strategy
}

func (eb *enemyBuilder) position(x int, y int) iEnemyBuilder {
	eb.x = x
	eb.y = y
	return eb
}
func (eb *enemyBuilder) displayFormat(r rune, s string) iEnemyBuilder {
	eb.char = r
	switch s {
	case "RED":
		eb.color = termbox.ColorRed
	case "BLUE":
		eb.color = termbox.ColorBlue
	}
	return eb
}
func (eb *enemyBuilder) speed(i int) iEnemyBuilder {
	eb.waitingTime = i
	eb.oneActionInN = i
	return eb
}
func (eb *enemyBuilder) movable(fn func(int, int) bool) iEnemyBuilder {
	eb.canMove = fn
	return eb
}
func (eb *enemyBuilder) strategize(s strategy) iEnemyBuilder {
	eb.strategy = s
	return eb
}
func (eb *enemyBuilder) defaultHunter() iEnemyBuilder {
	fn := func(x, y int) bool {
		return !isCharWall(x, y) && !isCharEnemy(x, y)
	}
	return eb.displayFormat(chHunter, "RED").speed(1).movable(fn).strategize(&assault{})
}
func (eb *enemyBuilder) defaultGhost() iEnemyBuilder {
	fn := func(x, y int) bool {
		return !isCharBorder(x, y) && !isCharEnemy(x, y)
	}
	return eb.displayFormat(chGhost, "BLUE").speed(2).movable(fn).strategize(&assault{})
}
func newEnemyBuilder() iEnemyBuilder {
	return &enemyBuilder{}
}
func (eb *enemyBuilder) build() iEnemy {
	return &enemy{
		x:            eb.x,
		y:            eb.y,
		char:         eb.char,
		color:        eb.color,
		waitingTime:  eb.waitingTime,
		oneActionInN: eb.oneActionInN,
		canMove:      eb.canMove,
		strategy:     eb.strategy,
	}
}

func (e *enemy) getPosition() (x, y int) {
	return e.x, e.y
}
func (e *enemy) setPosition(x, y int) {
	e.x = x
	e.y = y
}

func (e *enemy) getDisplayFormat() (rune, termbox.Attribute) {
	return e.char, e.color
}

func (e *enemy) think(p *player) (int, int) {
	x, y := e.getPosition()
	// Calculate the evaluation value for movement
	up := e.eval(p, x, y-1)
	down := e.eval(p, x, y+1)
	left := e.eval(p, x-1, y)
	right := e.eval(p, x+1, y)

	// The evaluation value is the distance between the player and the enemy
	// So prefer lower (closer) values
	if up <= down && up <= left && up <= right {
		return x, y - 1 // up
	} else if down <= left && down <= right {
		return x, y + 1 // down
	} else if left <= right {
		return x - 1, y // left
	} else {
		return x + 1, y // right
	}
}

func (e *enemy) move(x, y int) {
	e.waitingTime--
	if e.waitingTime <= 0 && e.canMove(x, y) {
		winWidth, _ := termbox.Size()
		// Set the original character in the original cell
		termbox.SetCell(e.x, e.y, e.underRune.char, e.underRune.color, termbox.ColorBlack)
		// Retains destination cell information
		// Because it is necessary to set the original character at the next move
		cell := termbox.CellBuffer()[(winWidth*y)+x]
		e.x = x
		e.y = y
		e.underRune.char = cell.Ch
		e.underRune.color = cell.Fg
		termbox.SetCell(x, y, e.char, e.color, termbox.ColorBlack)
		e.waitingTime = e.oneActionInN
	}
}

func (e *enemy) hasCaptured(p *player) {
	if e.x == p.x && e.y == p.y {
		gameState = lose
	}
}

func (e *enemy) eval(p *player, x, y int) float64 {
	if !e.canMove(x, y) {
		// Returns a large enough value if it can't move
		return 1000
	}
	return e.strategy.eval(p, x, y)
}
func (s *assault) eval(p *player, x, y int) float64 {
	// Distance between two points
	return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
}
func (s *tricky) eval(p *player, x, y int) float64 {
	if random(0, 4) == 0 {
		// Ignore 1 out of 5 times
		return 0
	} else {
		return math.Sqrt(math.Pow(float64(p.y-y), 2) + math.Pow(float64(p.x-x), 2))
	}
}

// Return a value between min and max
// e.g. random(0, 4) returns 0,1,2,3,4
func random(min, max int) int {
	rand.New(rand.NewSource(time.Now().UnixNano()))
	return rand.Intn(max-min+1) + min
}
