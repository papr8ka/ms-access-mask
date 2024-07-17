package main

import (
	"access-mask/gui"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.design/x/clipboard"
	"log"
)

func main() {
	if err := clipboard.Init(); err != nil {
		log.Fatal("could not initialize clipboard:", err)
	}

	ebiten.SetWindowSize(1055, 600)
	ebiten.SetWindowResizingMode(ebiten.WindowResizingModeEnabled)
	ebiten.SetWindowTitle("Access mask formatter")
	if err := ebiten.RunGame(gui.New()); err != nil {
		log.Fatal("could not create window:", err)
	}
}
