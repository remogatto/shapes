package shapes

type Shape interface {
	// Rotate shape
	Rotate(angle float32)

	// Scale shape
	Scale(sx, sy float32)

	// Set position of the (center of the) shape
	SetPosition(x, y float32)

	// Renders
	Draw()

	// Returns the center coords
	Center() (float32, float32)

	// Returns the rotation angle
	Angle() float32

	// Returns the bounds (width, height)
	Bounds() (float32, float32)

	// String representation
	String() string

	// Puts the shape into the world
	AttachToWorld(world World)

	// Adds a texture
	AttachTexture(buf []byte, w, h int, texCoords []float32) error
}
