package config

import "github.com/hajimehoshi/ebiten/v2"

const (
	InnerWidth  int = 320
	InnerHeight int = 240

	CameraSpeed     int = 5
	CameraZoomSpeed int = 5
)

var (
	KeyBindZoomOut = ebiten.KeyPageDown
	KeyBindZoomIn  = ebiten.KeyPageUp

	KeyBindMoveCamLeft  = ebiten.KeyArrowLeft
	KeyBindMoveCamRight = ebiten.KeyArrowRight
	KeyBindMoveCamUp    = ebiten.KeyArrowUp
	KeyBindMoveCamDown  = ebiten.KeyArrowDown
)
