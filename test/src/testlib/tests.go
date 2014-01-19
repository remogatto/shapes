package testlib

import (
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"

	"github.com/remogatto/imagetest"
	"github.com/remogatto/mandala"
	"github.com/remogatto/mandala/test/src/testlib"
	"github.com/remogatto/mathgl"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/shapes"
)

const (
	expectedImgPath = "res/drawable"
	outputPath      = "output"
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

// Create an image containing both expected and actual images, side by
// side.
func saveExpAct(filename string, exp image.Image, act image.Image) {
	rect := exp.Bounds()
	dstRect := image.Rect(rect.Min.X, rect.Min.Y, rect.Max.X*2, rect.Max.Y)
	dstImage := image.NewRGBA(dstRect)

	for y := rect.Min.Y; y < rect.Max.Y; y++ {
		for x := rect.Min.X; x < rect.Max.X; x++ {
			dstImage.Set(x, y, exp.At(x, y))
			dstImage.Set(x+rect.Max.X, y, act.At(x, y))
			// Draw a white vertical line to split the
			// images
			if x == rect.Max.X-1 {
				dstImage.Set(x, y, color.White)
			}
		}
	}

	file, _ := os.Create(filepath.Join(outputPath, filename))
	png.Encode(file, dstImage)
	file.Close()
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
		box.Rotate(-20.0)
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
