package body

import (
	"image/color"
	"image/draw"
	"math"

	"github.com/hajimehoshi/ebiten/v2"
)

// A body of mass. (Such as a planet.)
type Body struct {
	Image  *ebiten.Image
	X, Y   float64 // position of the body on 2D plane
	R      float64 // radius of the body
	M      float64 // mass of the body
	Fx, Fy float64 // force fector
	Color  color.Color
}

type Point struct {
	X, Y float64
}

const GravityConst float64 = 1.0

// Returns x and y axis differences.
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

// Applies mutual gravity of two bodies.
func ApplyGravity(a, b *Body) {

	// var a, b *Body
	// if bodyA.X < bodyB.X {
	// 	a, b = bodyA, bodyB
	// } else {
	// 	a, b = bodyB, bodyA
	// }
	// a is now left of b
	rx, ry := PointDistanceXY(a.Center(), b.Center())
	//fmt.Println(rx, ry)
	var fx, fy float64
	fx = GravityConst * (a.M * b.M) / (rx * rx) // zero danger
	fy = GravityConst * (a.M * b.M) / (ry * ry) // zero danger
	aCenter := a.Center()
	bCenter := b.Center()
	if aCenter.X < bCenter.X {
		// a is left of b
		a.Fx += fx
		b.Fx += -fx
	} else if aCenter.X > bCenter.X {
		// b is left of a
		a.Fx += -fx
		b.Fx += fx
	} else {
		// exactly same x position
	}
	if aCenter.Y < bCenter.Y {
		// a is above b
		//a.Fy, b.Fy = fy, -fy
		a.Fy += fy
		b.Fy += -fy
	} else if a.Y > b.Y {
		// a is below b
		a.Fy += -fy
		b.Fy += fy
	} else {
		// exactly same y position
	}
}

// NewBody returns a pointer to a new body.
func NewBody(x, y, r, m, fx, fy float64, c color.Color) *Body {
	b := Body{}
	b.X = x
	b.Y = y
	b.R = r
	b.M = m
	b.Fx = fx
	b.Fy = fy
	b.Image = ebiten.NewImage(
		int(2.0*r),
		int(2.0*r),
	)
	b.Image.Fill(c)
	return &b
}

func (b *Body) Update() {
	// apply force vector to position
	b.X += b.Fx
	b.Y += b.Fy
}

func (b *Body) Draw(screen *ebiten.Image) {
	dio := ebiten.DrawImageOptions{}
	dio.GeoM.Translate(b.X, b.Y)
	screen.DrawImage(b.Image, &dio)
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
