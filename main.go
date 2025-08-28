package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"

	"github.com/Pyrdelic/orbital/body"
	"github.com/Pyrdelic/orbital/config"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	BodyCountMax int = 256
)

type Game struct {
	Bodies     []*body.Body
	Background *ebiten.Image // only used for the trail-effect (which is broken rn)

	ViewPortScale   int
	ViewPortOffsetX int
	ViewPortOffsetY int
}

var ErrExit error = errors.New("Game exited")

var m1Hold bool = false

func (g *Game) Update() error {
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrExit
	}

	// zoom controls
	// zoom out
	if inpututil.IsKeyJustPressed(config.KeyBindZoomOut) {
		g.ViewPortScale += 20
	}
	if inpututil.IsKeyJustPressed(config.KeyBindZoomIn) {
		g.ViewPortScale -= 20
		if g.ViewPortScale < 0 {
			g.ViewPortScale = 0
		}
	}

	// camera movement
	if ebiten.IsKeyPressed(config.KeyBindMoveCamDown) {
		g.ViewPortOffsetY += config.CameraSpeed
	}
	if ebiten.IsKeyPressed(config.KeyBindMoveCamUp) {
		g.ViewPortOffsetY -= config.CameraSpeed
	}
	if ebiten.IsKeyPressed(config.KeyBindMoveCamLeft) {
		g.ViewPortOffsetX -= config.CameraSpeed
	}
	if ebiten.IsKeyPressed(config.KeyBindMoveCamRight) {
		g.ViewPortOffsetX += config.CameraSpeed
	}

	// adding a body
	if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && !m1Hold {
		mx, my := ebiten.CursorPosition()
		fmt.Println("Add body to x:", mx, "y:", my)
		g.addBody(mx, my)
		m1Hold = true
	} else if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
		m1Hold = false
	}

	// update bodies
	// zero out the gravity vectors
	for i := 0; i < len(g.Bodies); i++ {
		g.Bodies[i].Fx, g.Bodies[i].Fy = 0.0, 0.0
	}
	// calculate new gravity vectors, go through unique pairs
	for i := 0; i < len(g.Bodies); i++ {
		for j := i + 1; j < len(g.Bodies); j++ {
			// apply gravity between i and j
			// F = G*((m1*m2)/(r*r))
			//Fx := GravityConst * ((g.bodies[i].M * g.bodies[j].M) / g.bodies)
			body.ApplyGravity(g.Bodies[i], g.Bodies[j])
		}

		g.Bodies[i].Update()
	}
	return nil
}

func (g *Game) debugPrintBodies(screen *ebiten.Image) {
	debugStr := ""
	for i := 0; i < len(g.Bodies); i++ {
		bodyStr := fmt.Sprintf(
			"%d - x: %.1f, y: %.1f Fx: %.6f, Fy: %.6f",
			i,
			g.Bodies[i].Center().X,
			g.Bodies[i].Center().Y,
			g.Bodies[i].Fx,
			g.Bodies[i].Fy,
		)
		if i < len(g.Bodies)-1 {
			bodyStr += "\n"
		}
		debugStr += bodyStr

	}
	//debugStr += fmt.Sprintf("\npd: %.2f", body.PointDistance())
	ebitenutil.DebugPrint(screen, debugStr)
}

func (g *Game) Draw(screen *ebiten.Image) {
	//g.debugPrintBodies(screen)
	var fadeImg *ebiten.Image = ebiten.NewImage(config.InnerWidth, config.InnerHeight)
	//fadeImg := ebiten.NewImage(config.InnerWidth, config.InnerHeight)
	fadeImg.Fill(color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	})
	bgFadeOpts := ebiten.DrawImageOptions{}
	bgFadeOpts.ColorScale.Scale(1, 1, 1, 0.01)
	//bgFadeOpts.ColorScale.ScaleWithColor(color.RGBA{A: 1})
	g.Background.DrawImage(fadeImg, &bgFadeOpts)
	// screenOpts := ebiten.DrawImageOptions{}
	zoomOffset := float64(g.ViewPortScale)

	for i := 0; i < len(g.Bodies); i++ {
		// g.Bodies[i].Draw(screen, zoomOffset)
		bodyDIO := ebiten.DrawImageOptions{}
		bodyDIO.GeoM.Translate(
			(g.Bodies[i].X+zoomOffset/2)-float64(g.ViewPortOffsetX),
			(g.Bodies[i].Y+zoomOffset/2)-float64(g.ViewPortOffsetY),
		)
		g.Bodies[i].Draw(screen, &bodyDIO)
	}
	//fmt.Println(g.ViewPortOffsetX, g.ViewPortOffsetY)

	//fmt.Println("ViewportScale:", g.ViewPortScale, "ZoomOffset:", zoomOffset)
	//screen.DrawImage(g.Background, &screenOpts)
	//screen.DrawImage()
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Bodies: %d", len(g.Bodies)))
}

var aspectRatioX float64 = 4
var aspectRatioY float64 = 3

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return config.InnerWidth + g.ViewPortScale*int(aspectRatioX),
		config.InnerHeight + g.ViewPortScale*int(aspectRatioY)
}

func (g *Game) addBody(x, y int) {

	g.Bodies = append(g.Bodies, body.NewBody(
		float64((x-g.ViewPortScale/2)+g.ViewPortOffsetX), // x
		float64((y-g.ViewPortScale/2)+g.ViewPortOffsetY), // y
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
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("N Body Problem")

	g := Game{}
	g.Bodies = make([]*body.Body, 0)
	g.ViewPortScale = 0
	g.Background = ebiten.NewImage(config.InnerWidth, config.InnerHeight)
	g.Background.Fill(color.RGBA{
		R: 0,
		G: 0,
		B: 0,
		A: 255,
	})

	if err := ebiten.RunGame(&g); err != nil {
		if err == ErrExit {
			fmt.Println(err)
			return
		}
		log.Fatal(err)
	}
}
