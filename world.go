package shapes

import "github.com/remogatto/mathgl"

// World in an interface for projection and view matrix.
type World interface {
	// Projection returns the projection matrix used to render
	// the objects in the World.
	Projection() mathgl.Mat4f

	// View returns the view matrix used to render the World from
	// the point-of-view of a camera.
	View() mathgl.Mat4f
}
