package testlib

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"path/filepath"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
	"github.com/remogatto/mandala/test/src/testlib"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shapes"
)

const (
	distanceThreshold = 0.002
)

func distanceError(distance float64, filename string) string {
	return fmt.Sprintf("Image differs by distance %f, result saved in %s", distance, filename)
}

// Compare the result of rendering against the saved expected image.
func testImage(filename string, act image.Image) (float64, image.Image, image.Image, error) {
	request := mandala.LoadAssetRequest{
		Filename: filepath.Join(expectedImgPath, filename),
		Response: make(chan mandala.LoadAssetResponse),
	}

	mandala.AssetManager() <- request
	response := <-request.Response
	buffer := response.Buffer

	if response.Error != nil {
		return 1, nil, nil, response.Error
	}

	exp, err := png.Decode(buffer)
	if err != nil {
		return 1, nil, nil, err
	}
	return imagetest.CompareDistance(exp, act), exp, act, nil
}

func (t *TestSuite) TestBox() {
	filename := "expected_box.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a box
		box := shapes.NewBox(100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Position(float32(w/2), 0)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, outputPath))
	if t.Failed() {
		saveExpAct("failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestRotatedBox() {
	filename := "expected_box_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a 100x100 pixel² box
		box := shapes.NewBox(100, 100)
		box.AttachToWorld(world)
		// Place the box at the center of the screen
		box.Position(float32(w/2), 0)
		// Rotate the box 20 degrees
		box.Rotate(20.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, outputPath))
	if t.Failed() {
		saveExpAct("failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestTranslatedBox() {
	filename := "expected_box_translated_10_10.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Place a box on the center of the window
		box := shapes.NewBox(100, 100)
		box.AttachToWorld(world)
		box.Position(100, -140)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, outputPath))
	if t.Failed() {
		saveExpAct("failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestColoredBox() {
	filename := "expected_box_yellow.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		box := shapes.NewBox(100, 100)
		// Color is yellow
		box.Color(color.RGBA{255, 255, 0, 255})
		box.AttachToWorld(world)
		box.Position(float32(w/2), 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, outputPath))
	if t.Failed() {
		saveExpAct("failed_"+filename, exp, act)
	}
}
