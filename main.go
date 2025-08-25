package main

import (
	"fmt"
	"image/color"
	"log"

	"github.com/Pyrdelic/orbital/body"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BodyCountMax int = 256
)

type Game struct {
	bodies []*body.Body
}

func (g *Game) Update() error {
	// update bodies
	// gravity, go through unique pairs
	for i := 0; i < len(g.bodies); i++ {
		for j := i + 1; j < len(g.bodies); j++ {
			// apply gravity between i and j
			// F = G*((m1*m2)/(r*r))
			//Fx := GravityConst * ((g.bodies[i].M * g.bodies[j].M) / g.bodies)
			body.ApplyGravity(g.bodies[i], g.bodies[j])
		}

		g.bodies[i].Update()
	}
	return nil
}

func (g *Game) debugPrintBodies(screen *ebiten.Image) {
	debugStr := ""
	for i := 0; i < len(g.bodies); i++ {
		bodyStr := fmt.Sprintf(
			"%d - x: %.2f, y: %.2f Fx: %.2f, Fy: %.2f",
			i,
			g.bodies[i].X,
			g.bodies[i].Y,
			g.bodies[i].Fx,
			g.bodies[i].Fy,
		)
		if i < len(g.bodies)-1 {
			bodyStr += "\n"
		}
		debugStr += bodyStr

	}
	ebitenutil.DebugPrint(screen, debugStr)
}

func (g *Game) Draw(screen *ebiten.Image) {

	for i := 0; i < len(g.bodies); i++ {
		//fmt.Println(i, g.bodies[i].X, g.bodies[i].Y)
		g.bodies[i].Draw(screen)
	}
	g.debugPrintBodies(screen)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func main() {
	if false {
		fmt.Println(body.PointDistance(
			body.Point{X: 10, Y: 5},
			body.Point{X: 0, Y: 0},
		))
		return
	}
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("Hello, World!")
	game := Game{}
	bodies := make([]*body.Body, 0)
	bodies = append(bodies, body.NewBody(
		float64(200.0), // x
		float64(200.0), // y
		float64(10),    // r
		float64(1.0),   // m
		float64(0.2),   // fx
		float64(-0.2),  // fy
		color.White,
	))
	bodies = append(bodies, body.NewBody(
		float64(100), // x
		float64(100), // y
		float64(10),  // r
		float64(1.0), // m
		float64(0.2), // fx
		float64(0.5), // fy
		color.RGBA{
			R: 127,
			G: 0,
			B: 127,
			A: 255,
		},
	))
	game.bodies = bodies
	if err := ebiten.RunGame(&game); err != nil {
		log.Fatal(err)
	}
}
