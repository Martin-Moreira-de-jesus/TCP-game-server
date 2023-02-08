package main

import (
	"image/color"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {

	// The gopher's position
	x16  int
	y16  int
	vy16 int

	// Camera
	cameraX int
	cameraY int
}

var sprite *ebiten.Image

func init() {
	img, _, err := ebitenutil.NewImageFromFile("gm.png")
	if err != nil {
		log.Fatal(err)
	}
	sprite = ebiten.NewImageFromImage(img)
}

func (g *Game) init() {
	g.x16 = 0
	g.y16 = 100 * 16
	g.cameraX = -240
	g.cameraY = 0
}

func (g *Game) drawSprite(screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, h := sprite.Bounds().Dx(), sprite.Bounds().Dy()
	op.GeoM.Translate(-float64(w)/2.0, -float64(h)/2.0)
	op.GeoM.Translate(float64(w)/2.0, float64(h)/2.0)
	op.GeoM.Translate(0, float64(g.y16/16.0))
	op.Filter = ebiten.FilterLinear
	screen.DrawImage(sprite, op)
}

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		g.y16 -= 50
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		g.y16 += 50
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	screen.Fill(color.RGBA{0, 0, 0, 0})
	g.drawSprite(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	ebiten.SetWindowSize(900, 900)
	ebiten.SetWindowTitle("Dino nodino")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
