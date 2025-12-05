package main

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
	"image/color"
	"log"
)

const (
	screenWidth  = 640
	screenHeight = 480
	fishSpeed    = 3
	netSpeed     = 6
)

type Object struct {
	X, Y, W, H int
}

type Net struct {
	Object
}

type Fish struct {
	Object
	dxdt int // x velocity per tick
	dydt int // y velocity per tick
}

type Game struct {
	net       Net
	fish      Fish
	score     int
	highScore int
}

func main() {
	ebiten.SetWindowTitle("Blub in Go")
	ebiten.SetWindowSize(screenWidth, screenHeight)
	net := Net{
		Object: Object{
			X: 600,
			Y: 200,
			W: 15,
			H: 100,
		},
	}
	fish := Fish{
		Object: Object{
			X: 0,
			Y: 0,
			W: 15,
			H: 15,
		},
		dxdt: fishSpeed,
		dydt: fishSpeed,
	}
	g := Game{
		net:  net,
		fish: fish,
	}

	err := ebiten.RunGame(&g)
	if err != nil {
		log.Fatal(err)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Draw(screen *ebiten.Image) {
	vector.FillRect(screen,
		float32(0), float32(0),
		float32(screenWidth), float32(screenHeight),
		color.RGBA{41, 170, 240, 255}, false,
	)

	vector.FillRect(screen,
		float32(g.net.X), float32(g.net.Y),
		float32(g.net.W), float32(g.net.H),
		color.RGBA{212, 51, 52, 255}, false,
	)

	vector.FillRect(screen,
		float32(g.fish.X), float32(g.fish.Y),
		float32(g.fish.W), float32(g.fish.H),
		color.RGBA{212, 51, 52, 255}, false,
	)

	scoreStr := "Score: " + fmt.Sprint(g.score)
	text.Draw(screen, scoreStr, basicfont.Face7x13, 10, 10, color.Black)

	highScoreStr := "High Score: " + fmt.Sprint(g.highScore)
	text.Draw(screen, highScoreStr, basicfont.Face7x13, 10, 30, color.Black)
}

func (g *Game) Update() error {
	g.net.MoveOnKeyPress()
	g.fish.Move()
	g.CollideWithWall()
	g.CollideWithNet()
	return nil
}

func (n *Net) MoveOnKeyPress() {
	if ebiten.IsKeyPressed(ebiten.KeyArrowDown) {
		n.Y += netSpeed
	}
	if ebiten.IsKeyPressed(ebiten.KeyArrowUp) {
		n.Y -= netSpeed
	}
}

func (f *Fish) Move() {
	f.X += f.dxdt
	f.Y += f.dydt
}

func (g *Game) Reset() {
	g.fish.X = 0
	g.fish.Y = 0

	g.score = 0
}

func (g *Game) CollideWithWall() {

	// Right wall causes a game over
	if g.fish.X >= screenWidth {
		g.Reset()
	} else if g.fish.X < 0 {
		g.fish.dxdt = fishSpeed
	} else if g.fish.Y <= 0 {
		g.fish.dydt = fishSpeed
	} else if g.fish.Y >= screenHeight {
		g.fish.dydt = -fishSpeed
	}
}

func (g *Game) CollideWithNet() {
	if g.fish.X >= g.net.X && g.fish.Y >= g.net.Y && g.fish.Y <= g.net.Y+g.net.H {
		g.fish.dxdt = -g.fish.dxdt
		g.score++
		if g.score > g.highScore {
			g.highScore = g.score
		}
	}
}
