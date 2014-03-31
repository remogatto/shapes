package shapes

import (
	"image"
	"image/color"

	"github.com/remogatto/mathgl"
	"github.com/remogatto/shaders"
)

// Base represent a basic structure for shapes.
type Base struct {
	// Vertices of the generic shape
	vertices []float32

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// Center of the shape
	x, y float32

	// Angle
	angle float32

	// Bounds
	bounds image.Rectangle

	// Color
	color color.Color

	// Normalized RGBA color
	nColor [4]float32

	// Color matrix (four color component for each vertex)
	vColor []float32

	// Texture
	texBuffer uint32
	texCoords []float32

	// GLSL program
	program shaders.Program

	// GLSL variables IDs
	colorId       uint32
	posId         uint32
	projMatrixId  uint32
	modelMatrixId uint32
	viewMatrixId  uint32
	texInId       uint32
	texRatioId    uint32
	textureId     uint32
}

// Rotate rotates the shape around its center, by the given angle in
// degrees.
func (b *Base) Rotate(angle float32) {
	b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	b.angle = angle
}

// RotateAround rotates the shape around the given point, by the given
// angle in degrees.
func (b *Base) RotateAround(x, y, angle float32) {
	dx, dy := x-b.x, y-b.y
	b.modelMatrix = mathgl.Translate3D(x, y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	b.modelMatrix = b.modelMatrix.Mul4(mathgl.Translate3D(-dx, -dy, 0))
	b.angle = angle
}

// Scale scales the shape relative to its center, by the given
// factors.
func (b *Base) Scale(sx, sy float32) {
	b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.Scale3D(sx, sy, 1.0))

}

// Move moves the shape by dx, dy.
func (b *Base) Move(dx, dy float32) {
	b.modelMatrix = b.modelMatrix.Mul4(mathgl.Translate3D(dx, dy, 0))
	b.bounds = b.bounds.Add(image.Point{int(dx), int(dy)})
	b.x, b.y = b.x+dx, b.y+dy
}

// MoveTo moves the shape in x, y position.
func (b *Base) MoveTo(x, y float32) {
	b.modelMatrix = mathgl.Translate3D(x, y, 0)
	dx := x - b.x
	dy := y - b.y
	b.x, b.y = x, y
	b.bounds = b.bounds.Add(image.Point{int(dx), int(dy)})
}

// Vertices returns the vertices slice.
func (b *Base) Vertices() []float32 {
	return b.vertices
}

// Center returns the coordinates of the transformed center of the
// shape.
func (b *Base) Center() (float32, float32) {
	return b.x, b.y
}

// SetCenter sets the coordinates of the center of the shape.
func (b *Base) SetCenter(x, y float32) {
	b.x, b.y = x, y
}

// Angle returns the current angle of the shape in degrees.
func (b *Base) Angle() float32 {
	return b.angle
}

// Bounds returns the bounds of the shape as a Rectangle.
func (b *Base) Bounds() image.Rectangle {
	return b.bounds
}

// Color returns the color of the shape.
func (b *Base) Color() color.Color {
	return b.color
}

// NColor returns the color of the shape as a normalized float32
// array.
func (b *Base) NColor() [4]float32 {
	return b.nColor
}

// SetColor sets the color of the shape.
func (s *Base) SetColor(c color.Color) {

	s.color = c

	// Convert to RGBA
	rgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	r, g, b, a := rgba.R, rgba.G, rgba.B, rgba.A

	// Normalize the color components
	s.nColor = [4]float32{
		float32(r) / 255,
		float32(g) / 255,
		float32(b) / 255,
		float32(a) / 255,
	}

	// TODO improve code
	vCount := len(s.vertices) / 2
	s.vColor = s.vColor[:0]
	for i := 0; i < vCount; i++ {
		s.vColor = append(s.vColor, s.nColor[0], s.nColor[1], s.nColor[2], s.nColor[3])
	}
}

// AttachToWorld fills projection and view matrices with world's
// matrices.
func (b *Base) AttachToWorld(world World) {
	b.projMatrix = world.Projection()
	b.viewMatrix = world.View()
}

// SetTexture sets a texture for the shape. Texture argument is an
// uint32 value returned by the OpenGL context. It's a client-code
// responsibility to provide that value.
func (b *Base) SetTexture(texture uint32, texCoords []float32) error {
	b.texCoords = texCoords
	b.texBuffer = texture
	return nil
}

// String returns a string representation of the shape.
func (b *Base) String() string {
	return b.bounds.String()
}
