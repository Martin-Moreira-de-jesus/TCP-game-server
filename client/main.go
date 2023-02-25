package main

import (
	"log"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

type Game struct {

	// The caracter's position
	myposx int
	myposy int

	// Other's pos
	otherpos int

	// Pipes informations
	pipeX      int
	obstacleY1 int
	obstacleY2 int

	//viewport
	viewport viewport
}

var GameState = Game{
	myposy:     450,
	myposx:     0,
	otherpos:   450,
	pipeX:      0,
	obstacleY1: 0,
	obstacleY2: 0,
}

var sprite *ebiten.Image
var pipeDown *ebiten.Image
var pipeUp *ebiten.Image
var background *ebiten.Image

func init() {
	img, _, err := ebitenutil.NewImageFromFile("gm.png")
	if err != nil {
		log.Fatal(err)
	}

	sprite = ebiten.NewImageFromImage(img)

	img2, _, err := ebitenutil.NewImageFromFile("back.png")
	if err != nil {
		log.Fatal(err)
	}

	img3, _, err := ebitenutil.NewImageFromFile("pipeUp.png")
	if err != nil {
		log.Fatal(err)
	}
	pipeUp = ebiten.NewImageFromImage(img3)

	img4, _, err := ebitenutil.NewImageFromFile("pipeDown.png")
	if err != nil {
		log.Fatal(err)
	}
	pipeDown = ebiten.NewImageFromImage(img4)

	background = ebiten.NewImageFromImage(img2)

}

func (g *Game) drawSprite(screen *ebiten.Image) {
	_, h := sprite.Bounds().Dx(), sprite.Bounds().Dy()

	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Scale(0.2, 0.2)
	op1.GeoM.Translate(0, -float64(h)/2.0) // set image to image center
	op1.GeoM.Translate(0, float64(GameState.myposy))
	op1.Filter = ebiten.FilterLinear
	screen.DrawImage(sprite, op1)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Scale(0.2, 0.2)
	op2.GeoM.Translate(0, -float64(h)/2.0) // set image to image center
	op2.GeoM.Translate(0, float64(GameState.otherpos))
	op2.Filter = ebiten.FilterLinear
	screen.DrawImage(sprite, op2)
}

func (g *Game) drawPipes(screen *ebiten.Image) {
	w, h := pipeUp.Bounds().Dx(), sprite.Bounds().Dy()

	op1 := &ebiten.DrawImageOptions{}
	op1.GeoM.Translate(float64(w), float64(h)/2.0) // set image to image center
	op1.GeoM.Translate(float64(GameState.pipeX), float64(GameState.obstacleY1)-800)
	//op.GeoM.Scale(0.3, 0.3)
	op1.Filter = ebiten.FilterLinear
	screen.DrawImage(pipeUp, op1)

	op2 := &ebiten.DrawImageOptions{}
	op2.GeoM.Translate(float64(w), float64(h)/2.0) // set image to image center
	op2.GeoM.Translate(float64(GameState.pipeX), float64(GameState.obstacleY2)-800)
	//op2.GeoM.Scale(0.3, 0.3)
	op2.Filter = ebiten.FilterLinear
	screen.DrawImage(pipeDown, op2)
}

type viewport struct {
	x16 int
	y16 int
}

func (p *viewport) Move() {
	w, h := background.Size()
	maxX16 := w * 16
	maxY16 := h * 16

	p.x16 += w / 32
	p.y16 += h / 32
	p.x16 %= maxX16
	p.y16 %= maxY16
}

func (p *viewport) Position() (int, int) {
	return p.x16, p.y16
}

func (g *Game) drawBackground(screen *ebiten.Image) {
	x16, y16 := g.viewport.Position()
	offsetX, _ := float64(-x16)/8, float64(-y16)/8

	// Draw bgImage on the screen repeatedly.
	const repeat = 6
	w, _ := background.Size()
	for j := 0; j < repeat; j++ {
		for i := 0; i < repeat; i++ {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(w*i), 0)
			op.GeoM.Translate(offsetX, 0)
			screen.DrawImage(background, op)
		}
	}
}

func (g *Game) Update() error {
	//c := SafeCounter{} Declared but not used
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		//g.myposy -= 100
		println("test", g.myposy)
		SendButtonPressed("up") // Shitty correction for message
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		//g.myposy += 100
		println("text", g.myposy)
		SendButtonPressed("down") // Shitty correction for message
	}
	g.viewport.Move()
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.drawBackground(screen)
	g.drawSprite(screen)
	g.drawPipes(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1500, 800
}

func main() {
	go ClientInfos{}.Client()
	ebiten.SetWindowSize(1500, 800)
	ebiten.SetWindowTitle("Dino nodino")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
