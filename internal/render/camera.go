package render

import (
	"TheFiaskoTest/internal/config"
	"TheFiaskoTest/internal/core"
	"math"
)

type Camera struct {
	position    core.Vec3
	focalLength float64
	horizonY    float64
	screenW     float64
	screenH     float64
	config      config.CameraConfig
}

// Position returns the camera position
func (c *Camera) Position() core.Vec3 {
	return c.position
}

// FocalLength returns the camera focal length
func (c *Camera) FocalLength() float64 {
	return c.focalLength
}

// HorizonY returns the horizon Y coordinate
func (c *Camera) HorizonY() float64 {
	return c.horizonY
}

// ScreenW returns the screen width
func (c *Camera) ScreenW() float64 {
	return c.screenW
}

// ScreenH returns the screen height
func (c *Camera) ScreenH() float64 {
	return c.screenH
}

func NewCamera(screenW, screenH float64, cfg config.CameraConfig) *Camera {
	return &Camera{
		position: core.Vec3{
			X: 0,
			Y: cfg.DefaultPositionY,
			Z: cfg.DefaultPositionZ,
		},
		focalLength: screenH,
		horizonY:    screenH * cfg.HorizonRatio,
		screenW:     screenW,
		screenH:     screenH,
		config:      cfg,
	}
}

// Project проецирует 3D точку на экран с учётом конфигурации камеры
func (c *Camera) Project(point core.Vec3) (float64, float64, float64) {
	// Переводим в координаты камеры
	relX := point.X - c.position.X
	relY := point.Y - c.position.Y
	relZ := point.Z - c.position.Z

	if relZ <= 0.1 {
		return 0, 0, 0
	}

	// Базовая перспективная проекция
	scale := c.focalLength / relZ

	// === ГОРИЗОНТАЛЬНОЕ ИСКРИВЛЕНИЕ (влево) ===
	curveOffset := 0.0
	if c.config.CurveStrength > 0 && c.config.CurveDepth > 0 {
		// Нормализуем глубину (0..1) с учётом максимальной глубины
		normalizedDepth := math.Min(relZ/c.config.CurveDepth, 1.0)
		// Применяем квадратичную функцию для плавного закругления
		curveOffset = -c.config.CurveStrength * normalizedDepth * normalizedDepth * c.screenW
	}

	// === ВЕРТИКАЛЬНОЕ ИСКРИВЛЕНИЕ (подъём и опускание) ===
	verticalOffset := 0.0
	if c.config.VerticalCurveStrength > 0 && c.config.CurveDepth > 0 {
		normalizedDepth := math.Min(relZ/c.config.CurveDepth, 1.0)
		peakPoint := c.config.VerticalCurvePeak

		if normalizedDepth <= peakPoint {
			// Фаза подъёма: от 0 до пика (параболический подъём)
			t := normalizedDepth / peakPoint // 0..1 до точки пика
			verticalOffset = c.config.VerticalCurveStrength * t * (2.0 - t) * c.screenH
		} else {
			// Фаза опускания: от пика до конца
			t := (normalizedDepth - peakPoint) / (1.0 - peakPoint) // 0..1 после пика
			// Плавное опускание с усилением к концу
			dropAmount := c.config.VerticalCurveDrop * t * t
			riseAmount := c.config.VerticalCurveStrength * (1.0 - t)
			verticalOffset = (riseAmount - dropAmount) * c.screenH
		}
	}

	screenX := c.screenW/2 + relX*scale + curveOffset
	screenY := c.horizonY - relY*scale - verticalOffset

	return screenX, screenY, scale
}
