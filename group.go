package shapes

import (
	"image"
	"sync"
)

// Group implements Shape too.
type Group struct {
	// Center
	x, y float32

	// Angle
	angle float32

	// Bounds
	bounds image.Rectangle

	// rwMutex handle councurrent access to children slice
	rwMutex sync.RWMutex

	children []Shape
}

func NewGroup() *Group {
	return &Group{
		children: make([]Shape, 0),
	}
}

func (g *Group) Append(s Shape) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()

	g.children = append(g.children, s)

	// TODO recalculate center and bounds
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

func (g *Group) Draw() {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()
	for _, s := range g.children {
		s.Draw()
	}
}

func (g *Group) Rotate(angle float32) {
	// cx, cy := g.Center()
	for _, s := range g.children {
		s.Rotate(angle)
	}
}

func (g *Group) Scale(sx, sy float32) {
	for _, s := range g.children {
		s.Scale(sx, sy)
	}
}

func (g *Group) Move(dx, dy float32) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()
	for _, s := range g.children {
		s.Move(dx, dy)
	}
	g.x, g.y = g.x+dx, g.y+dy
}

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

func (g *Group) Center() (float32, float32) {
	return g.x, g.y
}

func (g *Group) SetCenter(x, y float32) {
	g.x, g.y = x, y
}

func (g *Group) Angle() float32 {
	return g.angle
}

func (g *Group) Bounds() image.Rectangle {
	return g.bounds
}

func (g *Group) String() string {
	str := ""

	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	for _, s := range g.children {
		str += s.String() + "\n"
	}

	return str
}

func (g *Group) AttachToWorld(world World) {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()

	for _, s := range g.children {
		s.AttachToWorld(world)
	}
}

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
