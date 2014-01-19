# Shapes

<tt>shapes</tt> is a small package for drawing simple two-dimensional
shapes on an OpenGL ES 2 context.

# Example API

~~~go
// Place 100x100 pixels² box on the center of the screen
box := NewBox(0, 0, 100, 100)

// Attach the box to a world object (see World interface)
box.AttachToWorld(world)

// Render the box on the surface
box.Draw()
~~~

# Test

To run the black-box test you need to install
[Mandala](https://github.com/remogatto/mandala). Then simply use the
following commands:

<pre>
cd test
gotask run xorg # or
gotask run android
</pre>

# LICENSE

See [LICENSE](LICENSE)
