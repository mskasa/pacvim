package main

import (
	"log"

	"github.com/hajimehoshi/ebiten/v2"

	"github.com/masahiro-kasatani/pacvim/renderer"
)

func main() {
	ebiten.SetWindowSize(renderer.ScreenWidth, renderer.ScreenHeight)
	ebiten.SetWindowTitle("PacVim")
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeDisabled)

	g, err := NewGame()
	if err != nil {
		log.Fatal(err)
	}
	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
