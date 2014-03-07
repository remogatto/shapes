package shapes

import (
	"github.com/remogatto/mathgl"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shaders"
)

var (
	DefaultBoxVS = (shaders.VertexShader)(
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
	DefaultBoxFS = (shaders.FragmentShader)(
		`precision mediump float;
                 varying vec4 vColor;
                 void main() {
                     gl_FragColor = vColor;
                 }`)
)

// A Box
type Box struct {
	shape

	// 4x4 color matrix (four color component for each vertex)
	vertexColor [16]float32

	// Vertices of the box
	vertices [8]float32
}

// NewBox creates a new box of given sizes.
func NewBox(program shaders.Program, width, height float32) *Box {

	box := new(Box)

	// The box is built around its center at (0, 0)
	box.vertices = [8]float32{
		-width / 2, -height / 2,
		width / 2, -height / 2,
		-width / 2, height / 2,
		width / 2, height / 2,
	}

	// Set the default color
	box.Color(DefaultColor)

	box.program = program
	box.program.Use()

	// Get variables IDs from shaders
	box.posId = box.program.GetAttribute("pos")
	box.colorId = box.program.GetAttribute("color")
	box.projMatrixId = box.program.GetUniform("projection")
	box.modelMatrixId = box.program.GetUniform("model")
	box.viewMatrixId = box.program.GetUniform("view")

	// Fill the model matrix with the identity.
	box.modelMatrix = mathgl.Ident4f()

	// Size of the box
	box.width = width
	box.height = height

	return box
}

// Draw actually renders the shape on the surface.
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
