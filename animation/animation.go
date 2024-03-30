package main

import (
	"fmt"
	_ "image/png"

	"github.com/hajimehoshi/ebiten/v2"
)

type Game struct {
	mouseX int
	mouseY int
}

func (g *Game) Update() error {

	ebiten.SetWindowTitle(fmt.Sprintf("Game Title | %.2ffps", ebiten.ActualFPS()))

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1280, 720
}

func main() {
	ebiten.SetWindowSize(1280, 720)
	ebiten.SetWindowTitle("Animation")
	if err := ebiten.RunGame(&Game{}); err != nil {
		panic(err)
	}
}
