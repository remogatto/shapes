package shapes

import "image/color"

var (
	// The default color for shapes is blue.
	DefaultColor = color.RGBA{0, 0, 0xff, 0xff}
)

// // Calculate bounds for a generic, possibly merged, slice of vertices
// func calcBounds(vertices []float32) (float32, float32) {
// 	minX, minY := math.Inf(1), math.Inf(1)
// 	maxX, maxY := math.Inf(-1), math.Inf(-1)

// 	for i, v := range vertices {
// 		if i%2 == 0 && float64(v) < minX {
// 			minX = float64(v)
// 		}
// 		if i%2 == 1 && float64(v) < minY {
// 			minY = float64(v)
// 		}
// 		if i%2 == 0 && float64(v) > maxX {
// 			maxX = float64(v)
// 		}
// 		if i%2 == 1 && float64(v) > maxY {
// 			maxY = float64(v)
// 		}
// 	}

// 	w = float32(math.Abs(maxX - minX))
// 	h = float32(math.Abs(maxY - minY))

// 	return w, h
// }
