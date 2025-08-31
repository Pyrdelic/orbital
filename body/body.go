package body

import (
	"image/color"
	"image/draw"
	"math"
	"sync"

	"github.com/hajimehoshi/ebiten/v2"
)

// A body of mass. (Such as a planet.)
type Body struct {
	Image  *ebiten.Image
	X, Y   float64 // position of the body on 2D plane
	R      float64 // radius of the body
	M      float64 // mass of the body
	Vx, Vy float64 // velocity vector (actual movement)
	Fx, Fy float64 // force fector (effect of gravity)
	Color  color.Color
}

type Point struct {
	X, Y float64
}

// Returns x and y axis difference from a to b.
func PointDistanceXY(a, b Point) (float64, float64) {
	var xDiff, yDiff float64 // difference
	if a.X < b.X {
		xDiff = b.X - a.X
	} else if a.X > b.X {
		xDiff = a.X - b.X
	} else {
		xDiff = 0
	}
	if a.Y < b.Y {
		yDiff = b.Y - a.Y
	} else if a.Y > b.Y {
		yDiff = a.Y - b.Y
	} else {
		yDiff = 0
	}
	return xDiff, yDiff
}

func PointDistance(a, b Point) float64 {
	xDiff, yDiff := PointDistanceXY(a, b)
	distanceSq := xDiff*xDiff + yDiff*yDiff
	return math.Sqrt(distanceSq)
}

// // Returns gravityvector x, y from a to b.
// func atobGetDiff(a, b Point, f float64) (float64, float64) {
// 	var xDir, yDir float64 // direction multipliers from a to b
// 	var xDiff, yDiff float64
// 	if a.X < b.X {
// 		// a is left of b
// 		xDiff = b.X - a.X
// 		xDir = float64(1) // right
// 		if a.Y < b.Y {
// 			// a is above b (+, +)
// 			yDiff = float64(b.Y - a.Y)
// 			yDir = float64(1) // down

// 		} else if a.Y > b.Y {
// 			// a is below b (+, -)

// 			yDir = float64(-1) // up
// 		} else {
// 			// a and b are vertically aligned (+, 0)
// 			yDir = 0
// 		}
// 	} else if a.X > b.X {
// 		// a is right of b (-, ?)
// 		xDir = float64(-1) // left
// 		if a.Y < b.Y {
// 			// a is above b (-, +)
// 			yDir = float64(1) // down
// 		} else if a.Y > b.Y {
// 			// a is below b (-, -)
// 			yDir = float64(-1) // up
// 		} else {
// 			// a and b are vertically aligned (-, 0)
// 			yDir = 0
// 		}
// 	} else {
// 		// a and b are horizontally aligned (0, ?)
// 		xDir = 0
// 		if a.Y < b.Y {
// 			// a is above b (0, +)
// 			yDir = float64(1) // down
// 		} else if a.Y > b.Y {
// 			// a is below b (0, -)
// 			yDir = float64(-1) // up
// 		} else {
// 			// a and b are vertically aligned (0, 0)
// 			yDir = 0
// 		}
// 	}
// 	return xDir, yDir
// }

func getatobDiff(a, b Point) (float64, float64) {
	return b.X - a.X, b.Y - a.Y
}

const GravityConst float64 = 1.0

// Applies mutual gravity of two bodies to both bodies.
var applyGravityMutex sync.Mutex

func ApplyGravity(a, b *Body) {
	applyGravityMutex.Lock()
	r := PointDistance(a.Center(), b.Center())
	//f := GravityConst * ((a.M * b.M) / (r * r)) // length of the gravity vector
	rx, ry := getatobDiff(a.Center(), b.Center())
	//fmt.Printf("%.6f\t%.6f\t%.6f\n", rx, ry, r)
	f := GravityConst * ((a.M * b.M) / (r * r))
	fx := rx * f
	fy := ry * f

	a.Fx += fx
	a.Fy += fy
	b.Fx += -fx
	b.Fy += -fy
	applyGravityMutex.Unlock()
}

// NewBody returns a pointer to a new body.
func NewBody(x, y, r, m, vx, vy float64, c color.Color) *Body {
	b := Body{}
	b.X = x
	b.Y = y
	b.R = r
	b.M = m
	b.Vx = vx
	b.Vy = vy
	b.Image = ebiten.NewImage(
		int(2.0*r),
		int(2.0*r),
	)
	b.Image.Fill(c)
	return &b
}

func (b *Body) Update() {

	// apply force vector to velocity vector
	b.Vx += b.Fx
	b.Vy += b.Fy
	// apply velocity vector to position
	b.X += b.Vx
	b.Y += b.Vy
}

func (b *Body) Draw(screen *ebiten.Image, bodyDIO *ebiten.DrawImageOptions) {
	// dio := ebiten.DrawImageOptions{}
	// dio.GeoM.Translate(
	// 	b.X+zoomOffset/2,
	// 	b.Y+zoomOffset/2)
	//dio.GeoM.Translate(b.X, b.Y)
	// bodyDIO.ColorScale.Scale(
	// 	float32(b.Vx),
	// 	1.0,
	// 	float32(b.Vy),
	// 	1.0,
	// )
	screen.DrawImage(b.Image, bodyDIO)
}

func (b *Body) Center() Point {
	return Point{
		X: b.X + b.R,
		Y: b.Y + b.R,
	}
}

func drawCircle(img draw.Image, x0, y0, r int, c color.Color) {
	x, y, dx, dy := r-1, 0, 1, 1
	err := dx - (r * 2)

	for x > y {
		img.Set(x0+x, y0+y, c)
		img.Set(x0+y, y0+x, c)
		img.Set(x0-y, y0+x, c)
		img.Set(x0-x, y0+y, c)
		img.Set(x0-x, y0-y, c)
		img.Set(x0-y, y0-x, c)
		img.Set(x0+y, y0-x, c)
		img.Set(x0+x, y0-y, c)

		if err <= 0 {
			y++
			err += dy
			dy += 2
		}
		if err > 0 {
			x--
			dx += 2
			err += dx - (r * 2)
		}
	}
}
