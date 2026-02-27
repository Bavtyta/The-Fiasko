package render

import "TheFiaskoTest/internal/core"

type Camera struct {
	Position    core.Vec3
	FocalLength float64
	HorizonY    float64
	ScreenW     float64
	ScreenH     float64
}

func NewCamera(screenW, screenH float64) *Camera {
	return &Camera{
		Position: core.Vec3{
			X: 0,
			Y: 1.5,
			Z: -2,
		},
		FocalLength: screenH,
		HorizonY:    screenH * 2.0 / 3.0,
		ScreenW:     screenW,
		ScreenH:     screenH,
	}
}
