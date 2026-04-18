package core

import (
	"math"
	"testing"
)

func TestVec3Add(t *testing.T) {
	v1 := Vec3{X: 1, Y: 2, Z: 3}
	v2 := Vec3{X: 4, Y: 5, Z: 6}
	result := v1.Add(v2)

	if result.X != 5 || result.Y != 7 || result.Z != 9 {
		t.Errorf("Expected (5,7,9), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

func TestVec3Sub(t *testing.T) {
	v1 := Vec3{X: 5, Y: 7, Z: 9}
	v2 := Vec3{X: 1, Y: 2, Z: 3}
	result := v1.Sub(v2)

	if result.X != 4 || result.Y != 5 || result.Z != 6 {
		t.Errorf("Expected (4,5,6), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

func TestVec3Scale(t *testing.T) {
	v := Vec3{X: 1, Y: 2, Z: 3}
	result := v.Scale(2.0)

	if result.X != 2 || result.Y != 4 || result.Z != 6 {
		t.Errorf("Expected (2,4,6), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

func TestVec3Length(t *testing.T) {
	v := Vec3{X: 3, Y: 4, Z: 0}
	length := v.Length()

	if length != 5.0 {
		t.Errorf("Expected 5.0, got %f", length)
	}
}

func TestVec3Normalize(t *testing.T) {
	v := Vec3{X: 3, Y: 4, Z: 0}
	normalized := v.Normalize()

	expectedLength := 1.0
	actualLength := normalized.Length()

	if math.Abs(actualLength-expectedLength) > 1e-10 {
		t.Errorf("Expected normalized length 1.0, got %f", actualLength)
	}

	// Check direction is preserved
	if math.Abs(normalized.X-0.6) > 1e-10 || math.Abs(normalized.Y-0.8) > 1e-10 {
		t.Errorf("Expected (0.6,0.8,0), got (%f,%f,%f)", normalized.X, normalized.Y, normalized.Z)
	}
}

func TestVec3NormalizeZeroVector(t *testing.T) {
	v := Vec3{X: 0, Y: 0, Z: 0}
	normalized := v.Normalize()

	if normalized.X != 0 || normalized.Y != 0 || normalized.Z != 0 {
		t.Errorf("Expected zero vector (0,0,0), got (%f,%f,%f)", normalized.X, normalized.Y, normalized.Z)
	}
}

func TestVec3Dot(t *testing.T) {
	v1 := Vec3{X: 1, Y: 2, Z: 3}
	v2 := Vec3{X: 4, Y: 5, Z: 6}
	result := v1.Dot(v2)

	expected := 1*4 + 2*5 + 3*6 // 4 + 10 + 18 = 32
	if result != float64(expected) {
		t.Errorf("Expected %f, got %f", float64(expected), result)
	}
}

func TestVec3Cross(t *testing.T) {
	v1 := Vec3{X: 1, Y: 0, Z: 0}
	v2 := Vec3{X: 0, Y: 1, Z: 0}
	result := v1.Cross(v2)

	// i × j = k
	if result.X != 0 || result.Y != 0 || result.Z != 1 {
		t.Errorf("Expected (0,0,1), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

func TestVec3MulScalar(t *testing.T) {
	v := Vec3{X: 2, Y: 3, Z: 4}
	result := v.MulScalar(3.0)

	if result.X != 6 || result.Y != 9 || result.Z != 12 {
		t.Errorf("Expected (6,9,12), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

func TestVec3DivScalar(t *testing.T) {
	v := Vec3{X: 6, Y: 9, Z: 12}
	result := v.DivScalar(3.0)

	if result.X != 2 || result.Y != 3 || result.Z != 4 {
		t.Errorf("Expected (2,3,4), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

// Edge case: negative values
func TestVec3NegativeValues(t *testing.T) {
	v1 := Vec3{X: -1, Y: -2, Z: -3}
	v2 := Vec3{X: 1, Y: 2, Z: 3}

	// Add with negative
	result := v1.Add(v2)
	if result.X != 0 || result.Y != 0 || result.Z != 0 {
		t.Errorf("Expected (0,0,0), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}

	// Scale with negative
	scaled := v2.Scale(-1)
	if scaled.X != -1 || scaled.Y != -2 || scaled.Z != -3 {
		t.Errorf("Expected (-1,-2,-3), got (%f,%f,%f)", scaled.X, scaled.Y, scaled.Z)
	}
}

// Edge case: zero scalar multiplication
func TestVec3ScaleByZero(t *testing.T) {
	v := Vec3{X: 5, Y: 10, Z: 15}
	result := v.Scale(0)

	if result.X != 0 || result.Y != 0 || result.Z != 0 {
		t.Errorf("Expected (0,0,0), got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}

// Edge case: very small values (near zero)
func TestVec3SmallValues(t *testing.T) {
	v := Vec3{X: 1e-10, Y: 1e-10, Z: 1e-10}
	length := v.Length()

	if length < 0 {
		t.Errorf("Length should be non-negative, got %f", length)
	}
}

// Edge case: large values
func TestVec3LargeValues(t *testing.T) {
	v := Vec3{X: 1e10, Y: 1e10, Z: 1e10}
	normalized := v.Normalize()

	expectedLength := 1.0
	actualLength := normalized.Length()

	if math.Abs(actualLength-expectedLength) > 1e-9 {
		t.Errorf("Expected normalized length 1.0, got %f", actualLength)
	}
}

// Edge case: perpendicular vectors dot product
func TestVec3DotPerpendicular(t *testing.T) {
	v1 := Vec3{X: 1, Y: 0, Z: 0}
	v2 := Vec3{X: 0, Y: 1, Z: 0}
	result := v1.Dot(v2)

	if result != 0 {
		t.Errorf("Expected 0 for perpendicular vectors, got %f", result)
	}
}

// Edge case: parallel vectors cross product
func TestVec3CrossParallel(t *testing.T) {
	v1 := Vec3{X: 1, Y: 2, Z: 3}
	v2 := Vec3{X: 2, Y: 4, Z: 6}
	result := v1.Cross(v2)

	// Parallel vectors should have zero cross product
	if result.X != 0 || result.Y != 0 || result.Z != 0 {
		t.Errorf("Expected (0,0,0) for parallel vectors, got (%f,%f,%f)", result.X, result.Y, result.Z)
	}
}
