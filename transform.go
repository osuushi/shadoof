package shadoof

import (
	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

type Side int

const (
	Left Side = iota
	Right
	Top
	Bottom
	Front
	Back
)

// Convenience function for translating a 2D SDF, with less boilerplate
func Translate2D(s sdf.SDF2, x float64, y float64) sdf.SDF2 {
	return sdf.Transform2D(s, sdf.Translate2d(v2.Vec{X: x, Y: y}))
}

// Convenience function for translating a 3D SDF, with less boilerplate
func Translate3D(s sdf.SDF3, x float64, y float64, z float64) sdf.SDF3 {
	return sdf.Transform3D(s, sdf.Translate3d(v3.Vec{X: x, Y: y, Z: z}))
}

// Set an SDF on the XY plane.
func SetOnXYPlane(s sdf.SDF3) sdf.SDF3 {
	bb := s.BoundingBox()
	return Translate3D(s, 0, 0, -bb.Min.Z)
}

// Move s2 so that its bounding box is centered on s1's bounding box.
func CenterWith2D(s1, s2 sdf.SDF2) sdf.SDF2 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	offset := bb1.Center().Sub(bb2.Center())
	return Translate2D(s2, offset.X, offset.Y)
}

// Move s2 so that its bounding box is centered on s1's bounding box. That is,
// it centers the X and Y coordinates, ignoring s1's Z coordinate.
func Center2DWith3D(s1 sdf.SDF3, s2 sdf.SDF2) sdf.SDF2 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	center1 := bb1.Center()
	offset := v2.Vec{X: center1.X, Y: center1.Y}.Sub(bb2.Center())
	return Translate2D(s2, offset.X, offset.Y)
}

// Move s2 so that its bounding box is centered on s1's bounding box.
func CenterWith3D(s1, s2 sdf.SDF3) sdf.SDF3 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	offset := bb1.Center().Sub(bb2.Center())
	return Translate3D(s2, offset.X, offset.Y, offset.Z)
}

// Move s2 so that its bounding box abuts the given side of s1's bounding box.
// Only one coordinate is ever translated.
func AbutWith2D(s1, s2 sdf.SDF2, side Side) sdf.SDF2 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	switch side {
	case Left:
		// Move the right side of bb2 to the left side of bb1
		return Translate2D(s2, bb1.Min.X-bb2.Max.X, 0)
	case Right:
		// Move the left side of bb2 to the right side of bb1
		return Translate2D(s2, bb1.Max.X-bb2.Min.X, 0)
	case Bottom:
		// Move the top side of bb2 to the bottom side of bb1
		return Translate2D(s2, 0, bb1.Max.Y-bb2.Min.Y)
	case Top:
		// Move the bottom side of bb2 to the top side of bb1
		return Translate2D(s2, 0, bb1.Min.Y-bb2.Max.Y)
	default:
		panic("Invalid side for 2D abut")
	}
}

// Move s2 so that its bounding box abuts the given side of s1's bounding box.
// Only one coordinate is ever translated. Note that "Top" and "Bottom" refer to
// the Z coordinate, unlike in the 2D case where they refer to the Y coordinate.
// For the 3D case, "Front" refers to negative Y.
func AbutWith3D(s1, s2 sdf.SDF3, side Side) sdf.SDF3 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	switch side {
	case Left:
		// Move the right side of bb2 to the left side of bb1
		return Translate3D(s2, bb1.Min.X-bb2.Max.X, 0, 0)
	case Right:
		// Move the left side of bb2 to the right side of bb1
		return Translate3D(s2, bb1.Max.X-bb2.Min.X, 0, 0)
	case Front:
		// Move the back side of bb2 to the front side of bb1
		return Translate3D(s2, 0, bb1.Min.Y-bb2.Max.Y, 0)
	case Back:
		// Move the front side of bb2 to the back side of bb1
		return Translate3D(s2, 0, bb1.Max.Y-bb2.Min.Y, 0)
	case Bottom:
		// Move the top side of bb2 to the bottom side of bb1
		return Translate3D(s2, 0, 0, bb1.Min.Z-bb2.Max.Z)
	case Top:
		// Move the bottom side of bb2 to the top side of bb1
		return Translate3D(s2, 0, 0, bb1.Max.Z-bb2.Min.Z)
	default:
		panic("Invalid side for 3D abut")
	}
}

// Move s2 so that the given side of its bounding box aligns with the same side
// of s1's bounding box.
func AlignWith2D(s1, s2 sdf.SDF2, side Side) sdf.SDF2 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	switch side {
	case Left:
		return Translate2D(s2, bb1.Min.X-bb2.Min.X, 0)
	case Right:
		return Translate2D(s2, bb1.Max.X-bb2.Max.X, 0)
	case Bottom:
		return Translate2D(s2, 0, bb1.Min.Y-bb2.Min.Y)
	case Top:
		return Translate2D(s2, 0, bb1.Max.Y-bb2.Max.Y)
	default:
		panic("Invalid side for 2D align")
	}
}

// Move s2 so that the given side of its bounding box aligns with the same side
// of s1's bounding box. Note that "Top" and "Bottom" refer to the Z coordinate,
// unlike in the 2D case where they refer to the Y coordinate. For the 3D case,
// "Front" refers to negative Y.
func AlignWith3D(s1, s2 sdf.SDF3, side Side) sdf.SDF3 {
	bb1 := s1.BoundingBox()
	bb2 := s2.BoundingBox()
	switch side {
	case Left:
		return Translate3D(s2, bb1.Min.X-bb2.Min.X, 0, 0)
	case Right:
		return Translate3D(s2, bb1.Max.X-bb2.Max.X, 0, 0)
	case Front:
		return Translate3D(s2, 0, bb1.Min.Y-bb2.Min.Y, 0)
	case Back:
		return Translate3D(s2, 0, bb1.Max.Y-bb2.Max.Y, 0)
	case Bottom:
		return Translate3D(s2, 0, 0, bb1.Min.Z-bb2.Min.Z)
	case Top:
		return Translate3D(s2, 0, 0, bb1.Max.Z-bb2.Max.Z)
	default:
		panic("Invalid side for 3D align")
	}
}
