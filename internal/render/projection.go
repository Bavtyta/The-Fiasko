package render

import "TheFiaskoTest/internal/core"

func Project(point core.Vec3, cam *Camera) (float64, float64, float64) {

	// переводим в координаты камеры
	relX := point.X - cam.Position.X
	relY := point.Y - cam.Position.Y
	relZ := point.Z - cam.Position.Z

	if relZ <= 0.1 {
		return 0, 0, 0
	}

	scale := cam.FocalLength / relZ

	screenX := cam.ScreenW/2 + relX*scale
	screenY := cam.HorizonY - relY*scale

	return screenX, screenY, scale
}
