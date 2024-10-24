package collision

import (
    "math"
)

// Box represents an axis-aligned bounding box
type Box struct {
    Width, Height, Depth float64
}

// CalculateRadius calculates the radius of the bounding sphere for the box
func (b Box) CalculateRadius() float64 {
    diagonal := math.Sqrt(b.Width * b.Width + b.Height * b.Height + b.Depth * b.Depth)
    return diagonal / 2
}