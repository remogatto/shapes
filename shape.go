package shapes

import "image"

type Shape interface {
	// Rotate shape
	Rotate(angle float32)

	// Scale shape
	Scale(sx, sy float32)

	// Move the shape by dx, dy
	Move(dx, dy float32)

	// Move the (center of the) shape to x, y
	MoveTo(x, y float32)

	// Renders
	Draw()

	// Returns vertices
	Vertices() []float32

	// Returns the center coords
	Center() (float32, float32)

	// Returns the rotation angle
	Angle() float32

	// Returns the bounds (width, height)
	Bounds() image.Rectangle

	// String representation
	String() string

	// Puts the shape into the world
	AttachToWorld(world World)

	// Clone clones the current shape and returns a new shape
	Clone() Shape

	// Adds a texture
	// AttachTexture(buf []byte, w, h int, texCoords []float32) error
}
