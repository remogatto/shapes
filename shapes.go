package shapes

import (
	"fmt"
	"image/color"
	"math"

	"github.com/remogatto/mathgl"
	gl "github.com/remogatto/opengles2"
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
}

func (shape *shape) GetSize() (float32, float32) {
	return shape.width, shape.height
}

func (shape *shape) Center() (float32, float32) {
	return shape.x, shape.y
}

func (shape *shape) Angle() float32 {
	return shape.angle
}

// AttachToWorld fills projection and view matrices.
func (shape *shape) AttachToWorld(world World) {
	shape.projMatrix = world.Projection()
	shape.viewMatrix = world.View()
}

// Rotate the box around its center, by the given angle (in degrees).
func (shape *shape) Rotate(angle float32) {
	shape.modelMatrix = mathgl.Translate3D(shape.x, shape.y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
	shape.angle = angle
}

// Place the box at the given position
func (shape *shape) Position(x, y float32) {
	shape.modelMatrix = mathgl.Translate3D(x, y, 0)
	shape.x, shape.y = x, y
}

// Set the color of the shape.
func (shape *shape) Color(c color.Color) {
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

type Box struct {
	shape

	// 4x4 color matrix (four color component for each vertex)
	vertexColor [16]float32

	// Vertices of the box
	vertices [8]float32
}

func NewBox(width, height float32) *Box {

	box := new(Box)

	// The box is built around its center
	box.vertices = [8]float32{
		-width / 2, -height / 2,
		width / 2, -height / 2,
		-width / 2, height / 2,
		width / 2, height / 2,
	}

	// Set the default color
	box.Color(DefaultColor)

	// Shader sources

	vShaderSrc := (shaders.VertexShader)(
		`precision mediump float;
                 attribute vec4 pos;
                 attribute vec4 color;
                 varying vec4 vColor;
                 uniform mat4 model;
                 uniform mat4 projection;
                 uniform mat4 view;
                 void main() {
                     gl_Position = projection*model*view*pos;
                     vColor = color;
                 }`)
	fShaderSrc := (shaders.FragmentShader)(
		`precision mediump float;
                 varying vec4 vColor;
                 void main() {
                     gl_FragColor = vColor;
                 }`)

	// Link the program
	program := shaders.NewProgram(vShaderSrc.Compile(), fShaderSrc.Compile())
	box.program = program
	box.program.Use()

	// Get variables IDs from shaders
	box.posId = program.GetAttribute("pos")
	box.colorId = program.GetAttribute("color")
	box.projMatrixId = program.GetUniform("projection")
	box.modelMatrixId = program.GetUniform("model")
	box.viewMatrixId = program.GetUniform("view")

	// Fill the model matrix with the identity.
	box.modelMatrix = mathgl.Ident4f()

	// Size of the box
	box.width = width
	box.height = height

	// Center of the box
	box.x = width / 2
	box.y = height / 2

	return box
}

// Draw actually renders the object on the surface.
func (box *Box) Draw() {

	// Color is the same for each vertex
	vertexColor := [16]float32{
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
		box.nColor[0], box.nColor[1], box.nColor[2], box.nColor[3],
	}

	box.program.Use()

	gl.VertexAttribPointer(box.posId, 2, gl.FLOAT, false, 0, &box.vertices[0])
	gl.EnableVertexAttribArray(box.posId)

	gl.VertexAttribPointer(box.colorId, 4, gl.FLOAT, false, 0, &vertexColor[0])
	gl.EnableVertexAttribArray(box.colorId)

	gl.UniformMatrix4fv(int32(box.modelMatrixId), 1, false, (*float32)(&box.modelMatrix[0]))
	gl.UniformMatrix4fv(int32(box.projMatrixId), 1, false, (*float32)(&box.projMatrix[0]))
	gl.UniformMatrix4fv(int32(box.viewMatrixId), 1, false, (*float32)(&box.viewMatrix[0]))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.Flush()
	gl.Finish()
}

// String return a string representation of the original box vertices
// (before transformation).
func (box *Box) String() string {
	return fmt.Sprintf("%v", box.vertices)
}

type Line struct {
	shape

	// Points of the line
	x1, y1, x2, y2 float32
	vertices       [4]float32
}

func NewLine(x1, y1, x2, y2 float32) *Line {

	line := new(Line)

	// Set the default color

	line.Color(DefaultColor)

	// Set the geometry

	line.x1, line.x2 = x1, x2
	line.y1, line.y2 = y1, y2

	line.vertices = [4]float32{
		line.x1, line.y1,
		line.x2, line.y2,
	}

	// Size of the line bounding box

	line.width = float32(math.Abs(float64(x1 - x2)))
	line.height = float32(math.Abs(float64(y1 - y2)))

	// Center of the line
	line.x = (line.x1 + line.x2) / 2
	line.y = (line.y1 + line.y2) / 2

	// Shader sources

	vShaderSrc := (shaders.VertexShader)(
		`precision mediump float;
                 attribute vec4 pos;
                 attribute vec4 color;
                 varying vec4 vColor;
                 uniform mat4 model;
                 uniform mat4 projection;
                 uniform mat4 view;
                 void main() {
                     gl_Position = projection*model*view*pos;
                     vColor = color;
                 }`)
	fShaderSrc := (shaders.FragmentShader)(
		`precision mediump float;
                 varying vec4 vColor;
                 void main() {
                     gl_FragColor = vColor;
                 }`)

	// Link the program
	line.program = shaders.NewProgram(vShaderSrc.Compile(), fShaderSrc.Compile())
	line.program.Use()

	// Get variables IDs from shaders
	line.posId = line.program.GetAttribute("pos")
	line.colorId = line.program.GetAttribute("color")
	line.projMatrixId = line.program.GetUniform("projection")
	line.modelMatrixId = line.program.GetUniform("model")
	line.viewMatrixId = line.program.GetUniform("view")

	// Fill the model matrix with the identity.
	line.modelMatrix = mathgl.Ident4f()

	return line
}

// Draw actually renders the object on the surface.
func (line *Line) Draw() {
	// Color is the same for each vertex
	vertexColor := [8]float32{
		line.nColor[0], line.nColor[1], line.nColor[2], line.nColor[3],
		line.nColor[0], line.nColor[1], line.nColor[2], line.nColor[3],
	}

	line.program.Use()
	gl.VertexAttribPointer(line.posId, 2, gl.FLOAT, false, 0, &line.vertices[0])
	gl.EnableVertexAttribArray(line.posId)

	gl.VertexAttribPointer(line.colorId, 4, gl.FLOAT, false, 0, &vertexColor[0])
	gl.EnableVertexAttribArray(line.colorId)

	gl.UniformMatrix4fv(int32(line.modelMatrixId), 1, false, (*float32)(&line.modelMatrix[0]))
	gl.UniformMatrix4fv(int32(line.projMatrixId), 1, false, (*float32)(&line.projMatrix[0]))
	gl.UniformMatrix4fv(int32(line.viewMatrixId), 1, false, (*float32)(&line.viewMatrix[0]))

	gl.DrawArrays(gl.LINES, 0, 2)

	gl.Flush()
	gl.Finish()
}

// String return a string representation of the original box vertices
// (before transformation).
func (line *Line) String() string {
	return fmt.Sprintf("x1: %f y1: %f x2: %f y2: %f", line.x1, line.y1, line.x2, line.y2)
}
