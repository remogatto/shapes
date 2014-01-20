package shapes

import (
	"fmt"
	"image/color"

	"github.com/remogatto/mathgl"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shaders"
)

type World interface {
	// Projection returns the projection matrix used to render
	// the objects in the World.
	Projection() mathgl.Mat4f

	// View returns the view matrix used to render the World from
	// the point-of-view of a camera.
	View() mathgl.Mat4f
}

type Box struct {
	// Center of the box
	X, Y float32

	// Size of the box
	Width, Height float32

	// Color of the box
	color [16]float32

	// Vertices of the box
	vertices [8]float32

	program shaders.Program

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// GLSL variables IDs
	colorId       uint32
	posId         uint32
	projMatrixId  uint32
	modelMatrixId uint32
	viewMatrixId  uint32
}

func NewBox(width, height float32) *Box {

	// The box is built around its center
	vertices := [8]float32{
		-width / 2, -height / 2,
		width / 2, -height / 2,
		-width / 2, height / 2,
		width / 2, height / 2,
	}

	// Default color is blue
	color := [16]float32{
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
		0, 0, 1, 1,
	}

	x, y := width/2, height/2

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
	program.Use()

	// Get variables IDs from shaders
	posId := program.GetAttribute("pos")
	colorId := program.GetAttribute("color")
	projMatrixId := program.GetUniform("projection")
	modelMatrixId := program.GetUniform("model")
	viewMatrixId := program.GetUniform("view")

	// Fill the model matrix with the identity.
	modelMatrix := mathgl.Ident4f()

	return &Box{
		X:             x,
		Y:             y,
		Width:         width,
		Height:        height,
		vertices:      vertices,
		color:         color,
		program:       program,
		posId:         posId,
		colorId:       colorId,
		projMatrixId:  projMatrixId,
		modelMatrixId: modelMatrixId,
		viewMatrixId:  viewMatrixId,
		modelMatrix:   modelMatrix,
	}
}

// AttachToWorld fills projection and view matrices.
func (box *Box) AttachToWorld(world World) {
	box.projMatrix = world.Projection()
	box.viewMatrix = world.View()
}

// Rotate the box around its center, by the given angle (in degrees).
func (box *Box) Rotate(angle float32) {
	box.modelMatrix = mathgl.Translate3D(box.X, box.Y, 0).Mul4(mathgl.HomogRotate3DZ(angle))
}

// Place the box at the given position
func (box *Box) Position(x, y float32) {
	box.X, box.Y = x, y
	box.modelMatrix = mathgl.Translate3D(x, y, 0)
}

// Draw actually renders the object on the surface.
func (box *Box) Draw() {
	box.program.Use()
	gl.VertexAttribPointer(box.posId, 2, gl.FLOAT, false, 0, &box.vertices[0])
	gl.EnableVertexAttribArray(box.posId)

	gl.VertexAttribPointer(box.colorId, 4, gl.FLOAT, false, 0, &box.color[0])
	gl.EnableVertexAttribArray(box.colorId)

	gl.UniformMatrix4fv(int32(box.modelMatrixId), 1, false, (*float32)(&box.modelMatrix[0]))
	gl.UniformMatrix4fv(int32(box.projMatrixId), 1, false, (*float32)(&box.projMatrix[0]))
	gl.UniformMatrix4fv(int32(box.viewMatrixId), 1, false, (*float32)(&box.viewMatrix[0]))

	gl.DrawArrays(gl.TRIANGLE_STRIP, 0, 4)

	gl.Flush()
	gl.Finish()
}

// Set the color of the shape.
func (box *Box) Color(c color.Color) {
	// Convert to RGBA
	rgba := color.RGBAModel.Convert(c)
	r, g, b, a := rgba.RGBA()
	// Normalize the color components
	box.color = [16]float32{
		float32(r / 65535), float32(g / 65535), float32(b / 65535), float32(a / 65535),
		float32(r / 65535), float32(g / 65535), float32(b / 65535), float32(a / 65535),
		float32(r / 65535), float32(g / 65535), float32(b / 65535), float32(a / 65535),
		float32(r / 65535), float32(g / 65535), float32(b / 65535), float32(a / 65535),
	}
}

// String return a string representation of the original box vertices
// (before transformation).
func (box *Box) String() string {
	return fmt.Sprintf("%v", box.vertices)
}

type Line struct {
	// Points of the line
	X1, Y1, X2, Y2 float32

	// Color of the line
	color [8]float32

	program shaders.Program

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// GLSL variables IDs
	colorId       uint32
	posId         uint32
	projMatrixId  uint32
	modelMatrixId uint32
	viewMatrixId  uint32
}

func NewLine(x1, y1, x2, y2 float32) *Line {

	// Default color is blue
	color := [8]float32{
		0, 0, 1, 1,
		0, 0, 1, 1,
	}

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
	program.Use()

	// Get variables IDs from shaders
	posId := program.GetAttribute("pos")
	colorId := program.GetAttribute("color")
	projMatrixId := program.GetUniform("projection")
	modelMatrixId := program.GetUniform("model")
	viewMatrixId := program.GetUniform("view")

	// Fill the model matrix with the identity.
	modelMatrix := mathgl.Ident4f()

	return &Line{
		X1:            x1,
		Y1:            y1,
		X2:            x2,
		Y2:            y2,
		color:         color,
		program:       program,
		posId:         posId,
		colorId:       colorId,
		projMatrixId:  projMatrixId,
		modelMatrixId: modelMatrixId,
		viewMatrixId:  viewMatrixId,
		modelMatrix:   modelMatrix,
	}
}

// AttachToWorld fills projection and view matrices.
func (line *Line) AttachToWorld(world World) {
	line.projMatrix = world.Projection()
	line.viewMatrix = world.View()
}

// Draw actually renders the object on the surface.
func (line *Line) Draw() {

	vertices := [4]float32{
		line.X1, line.Y1,
		line.X2, line.Y2,
	}

	line.program.Use()
	gl.VertexAttribPointer(line.posId, 2, gl.FLOAT, false, 0, &vertices[0])
	gl.EnableVertexAttribArray(line.posId)

	gl.VertexAttribPointer(line.colorId, 4, gl.FLOAT, false, 0, &line.color[0])
	gl.EnableVertexAttribArray(line.colorId)

	gl.UniformMatrix4fv(int32(line.modelMatrixId), 1, false, (*float32)(&line.modelMatrix[0]))
	gl.UniformMatrix4fv(int32(line.projMatrixId), 1, false, (*float32)(&line.projMatrix[0]))
	gl.UniformMatrix4fv(int32(line.viewMatrixId), 1, false, (*float32)(&line.viewMatrix[0]))

	gl.DrawArrays(gl.LINES, 0, 2)

	gl.Flush()
	gl.Finish()
}

// Set the color of the shape.
func (line *Line) Color(c color.Color) {
	// Convert to RGBA
	rgba := color.RGBAModel.Convert(c)
	r, g, b, a := rgba.RGBA()
	// Normalize the color components
	line.color = [8]float32{
		float32(r / 65535), float32(g / 65535), float32(b / 65535), float32(a / 65535),
		float32(r / 65535), float32(g / 65535), float32(b / 65535), float32(a / 65535),
	}
}

// String return a string representation of the original box vertices
// (before transformation).
func (line *Line) String() string {
	return fmt.Sprintf("x1: %f y1: %f x2: %f y2: %f", line.X1, line.Y1, line.X2, line.Y2)
}
