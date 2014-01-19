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
	// Center of the box in viewport coordinates
	X, Y float32

	// Size of the box
	Width, Height float32

	// Color of the box
	Color color.Color

	vertices [8]float32
	program  shaders.Program

	// Matrices
	projMatrix  mathgl.Mat4f
	modelMatrix mathgl.Mat4f
	viewMatrix  mathgl.Mat4f

	// GLSL variables IDs
	posId         uint32
	projMatrixId  uint32
	modelMatrixId uint32
	viewMatrixId  uint32
}

func NewBox(x, y, width, height float32) *Box {
	vertices := [8]float32{
		x - width/2, y - height/2,
		x + width/2, y - height/2,
		x - width/2, y + height/2,
		x + width/2, y + height/2,
	}

	// Shader sources

	vShaderSrc := (shaders.VertexShader)(
		`precision mediump float;
                 attribute vec4 pos;
                 uniform mat4 model;
                 uniform mat4 projection;
                 uniform mat4 view;
                 void main() {
                     gl_Position = projection*model*view*pos;
                 }`)
	fShaderSrc := (shaders.FragmentShader)(
		`precision mediump float;
                 void main() {
                     gl_FragColor = vec4(0.0, 0.0, 1.0, 1.0);
                 }`)

	// Link the program
	program := shaders.NewProgram(vShaderSrc.Compile(), fShaderSrc.Compile())
	program.Use()

	// Get variables IDs from shaders
	posId := program.GetAttribute("pos")
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
		program:       program,
		posId:         posId,
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

// Draw actually renders the object on the surface.
func (box *Box) Draw() {
	box.program.Use()
	gl.VertexAttribPointer(box.posId, 2, gl.FLOAT, false, 0, &box.vertices[0])
	gl.EnableVertexAttribArray(box.posId)

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
