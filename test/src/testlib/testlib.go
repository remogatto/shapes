package testlib

import (
	"fmt"
	"image"
	"runtime"
	"time"

	"git.tideland.biz/goas/loop"
	"github.com/remogatto/mandala"
	gl "github.com/remogatto/opengles2"
	"github.com/remogatto/prettytest"
)

const (
	// We don't need high framerate for testing
	FramesPerSecond = 15
)

type TestSuite struct {
	prettytest.Suite

	rlControl *renderLoopControl
	timeout   <-chan time.Time

	testDraw chan image.Image

	renderState *renderState
}

type renderLoopControl struct {
	window   chan mandala.Window
	drawFunc chan func()
}

type renderState struct {
	window mandala.Window
}

func (renderState *renderState) init(window mandala.Window) {
	window.MakeContextCurrent()

	renderState.window = window
	width, height := window.GetSize()

	// Set the viewport
	gl.Viewport(0, 0, gl.Sizei(width), gl.Sizei(height))
	gl.ClearColor(0.0, 0.0, 0.0, 1.0)
}

func newRenderLoopControl() *renderLoopControl {
	return &renderLoopControl{
		drawFunc: make(chan func()),
		window:   make(chan mandala.Window),
	}
}

// Timeout timeouts the tests after the given duration.
func (t *TestSuite) Timeout(timeout time.Duration) {
	t.timeout = time.After(timeout)
}

// Run runs renderLoop. The loop renders a frame and swaps the buffer
// at each tick received.
func (t *TestSuite) renderLoopFunc(control *renderLoopControl) loop.LoopFunc {
	return func(loop loop.Loop) error {

		// renderState stores rendering state variables such
		// as the EGL state
		t.renderState = new(renderState)

		// Lock/unlock the loop to the current OS thread. This is
		// necessary because OpenGL functions should be called from
		// the same thread.
		runtime.LockOSThread()
		defer runtime.UnlockOSThread()

		window := <-control.window
		t.renderState.init(window)

		for {
			select {
			case drawFunc := <-control.drawFunc:
				drawFunc()
			}
		}
	}
}

// eventLoopFunc is listening for events originating from the
// framwork.
func (t *TestSuite) eventLoopFunc(renderLoopControl *renderLoopControl) loop.LoopFunc {
	return func(loop loop.Loop) error {

		for {
			select {

			// Receive events from the framework.
			case untypedEvent := <-mandala.Events():

				switch event := untypedEvent.(type) {

				case mandala.CreateEvent:

				case mandala.StartEvent:

				case mandala.NativeWindowCreatedEvent:
					renderLoopControl.window <- event.Window

				case mandala.ActionUpDownEvent:

				case mandala.ActionMoveEvent:

				case mandala.NativeWindowDestroyedEvent:

				case mandala.DestroyEvent:

				case mandala.NativeWindowRedrawNeededEvent:

				case mandala.PauseEvent:

				case mandala.ResumeEvent:

				}
			}
		}
	}
}

func (t *TestSuite) timeoutLoopFunc() loop.LoopFunc {
	return func(loop loop.Loop) error {
		time := <-t.timeout
		err := fmt.Errorf("Tests timed out after %v", time)
		mandala.Logf("%s %s", err.Error(), mandala.Stacktrace())
		t.Error(err)
		return nil
	}
}

func (t *TestSuite) BeforeAll() {
	// Create rendering loop control channels
	t.rlControl = newRenderLoopControl()
	// Start the rendering loop
	loop.GoRecoverable(
		t.renderLoopFunc(t.rlControl),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				mandala.Logf("%s", r.Reason)
				mandala.Logf("%s", mandala.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)
	// Start the event loop
	loop.GoRecoverable(
		t.eventLoopFunc(t.rlControl),
		func(rs loop.Recoverings) (loop.Recoverings, error) {
			for _, r := range rs {
				mandala.Logf("%s", r.Reason)
				mandala.Logf("%s", mandala.Stacktrace())
			}
			return rs, fmt.Errorf("Unrecoverable loop\n")
		},
	)

	if t.timeout != nil {
		// Start the timeout loop
		loop.GoRecoverable(
			t.timeoutLoopFunc(),
			func(rs loop.Recoverings) (loop.Recoverings, error) {
				for _, r := range rs {
					mandala.Logf("%s", r.Reason)
					mandala.Logf("%s", mandala.Stacktrace())
				}
				return rs, fmt.Errorf("Unrecoverable loop\n")
			},
		)
	}

}

func NewTestSuite() *TestSuite {
	return &TestSuite{
		rlControl: newRenderLoopControl(),
		testDraw:  make(chan image.Image),
	}
}
