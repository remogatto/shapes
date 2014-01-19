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
	distanceThreshold = 0.001
)

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
		// Place a box on the center of the window
		box := shapes.NewBox(0, 0, 100, 100)
		box.Color = color.RGBA{1.0, 0.0, 0.0, 1.0}
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, fmt.Sprint("Image differs, result saved in %s", outputPath))
	if t.Failed() {
		saveExpAct("failed_"+filename, exp, act)
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
		box.Rotate(-20.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, fmt.Sprint("Image differs, result saved in %s", outputPath))
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
		box := shapes.NewBox(0, 0, 100, 100)
		box.Color = color.RGBA{1.0, 0.0, 0.0, 1.0}
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, fmt.Sprint("Image differs, result saved in %s", outputPath))
	if t.Failed() {
		saveExpAct("failed_"+filename, exp, act)
	}
}
