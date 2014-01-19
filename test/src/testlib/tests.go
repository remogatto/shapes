package testlib

import (
	"fmt"
	"image/color"
	"image/png"
	"path/filepath"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
	"github.com/remogatto/mandala/test/src/testlib"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shapes"
)

func (t *TestSuite) TestBox() {
	filename := "expected_box.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Place a box on the center of the window
		box := shapes.NewBox(0, 0, 100, 100)
		box.Color = color.RGBA{1.0, 0.0, 0.0, 1.0}
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	request := mandala.LoadAssetRequest{
		Filename: filepath.Join(expectedImgPath, "expected_box.png"),
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	if buffer != nil {
		exp, err := png.Decode(buffer)
		t.True(err == nil, "An error occured during png decoding")
		act := <-t.testDraw
		distance := imagetest.CompareDistance(exp, act)
		t.True(distance < 0.001, fmt.Sprintf("Distance is %f", distance))
		if t.Failed() {
			saveExpAct("failed_"+filename, exp, act)
		}
	}
}

func (t *TestSuite) TestRotatedBox() {
	filename := "expected_box_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Place a box on the center of the window
		box := shapes.NewBox(0, 0, 100, 100)
		box.Color = color.RGBA{1.0, 0.0, 0.0, 1.0}
		box.AttachToWorld(world)
		// Rotate the box 20 degrees
		box.Rotate(20.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	request := mandala.LoadAssetRequest{
		Filename: filepath.Join(expectedImgPath, "expected_box_rotated_20.png"),
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	t.True(response.Error == nil, "An error occured during resource opening")

	if buffer != nil {
		exp, err := png.Decode(buffer)
		t.True(err == nil, "An error occured during png decoding")
		act := <-t.testDraw
		distance := imagetest.CompareDistance(exp, act)
		t.True(distance < 0.001, fmt.Sprintf("Distance is %f", distance))

		if t.Failed() {
			saveExpAct("failed_"+filename, exp, act)
		}
	}
}
