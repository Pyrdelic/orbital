package vec

import (
	"image"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// // Apply applies a ForceVector to caller.
// func (f *ForceVector) ApplyForce(fv ForceVector) {
// 	f.Fx += fv.Fx
// 	f.Fy += fv.Fy
// }

type Vec2D struct {
	X, Y, PosX, PosY float64
}

var whiteImage = ebiten.NewImage(3, 3)
var whiteSubImage = whiteImage.SubImage(image.Rect(1, 1, 2, 2)).(*ebiten.Image)

// Sum returns v1 + v2
func Sum(v1, v2 Vec2D) Vec2D {
	return Vec2D{
		X: v1.X + v2.X,
		Y: v1.Y + v2.Y,
	}
}

func (v *Vec2D) Draw(screen *ebiten.Image) {
	var path vector.Path
	path.MoveTo(50, 50)
	path.LineTo(100, 100)

	op := &vector.StrokeOptions{}
	op.Width = 1
	op.LineJoin = vector.LineJoinRound

	vertices := make([]ebiten.Vertex, 0)
	indices := make([]uint16, 0)

	vertices, indices = path.AppendVerticesAndIndicesForStroke(vertices[:0], indices[:0], op)
	for i := range vertices {
		vertices[i].SrcX = 1
		vertices[i].SrcY = 1
		vertices[i].ColorR = 1
		vertices[i].ColorG = 1
		vertices[i].ColorB = 1
		vertices[i].ColorA = 1
	}
	screen.DrawTriangles(vertices, indices, whiteSubImage, &ebiten.DrawTrianglesOptions{
		AntiAlias: false,
	})
	//fmt.Println(vertices)
	//fmt.Println(indices)
}
