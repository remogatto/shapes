package testlib

import (
	"fmt"
	"image/color"

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

func (t *TestSuite) TestShape() {
	box := shapes.NewBox(10, 20)

	// Color

	c := box.GetColor()
	nc := box.GetNColor()
	t.Equal(color.RGBA{0, 0, 255, 255}, c)
	t.Equal([4]float32{0, 0, 1, 1}, nc)

	box.Color(color.RGBA{0xaa, 0xaa, 0xaa, 0xff})
	c = box.GetColor()
	nc = box.GetNColor()
	t.Equal(color.RGBA{170, 170, 170, 255}, c)
	t.Equal([4]float32{0.6666667, 0.6666667, 0.6666667, 1}, nc)

	// GetSize

	w, h := box.GetSize()
	t.True(w == 10)
	t.True(h == 20)

	// Center

	x, y := box.Center()
	t.Equal(float32(0), x)
	t.Equal(float32(0), y)

	// Center after translation

	winW, winH := t.renderState.window.GetSize()
	world := newWorld(winW, winH)
	box.AttachToWorld(world)
	box.Position(10, 20)

	x, y = box.Center()
	t.Equal(float32(10), x)
	t.Equal(float32(20), y)

	// Angle after rotation

	box.Rotate(10)
	angle := box.Angle()
	t.Equal(float32(10), angle)

	// String representation

	t.Equal("(10.000000,20.000000)-(10.000000,20.000000)", box.String())
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
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
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
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
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
		box.Position(111, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
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
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestScaledBox() {
	filename := "expected_box_scaled.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		box := shapes.NewBox(100, 100)
		// Color is yellow
		box.Color(color.RGBA{0, 0, 255, 255})
		box.AttachToWorld(world)
		box.Position(float32(w/2), 0)
		box.Scale(1.5, 1.5)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < distanceThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestSegment() {
	filename := "expected_line.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		segment := shapes.NewSegment(81.5, -40, 238.5, 44)

		// Color is yellow
		segment.Color(color.RGBA{255, 0, 0, 255})
		segment.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		segment.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw)
	if err != nil {
		panic(err)
	}
	t.True(distance < 0.0009, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestSegmentCenter() {
	segment := shapes.NewSegment(10, 15, 20, 20)

	x, y := segment.Center()
	t.Equal(float32(15), x)
	t.Equal(float32(17.5), y)

	w, h := segment.GetSize()
	t.Equal(float32(10), w)
	t.Equal(float32(5), h)
}
