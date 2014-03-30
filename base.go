package shapes

import (
	"fmt"
	"image"
	"image/color"

	"github.com/remogatto/mathgl"
	"github.com/remogatto/shaders"
)

// Base shape.
type Base struct {
	// Vertices of the generic shape
	vertices []float32

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// Center
	x, y float32

	// Angle
	angle float32

	// Bounds
	bounds image.Rectangle
	// w, h float32

	// Colors
	color color.Color
	// normalized RGBA color
	nColor [4]float32
	// color matrix (four color component for each vertex)
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

// Rotate the shape around its center, by the given angle in degrees.
func (b *Base) Rotate(angle float32) {
	// axis := mathgl.Vec3f{0, 0, 0}
	// b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.HomogRotate3D(angle, axis))
	b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	b.angle = angle
}

// Scale the shape relative to its center, by the given factors.
func (b *Base) Scale(sx, sy float32) {
	b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.Scale3D(sx, sy, 1.0))

}

// Move moves the shape by dx, dy.
func (b *Base) Move(dx, dy float32) {
	b.modelMatrix = b.modelMatrix.Mul4(mathgl.Translate3D(dx, dy, 0))
	b.bounds = b.bounds.Add(image.Point{int(dx), int(dy)})
	b.x, b.y = b.x+dx, b.y+dy
}

// Move the shape by x, y
func (b *Base) MoveTo(x, y float32) {
	b.modelMatrix = mathgl.Translate3D(x, y, 0)
	dx := x - b.x
	dy := y - b.y
	b.x, b.y = x, y
	b.bounds = b.bounds.Add(image.Point{int(dx), int(dy)})
}

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

func (b *Base) Bounds() image.Rectangle {
	// return b.w, b.h
	return b.bounds
}

// Get the color of the shape.
func (b *Base) Color() color.Color {
	return b.color
}

// Get the color as a normalized float32 array.
func (b *Base) NColor() [4]float32 {
	return b.nColor
}

// Set the color of the shape.
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

// AttachToWorld fills projection and view matrices.
func (b *Base) AttachToWorld(world World) {
	b.projMatrix = world.Projection()
	b.viewMatrix = world.View()
}

// SetTexture sets a texture for the shape.
func (b *Base) SetTexture(texture uint32, texCoords []float32) error {
	b.texCoords = texCoords
	b.texBuffer = texture
	return nil
}

// String return a string representation of the shape in the form
// "(cx, cy), (w, h)".
func (b *Base) String() string {
	size := b.bounds.Size()
	return fmt.Sprintf(
		"(%f,%f)-(%f,%f)\nVertices: %v\nProjection Matrix: %v\nView Matrix: %v\nModel matrix: %v\nTexture coordinates: %v\nTexture buffer: %v\nVertex color: %v %v",
		b.x,
		b.y,
		float32(size.X),
		float32(size.Y),
		b.vertices,
		b.projMatrix,
		b.viewMatrix,
		b.modelMatrix,
		b.texCoords,
		b.texBuffer,
		b.vColor,
		b.nColor,
	)
}
