package shapes

import (
	"fmt"
	"image/color"

	"github.com/remogatto/mathgl"
	gl "github.com/remogatto/opengles2"
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
	w, h float32

	// Colors
	color color.Color
	// normalized RGBA color
	nColor [4]float32
	// color matrix (four color component for each vertex)
	vertexColor []float32

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
	b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	b.angle = angle
}

// Scale the shape relative to its center, by the given factors.
func (b *Base) Scale(sx, sy float32) {
	b.modelMatrix = mathgl.Translate3D(b.x, b.y, 0).Mul4(mathgl.Scale3D(sx, sy, 1.0))
}

// Place the shape at the given position
func (b *Base) Move(x, y float32) {
	b.modelMatrix = mathgl.Translate3D(x, y, 0)
	b.x, b.y = x, y
}

// Center returns the coordinates of the transformed center of the
// shape.
func (b *Base) Center() (float32, float32) {
	return b.x, b.y
}

// Angle returns the current angle of the shape in degrees.
func (b *Base) Angle() float32 {
	return b.angle
}

func (b *Base) Bounds() (float32, float32) {
	return b.w, b.h
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
}

// AttachToWorld fills projection and view matrices.
func (b *Base) AttachToWorld(world World) {
	b.projMatrix = world.Projection()
	b.viewMatrix = world.View()
}

// Binds a texture to the shape.
func (b *Base) AttachTexture(buf []byte, w, h int, texCoords []float32) error {

	b.texCoords = texCoords

	gl.GenTextures(1, &b.texBuffer)
	gl.BindTexture(gl.TEXTURE_2D, b.texBuffer)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		gl.Sizei(w),
		gl.Sizei(h),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Void(&buf[0]),
	)

	return nil
}

// String return a string representation of the shape in the form
// "(cx, cy), (w, h)".
func (b *Base) String() string {
	return fmt.Sprintf("(%f,%f)-(%f,%f)", b.x, b.y, b.w, b.h)
}

// Dummy method, must be implemented by "concrete" shape
func (b *Base) Draw() {
	// TODO Change panic with something else
	panic("Cannot draw a generic shape")
}
