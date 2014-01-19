// +build !android

package main

import (
	"runtime"
	"testing"

	"testlib"
	glfw "github.com/go-gl/glfw3"
	"github.com/remogatto/mandala"
	"github.com/remogatto/prettytest"
)

const (
	Width  = 320
	Height = 480
)

func main() {
	runtime.LockOSThread()
	defer runtime.UnlockOSThread()
	defer glfw.Terminate()

	mandala.Verbose = true

	if !glfw.Init() {
		panic("Can't init glfw!")
	}

	// Enable OpenGL ES 2.0.
	glfw.WindowHint(glfw.ClientApi, glfw.OpenglEsApi)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	window, err := glfw.CreateWindow(Width, Height, "Shapes black-box testing", nil, nil)
	if err != nil {
		panic(err)
	}

	glfw.SwapInterval(0)

	mandala.Init(window)

	go prettytest.Run(new(testing.T), testlib.NewTestSuite())

	for !window.ShouldClose() {
		glfw.WaitEvents()
	}
}
