# Shapes

<tt>shapes</tt> is a small package for drawing simple two-dimensional
shapes on an OpenGL ES 2 context.

# Example API

~~~go
// Create a 100x100 pixels² box
box := NewBox(100, 100)

// Place the box at a given position
box.Position(10, 10)

// Color it
box.Color(color.White)

// Attach the box to a world object (see World interface)
box.AttachToWorld(world)

// Render the box on the surface
box.Draw()

// swap the buffers
~~~

# Test

See [test](test/) for a black-box testing approach.

# LICENSE

See [LICENSE](LICENSE)
