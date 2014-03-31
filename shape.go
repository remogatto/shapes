package shapes

import (
	"image"
	"image/color"
)

var (
	// The default color for shapes is blue.
	DefaultColor = color.RGBA{0, 0, 0xff, 0xff}
)

// Shape is an interface describing a shape.
type Shape interface {
	// Rotate rotates the shape by angle in degree.
	Rotate(angle float32)

	// RotateAround rotates the shape around the given point, by
	// the given angle in degrees.
	RotateAround(x, y, angle float32)

	// Scale scales the shape by (sx,sy) factor.
	Scale(sx, sy float32)

	// Move moves the shape by (dx, dy).
	Move(dx, dy float32)

	// MoveTo moves the (center of the) shape in position (x,y).
	MoveTo(x, y float32)

	// Draw renders the shape on the surface.
	Draw()

	// Vertices returns the vertices slice of the shape.
	Vertices() []float32

	// Center returns the center coordinates of the shape.
	Center() (float32, float32)

	// Angle returns the rotation angle of the shape.
	Angle() float32

	// Bounds returns the bounding rectangle of the shape.
	Bounds() image.Rectangle

	// String returns a string representation of the shape.
	String() string

	// AttachToWorld attaches the shape to a world.
	AttachToWorld(world World)

	// Clone clones the current shape and returns a new shape.
	Clone() Shape

	// SetTexture sets a texture for the shape.
	SetTexture(texture uint32, texCoords []float32) error
}
