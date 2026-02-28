// Пакет core содержит базовые типы и утилиты для 3D-вычислений.
package core

import (
	"math"
)

// Vec3 представляет трёхмерный вектор с координатами X, Y, Z.
// Используется для задания точек, направлений, цветов в RGB и других величин.
type Vec3 struct {
	X float64
	Y float64
	Z float64
}

// Сложение векторов: возвращает новый вектор (v + other).
func (v Vec3) Add(other Vec3) Vec3 {
	return Vec3{v.X + other.X, v.Y + other.Y, v.Z + other.Z}
}

// Вычитание векторов: возвращает новый вектор (v - other).
func (v Vec3) Sub(other Vec3) Vec3 {
	return Vec3{v.X - other.X, v.Y - other.Y, v.Z - other.Z}
}

// Умножение на скаляр: возвращает новый вектор, умноженный на t.
func (v Vec3) MulScalar(t float64) Vec3 {
	return Vec3{v.X * t, v.Y * t, v.Z * t}
}

// Деление на скаляр: возвращает новый вектор, поделённый на t (проверка на ноль не выполняется).
func (v Vec3) DivScalar(t float64) Vec3 {
	return Vec3{v.X / t, v.Y / t, v.Z / t}
}

// Dot вычисляет скалярное произведение двух векторов.
func (v Vec3) Dot(other Vec3) float64 {
	return v.X*other.X + v.Y*other.Y + v.Z*other.Z
}

// Cross вычисляет векторное произведение (v × other).
func (v Vec3) Cross(other Vec3) Vec3 {
	return Vec3{
		X: v.Y*other.Z - v.Z*other.Y,
		Y: v.Z*other.X - v.X*other.Z,
		Z: v.X*other.Y - v.Y*other.X,
	}
}

// Length возвращает длину (модуль) вектора.
func (v Vec3) Length() float64 {
	return math.Sqrt(v.X*v.X + v.Y*v.Y + v.Z*v.Z)
}

// Normalize возвращает нормализованный вектор (единичной длины).
// Если длина равна нулю, возвращается нулевой вектор.
func (v Vec3) Normalize() Vec3 {
	len := v.Length()
	if len == 0 {
		return Vec3{}
	}
	return v.DivScalar(len)
}
