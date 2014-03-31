package shapes

import (
	"image"
	"sync"
)

// Group is a structure for grouping shapes. It implements Shape.
type Group struct {
	// Center of the group
	x, y float32

	// Angle
	angle float32

	// Bounds of the group
	bounds image.Rectangle

	// rwMutex handle councurrent access to children slice
	rwMutex sync.RWMutex

	// children is the slice containing the shapes of the group
	children []Shape
}

// NewGroup instantiates a group object.
func NewGroup() *Group {
	return &Group{
		children: make([]Shape, 0),
	}
}

// Append appends a shape to the group.
func (g *Group) Append(s Shape) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()

	g.children = append(g.children, s)

	if len(g.children) == 1 {
		g.bounds = s.Bounds()
	} else {
		g.bounds = g.bounds.Union(s.Bounds())
	}

	g.x = float32((g.bounds.Min.X + g.bounds.Max.X) / 2)
	g.y = float32((g.bounds.Min.Y + g.bounds.Max.Y) / 2)
}

// func (g *Group) Remove(k string) error {
// 	g.rwMutex.Lock()
// 	defer g.rwMutex.Unlock()

// 	if _, exists := g.children[k]; !exists {
// 		return fmt.Errorf("cannot find a shape named '%s'", k)
// 	}

// 	g.children[k] = nil

// 	// TODO recalculate center and bounds

// 	return nil
// }

// func (g *Group) Child(k string) (Shape, error) {
// 	g.rwMutex.RLock()
// 	defer g.rwMutex.RUnlock()

// 	if _, exists := g.children[k]; !exists {
// 		return nil, fmt.Errorf("cannot find a shape named '%s'", k)
// 	}

// 	return g.children[k], nil
// }

// Draw draws all the shapes in the group calling their Draw method.
func (g *Group) Draw() {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()
	for _, s := range g.children {
		s.Draw()
	}
}

// Rotate rotates the group aroung its center.
func (g *Group) Rotate(angle float32) {
	for _, s := range g.children {
		s.Rotate(angle)
	}
}

// Scale scales each shape of the group.
func (g *Group) Scale(sx, sy float32) {
	for _, s := range g.children {
		s.Scale(sx, sy)
	}
}

// Move moves each shape of the group by dx, dy.
func (g *Group) Move(dx, dy float32) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()
	for _, s := range g.children {
		s.Move(dx, dy)
	}
	g.x, g.y = g.x+dx, g.y+dy
}

// MoveTo moves the group in the (x,y) position.
func (g *Group) MoveTo(x, y float32) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()
	dx, dy := x-g.x, y-g.y
	for _, s := range g.children {
		s.Move(dx, dy)
	}
	g.x, g.y = x, y
}

func (g *Group) Vertices() []float32 {
	v := []float32{}

	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	for _, s := range g.children {
		v = append(v, s.Vertices()...)
	}

	return v
}

// Center returns the center of the group.
func (g *Group) Center() (float32, float32) {
	return g.x, g.y
}

// SetCenter sets a new center for the group.
func (g *Group) SetCenter(x, y float32) {
	g.x, g.y = x, y
}

// Angle returns the rotation angle of the group.
func (g *Group) Angle() float32 {
	return g.angle
}

// Bounds returns the bounding rectangle of the group.
func (g *Group) Bounds() image.Rectangle {
	return g.bounds
}

// String returns a textual representation of the group.
func (g *Group) String() string {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	str := ""
	for _, s := range g.children {
		str += s.String() + "\n"
	}
	return str
}

// AttachToWorld attaches the group to a world.
func (g *Group) AttachToWorld(world World) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()

	for _, s := range g.children {
		s.AttachToWorld(world)
	}
}

// Clone returns a copy of the group.
func (g *Group) Clone() Shape {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	cg := NewGroup()

	for _, s := range g.children {
		cs := s.Clone()
		cg.Append(cs)
	}

	return cg
}

// SetTexture sets the same texture to all shapes in the group.
func (g *Group) SetTexture(texture uint32, texCoords []float32) error {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()
	for _, s := range g.children {
		s.SetTexture(texture, texCoords)
	}
	return nil
}
