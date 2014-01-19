package testlib

import (
	"fmt"
	"image/color"
	"image/png"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
	"github.com/remogatto/mathgl"
	"github.com/remogatto/shapes"
)

type world struct {
	width, height int
	projMatrix    mathgl.Mat4f
	viewMatrix    mathgl.Mat4f
}

func newWorld(width, height int) *world {
	return &world{
		width:      width,
		height:     height,
		projMatrix: mathgl.Ortho2D(-float32(width/2), float32(width/2), float32(height/2), -float32(height/2)),
		viewMatrix: mathgl.Ident4f(),
	}
}

func (w *world) Projection() mathgl.Mat4f {
	return w.projMatrix
}

func (w *world) View() mathgl.Mat4f {
	return w.viewMatrix
}

func (t *TestSuite) TestBox() {
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Place a box on the center of the window
		box := shapes.NewBox(0, 0, 100, 100)
		box.Color = color.RGBA{1.0, 0.0, 0.0, 1.0}
		box.AttachToWorld(world)
		box.Draw()
	}
	request := mandala.LoadAssetRequest{
		Filename: "res/drawable/expected_box.png",
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	if buffer != nil {
		exp, err := png.Decode(buffer)
		t.True(err == nil, "An error occured during png decoding")

		distance := imagetest.CompareDistance(exp, <-t.testDraw)
		t.True(distance < 0.001, fmt.Sprintf("Distance is %f", distance))
	}
}
