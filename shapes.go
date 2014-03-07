package shapes

import (
	"fmt"
	"image/color"
	gl "github.com/remogatto/opengles2"

	"github.com/remogatto/mathgl"
	"github.com/remogatto/shaders"
)

var (
	// The default color for shapes is blue.
	DefaultColor = color.RGBA{0, 0, 0xff, 0xff}
)

type World interface {
	// Projection returns the projection matrix used to render
	// the objects in the World.
	Projection() mathgl.Mat4f

	// View returns the view matrix used to render the World from
	// the point-of-view of a camera.
	View() mathgl.Mat4f
}

type shape struct {
	x, y          float32
	width, height float32

	angle float32
	color color.Color

	// normalized RGBA color
	nColor [4]float32

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// GLSL program
	program shaders.Program

	// GLSL variables IDs
	colorId       uint32
	posId         uint32
	projMatrixId  uint32
	modelMatrixId uint32
	viewMatrixId  uint32
	texInId       uint32

	// Other texture stuff
	texRatioId uint32
	texBuffer  uint32
	textureId  uint32

	texCoords []float32
}

func (shape *shape) GetSize() (float32, float32) {
	return shape.width, shape.height
}

// Center returns the coordinates of the transformed center of the
// shape.
func (shape *shape) Center() (float32, float32) {
	return shape.x, shape.y
}

// Angle returns the current angle of the shape in degrees.
func (shape *shape) Angle() float32 {
	return shape.angle
}

// AttachToWorld fills projection and view matrices.
func (shape *shape) AttachToWorld(world World) {
	shape.projMatrix = world.Projection()
	shape.viewMatrix = world.View()
}

// Rotate the shape around its center, by the given angle in degrees.
func (shape *shape) Rotate(angle float32) {
	shape.modelMatrix = mathgl.Translate3D(shape.x, shape.y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	shape.angle = angle
}

// Scale the shape relative to its center, by the given factors.
func (shape *shape) Scale(sx, sy float32) {
	shape.modelMatrix = mathgl.Translate3D(shape.x, shape.y, 0).Mul4(mathgl.Scale3D(sx, sy, 1.0))
}

// Place the shape at the given position
func (shape *shape) Position(x, y float32) {
	shape.modelMatrix = mathgl.Translate3D(x, y, 0)
	shape.x, shape.y = x, y
}

// Set the color of the shape.
func (shape *shape) Color(c color.Color) {

	shape.color = c

	// Convert to RGBA
	rgba := color.NRGBAModel.Convert(c).(color.NRGBA)
	r, g, b, a := rgba.R, rgba.G, rgba.B, rgba.A

	// Normalize the color components
	shape.nColor = [4]float32{
		float32(r) / 255,
		float32(g) / 255,
		float32(b) / 255,
		float32(a) / 255,
	}
}

// Get the color of the shape.
func (shape *shape) GetColor() color.Color {
	return shape.color
}

// Get the color as a normalized float32 array.
func (shape *shape) GetNColor() [4]float32 {
	return shape.nColor
}

// String return a string representation of the shape in the form
// "(cx,cy),(w,h)".
func (shape *shape) String() string {
	return fmt.Sprintf("(%f,%f)-(%f,%f)", shape.x, shape.y, shape.width, shape.height)
}

func (s *shape) AttachTexture(b []byte, w, h int, texCoords []float32) error {
	s.texCoords = texCoords

	gl.GenTextures(1, &s.texBuffer)
	gl.BindTexture(gl.TEXTURE_2D, s.texBuffer)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.NEAREST)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.NEAREST)
	gl.TexImage2D(gl.TEXTURE_2D, 0, gl.RGBA, gl.Sizei(w), gl.Sizei(h), 0, gl.RGBA, gl.UNSIGNED_BYTE, gl.Void(&b[0]))

	return nil
}
