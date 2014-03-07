package testlib

import (
	"fmt"
	"image"
	"image/color"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala/test/src/testlib"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shapes"
)

const (
	distanceThreshold = 0.002
	texFilename       = "gopher.png"
	texDistThreshold  = 0.004
)

func distanceError(distance float64, filename string) string {
	return fmt.Sprintf("Image differs by distance %f, result saved in %s", distance, filename)
}

func (t *TestSuite) TestShape() {
	box := shapes.NewBox(t.renderState.boxProgram, 10, 20)

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
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Position(float32(w/2), 0)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
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
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
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
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
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
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		box.Position(111, 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
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
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		// Color is yellow
		box.Color(color.RGBA{255, 255, 0, 255})
		box.AttachToWorld(world)
		box.Position(float32(w/2), 0)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
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
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
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
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
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
		segment := shapes.NewSegment(t.renderState.segmentProgram, 81.5, -40, 238.5, 44)

		// Color is yellow
		segment.Color(color.RGBA{255, 0, 0, 255})
		segment.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		segment.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < 0.0009, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestSegmentCenter() {
	segment := shapes.NewSegment(t.renderState.segmentProgram, 10, 15, 20, 20)

	x, y := segment.Center()
	t.Equal(float32(15), x)
	t.Equal(float32(17.5), y)

	w, h := segment.GetSize()
	t.Equal(float32(10), w)
	t.Equal(float32(5), h)
}

func (t *TestSuite) TestTexturedBox() {
	filename := "expected_box_textured.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Position(float32(w/2), 0)

		texImg, err := loadImageResource(texFilename)
		if err != nil {
			panic(err)
		}

		buffer, texImgWidth, texImgHeight := getBufferDataFromImage(texImg)

		texCoords := []float32{
			0, 0,
			1, 0,
			0, 1,
			1, 1,
		}

		box.AttachTexture(buffer, texImgWidth, texImgHeight, texCoords)

		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < texDistThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestTexturedRotatedBox() {
	filename := "expected_box_textured_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Position(float32(w/2), 0)

		texImg, err := loadImageResource(texFilename)
		if err != nil {
			panic(err)
		}

		buffer, texImgWidth, texImgHeight := getBufferDataFromImage(texImg)

		texCoords := []float32{
			0, 0,
			1, 0,
			0, 1,
			1, 1,
		}

		box.AttachTexture(buffer, texImgWidth, texImgHeight, texCoords)

		box.Rotate(20.0)

		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < texDistThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func (t *TestSuite) TestPartialTextureRotatedBox() {
	filename := "expected_box_partial_texture_rotated_20.png"
	t.rlControl.drawFunc <- func() {
		w, h := t.renderState.window.GetSize()
		world := newWorld(w, h)
		// Create a box
		box := shapes.NewBox(t.renderState.boxProgram, 100, 100)
		box.AttachToWorld(world)
		gl.Clear(gl.COLOR_BUFFER_BIT)
		box.Position(float32(w/2), 0)

		texImg, err := loadImageResource(texFilename)
		if err != nil {
			panic(err)
		}

		buffer, texImgWidth, texImgHeight := getBufferDataFromImage(texImg)

		texCoords := []float32{
			0, 0,
			0.5, 0,
			0, 0.5,
			0.5, 0.5,
		}

		box.AttachTexture(buffer, texImgWidth, texImgHeight, texCoords)

		box.Rotate(20.0)

		box.Draw()
		t.testDraw <- testlib.Screenshot(t.renderState.window)
		t.renderState.window.SwapBuffers()
	}
	distance, exp, act, err := testlib.TestImage(filename, <-t.testDraw, imagetest.Center)
	if err != nil {
		panic(err)
	}
	t.True(distance < texDistThreshold, distanceError(distance, filename))
	if t.Failed() {
		saveExpAct(t.outputPath, "failed_"+filename, exp, act)
	}
}

func getBufferDataFromImage(img image.Image) ([]byte, int, int) {
	bounds := img.Bounds()
	imgWidth, imgHeight := bounds.Size().X, bounds.Size().Y
	buffer := make([]byte, imgWidth*imgHeight*4)
	index := 0
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := img.At(x, y).RGBA()
			buffer[index] = byte(r)
			buffer[index+1] = byte(g)
			buffer[index+2] = byte(b)
			buffer[index+3] = byte(a)
			index += 4
		}
	}

	return buffer, imgWidth, imgHeight
}
