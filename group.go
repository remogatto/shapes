package shapes

import "fmt"

type children map[string]Shape

type Group struct {
	// Center
	x, y float32

	// Angle
	angle float32

	// Bounds
	w, h float32

	children children
}

func (g *Group) Add(k string, s Shape) error {
	if _, exists := g.children[k]; exists {
		return fmt.Errorf("a shape '%s' already exists", k)
	}

	g.children[k] = s

	// TODO recalculate center and bounds

	return nil
}

func (g *Group) Remove(k string) error {
	if _, exists := g.children[k]; !exists {
		return fmt.Errorf("cannot find a shape named '%s'", k)
	}

	g.children[k] = nil

	// TODO recalculate center and bounds

	return nil
}

func (g *Group) Draw() {
	for _, s := range g.children {
		s.Draw()
	}
}

func (g *Group) Rotate(angle float32) {
	// TODO
}

func (g *Group) Scale(sx, sy float32) {
	for _, s := range g.children {
		s.Scale(sx, sy)
	}
}

func (g *Group) Move(x, y float32) {
	// TODO
}

func (g *Group) Center() (float32, float32) {
	return g.x, g.y
}

func (g *Group) Angle() float32 {
	return g.angle
}

func (g *Group) Bounds() (float32, float32) {
	return g.w, g.h
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
