package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/Pyrdelic/orbital/body"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
)

const (
	BodyCountMax int = 256
)

type Game struct {
	bodies   []*body.Body
	vertices []ebiten.Vertex
	indices  []uint16
}

var ErrExit error = errors.New("Game exited")

var m1Hold bool = false

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrExit
	}
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !m1Hold {
		mx, my := ebiten.CursorPosition()
		g.addBody(mx, my)
		m1Hold = true
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		m1Hold = false
	}
	// update bodies
	// zero gravity vectors
	for i := 0; i < len(g.bodies); i++ {
		g.bodies[i].Fx, g.bodies[i].Fy = 0.0, 0.0
	}
	// calculate new gravity vectors, go through unique pairs
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
			"%d - x: %.1f, y: %.1f Fx: %.6f, Fy: %.6f",
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
	//debugStr += fmt.Sprintf("\npd: %.2f", body.PointDistance())
	ebitenutil.DebugPrint(screen, debugStr)
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.debugPrintBodies(screen)
	for i := 0; i < len(g.bodies); i++ {
		//fmt.Println(i, g.bodies[i].X, g.bodies[i].Y)
		g.bodies[i].Draw(screen)
	}
	// testVec := vec.Vec2D{X: 100, Y: 100, PosX: 150, PosY: 150}
	// testVec.Draw(screen)

}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 320, 240
}

func (g *Game) addBody(x, y int) {

	g.bodies = append(g.bodies, body.NewBody(
		float64(x),   // x
		float64(y),   // y
		float64(5),   // r
		float64(0.5), // m
		float64(0),   // vx
		float64(0),   // vy
		color.RGBA{
			R: uint8(rand.IntN(255)),
			G: uint8(rand.IntN(255)),
			B: uint8(126 + rand.IntN(126)),
			A: 255,
		},
	))
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
	ebiten.SetWindowTitle("N Body Problem")
	game := Game{}
	bodies := make([]*body.Body, 0)
	// bodies = append(bodies, body.NewBody(
	// 	float64(100), // x
	// 	float64(150), // y
	// 	float64(5),   // r
	// 	float64(0.5), // m
	// 	float64(0.1), // vx
	// 	float64(0),   // vy
	// 	color.RGBA{
	// 		R: 255,
	// 		G: 0,
	// 		B: 0,
	// 		A: 255,
	// 	},
	// ))
	// bodies = append(bodies, body.NewBody(
	// 	float64(200), // x
	// 	float64(150), // y
	// 	float64(5),   // r
	// 	float64(0.5), // m
	// 	float64(0),   // vx
	// 	float64(0),   // vy
	// 	color.RGBA{
	// 		R: 0,
	// 		G: 255,
	// 		B: 0,
	// 		A: 255,
	// 	},
	// ))
	// bodies = append(bodies, body.NewBody(
	// 	float64(150), // x
	// 	float64(175), // y
	// 	float64(5),   // r
	// 	float64(0.5), // m
	// 	float64(0.0), // vx
	// 	float64(0.0), // vy
	// 	color.RGBA{
	// 		R: 0,
	// 		G: 0,
	// 		B: 255,
	// 		A: 255,
	// 	},
	// ))
	game.bodies = bodies
	if err := ebiten.RunGame(&game); err != nil {
		if err == ErrExit {
			fmt.Println(err)
			return
		}
		log.Fatal(err)
	}
}
