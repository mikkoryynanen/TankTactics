package collision

type AABB struct {
    Min, Max Vector3 // Min and Max points of the box
}

type Vector3 struct {
    X, Y, Z float64
}

// Check if two AABBs are colliding
func IsAABBColliding(a, b AABB) bool {
    return (a.Min.X <= b.Max.X && a.Max.X >= b.Min.X) &&
           (a.Min.Y <= b.Max.Y && a.Max.Y >= b.Min.Y) &&
           (a.Min.Z <= b.Max.Z && a.Max.Z >= b.Min.Z)
}
