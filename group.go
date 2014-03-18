package shapes

import (
	"fmt"
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

	children map[string]Shape
}

func NewGroup() *Group {
	return &Group{
		children: make(map[string]Shape),
	}
}

func (g *Group) Add(k string, s Shape) error {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()

	if _, exists := g.children[k]; exists {
		return fmt.Errorf("a shape '%s' already exists", k)
	}

	g.children[k] = s

	// TODO recalculate center and bounds
	if len(g.children) == 1 {
		g.bounds = s.Bounds()
	} else {
		g.bounds = g.bounds.Union(s.Bounds())
	}

	g.x = float32((g.bounds.Min.X + g.bounds.Max.X) / 2)
	g.y = float32((g.bounds.Min.Y + g.bounds.Max.Y) / 2)

	return nil
}

func (g *Group) Remove(k string) error {
	g.rwMutex.Lock()
	defer g.rwMutex.Unlock()

	if _, exists := g.children[k]; !exists {
		return fmt.Errorf("cannot find a shape named '%s'", k)
	}

	g.children[k] = nil

	// TODO recalculate center and bounds

	return nil
}

func (g *Group) Child(k string) (Shape, error) {
	g.rwMutex.RLock()
	defer g.rwMutex.RUnlock()

	if _, exists := g.children[k]; !exists {
		return nil, fmt.Errorf("cannot find a shape named '%s'", k)
	}

	return g.children[k], nil
}

func (g *Group) Draw() {
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

func (g *Group) MoveTo(x, y float32) {
	for _, s := range g.children {
		sx, sy := s.Center()
		nx := x - (g.x - sx)
		ny := y - (g.y - sy)
		s.MoveTo(nx, ny)
	}
}

func (g *Group) Vertices() []float32 {
	v := []float32{}

	for _, s := range g.children {
		v = append(v, s.Vertices()...)
	}

	return v
}

func (g *Group) Center() (float32, float32) {
	return g.x, g.y
}

func (g *Group) Angle() float32 {
	return g.angle
}

func (g *Group) Bounds() image.Rectangle {
	return g.bounds
}

func (g *Group) String() string {
	str := ""

	for _, s := range g.children {
		str += s.String() + "\n"
	}

	return str
}

func (g *Group) AttachToWorld(world World) {
	for _, s := range g.children {
		s.AttachToWorld(world)
	}
}
