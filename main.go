package main

import (
	"errors"
	"fmt"
	"image/color"
	"log"
	"math/rand/v2"
	"sync"

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
	Bodies      []*body.Body
	uniquePairs [][2]*body.Body
	mutex       sync.Mutex
	Background  *ebiten.Image // only used for the trail-effect (which is broken rn)

	ViewPortScale   int
	ViewPortOffsetX int
	ViewPortOffsetY int
}

var ErrExit error = errors.New("Game exited")
var m1Holding bool = false
var slingStartX, slingStartY, slingEndX, slingEndY float64

func (g *Game) ProcessInput() error {
	// ESC exit
	if ebiten.IsKeyPressed(ebiten.KeyEscape) {
		return ErrExit
	}
	// zoom controls
	// zoom out
	if ebiten.IsKeyPressed(config.KeyBindZoomOut) {
		g.ViewPortScale += config.CameraZoomSpeed
	}
	if ebiten.IsKeyPressed(config.KeyBindZoomIn) {
		g.ViewPortScale -= config.CameraZoomSpeed
		if g.ViewPortScale < 0 {
			g.ViewPortScale = 0
		}
	}
	if inpututil.IsKeyJustPressed(ebiten.KeyEnter) {
		g.spawnToRandom(100)
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

	// adding a body via mouse click
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		m1Holding = true
		//mx, my := ebiten.CursorPosition()
		mx, my := ebiten.CursorPosition()
		slingStartX, slingStartY = float64(mx), float64(my)
		fmt.Println("Sling start at x:", mx, "y:", my)
		//g.addBody(mx, my)
	}
	if inpututil.IsMouseButtonJustReleased(ebiten.MouseButtonLeft) {
		m1Holding = false
		mx, my := ebiten.CursorPosition()
		slingEndX, slingEndY = float64(mx), float64(my)
		// get difference between sling start and sling end
		dx := slingStartX - slingEndX
		dy := slingStartY - slingEndY
		g.addBody(int(slingStartX), int(slingStartY), dx*config.SlingSpeed, dy*config.SlingSpeed)
	}
	return nil
}

func (g *Game) Update() error {
	//start := time.Now()
	if err := g.ProcessInput(); err != nil {
		if err == ErrExit {
			return ErrExit
		}
	}

	// update bodies

	// zero out the bodies' gravity vectors
	for i := 0; i < len(g.Bodies); i++ {
		g.Bodies[i].Fx, g.Bodies[i].Fy = 0.0, 0.0
	}

	// Calculate new gravity vectors between
	// all uniquely paired bodies.
	threadCount := 4
	if false { // multithreading is under construction
		workload := len(g.uniquePairs) / threadCount
		remWorkload := len(g.uniquePairs) % threadCount
		var pairsWg sync.WaitGroup
		for i := 0; i < threadCount; i++ {
			pairsWg.Add(1)
			go g.calcPairsGoroutine(i, &pairsWg, i*workload, workload+workload*i)
		}
		// calculate the possible remained here in the main thread
		// while waiting other threads to finish.
		for i := len(g.uniquePairs) - remWorkload - 1; i < len(g.uniquePairs); i++ {
			index := len(g.uniquePairs) - 1 - remWorkload + i
			if index < len(g.uniquePairs) && index >= 0 {
				fmt.Println("Main index:", index)
				body.ApplyGravity(g.uniquePairs[index][0], g.uniquePairs[index][1])
			}

		}
		pairsWg.Wait()
	} else {
		// calculate gravity in a single thread
		for i := 0; i < len(g.uniquePairs); i++ {
			body.ApplyGravity(g.uniquePairs[i][0], g.uniquePairs[i][1])
		}
	}

	// update bodies
	for i := 0; i < len(g.Bodies); i++ {
		g.Bodies[i].Update()
	}

	return nil // nil error
}

// calculates gravity between unique pairs, from indices a to b
func (g *Game) calcPairsGoroutine(id int, wg *sync.WaitGroup, a, b int) {
	for i := a; i < b; i++ {

		body.ApplyGravity(g.uniquePairs[a][0], g.uniquePairs[a][1])
	}
	wg.Done()
}

func (g *Game) spawnToRandom(n int) {
	for i := 0; i < n; i++ {
		g.addBody(rand.IntN(1000), rand.IntN(1000), 0, 0)
	}
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

func (g *Game) addBody(x, y int, vx, vy float64) {

	g.Bodies = append(g.Bodies, body.NewBody(
		float64((x-g.ViewPortScale/2)+g.ViewPortOffsetX), // x
		float64((y-g.ViewPortScale/2)+g.ViewPortOffsetY), // y
		float64(5),   // r
		float64(0.5), // m
		vx,           // vx
		vy,           // vy
		color.RGBA{
			R: uint8(rand.IntN(255)),
			G: uint8(rand.IntN(255)),
			B: uint8(126 + rand.IntN(126)),
			A: 255,
		},
	))
	// Every time a new body is added, pair it with already added ones.
	if len(g.Bodies) >= 2 {
		for i := 0; i < len(g.Bodies)-1; i++ {
			//bodieslen := len(g.Bodies)
			pair := make([][2]*body.Body, 1)
			pair[0][0] = g.Bodies[i]
			pair[0][1] = g.Bodies[len(g.Bodies)-1]
			g.uniquePairs = append(g.uniquePairs, pair...)

		}
		for i := 0; i < len(g.uniquePairs); i++ {
			//fmt.Println(g.uniquePairs[i])
		}
		//fmt.Println()
	}
}

func main() {
	ebiten.SetWindowSize(640, 480)
	ebiten.SetWindowTitle("N Body Problem")

	g := Game{}
	g.Bodies = make([]*body.Body, 0)
	g.uniquePairs = make([][2]*body.Body, 0)
	//g.spawnToRandom(1000)
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
