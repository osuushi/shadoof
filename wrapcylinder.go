package shadoof

import (
	"math"

	"github.com/deadsy/sdfx/sdf"
	v2 "github.com/deadsy/sdfx/vec/v2"
	v3 "github.com/deadsy/sdfx/vec/v3"
)

// An SDF2 wrapped around a cylinder
type WrapSylinderSDF struct {
	SDF2          sdf.SDF2
	width         float64
	outerCylinder sdf.SDF3
	innerCylinder sdf.SDF3
	innerMax      sdf.MaxFunc
	outerMax      sdf.MaxFunc
}

// Create an SDF3 by wrapping an SDF2 around a cylinder with a given thickness.
//
// Lengths are preserved on the inside of the wrapped SDF, and points are
// sampled for -r/2 < x < r/2, where r is the inner radius, so if the 2D SDF is
// wider than the circumference of the cylinder, it will be clipped.
//
// Both the inside and outside can be rounded independently via PolyMax.
//
// Note that if you wish to subtract this from a cylinder (for embossing, for
// example), you will need to subtract the thickness from the radius you
// provide. So if you are trying to create a seamless embossing around the
// entire cylinder, you must ensure that the 2D SDF's width is the 2*pi*r, where
// r is that adjusted radius.
func WrapAroundCylinder(sdf2d sdf.SDF2, innerRadius, thickness, roundInside, roundOutside float64) sdf.SDF3 {
	outerRadius := innerRadius + thickness
	// Compute the width, which is half the circumference of the inner cylinder.
	width := math.Pi * innerRadius
	sdfBb := sdf2d.BoundingBox()
	height := sdfBb.Max.Y - sdfBb.Min.Y
	outerCylinder, _ := sdf.Cylinder3D(height, outerRadius, 0)
	innerCylinder, _ := sdf.Cylinder3D(height, innerRadius, 0)

	// Move the cylinders so they are aligned with the bounding box in the y direction
	outerCylinder = Translate3D(outerCylinder, 0, 0, height/2+sdfBb.Min.Y)
	innerCylinder = Translate3D(innerCylinder, 0, 0, height/2+sdfBb.Min.Y)

	var innerMax, outerMax sdf.MaxFunc

	if roundInside == 0 {
		innerMax = math.Max
	} else {
		innerMax = sdf.PolyMax(roundInside)
	}

	if roundOutside == 0 {
		outerMax = math.Max
	} else {
		outerMax = sdf.PolyMax(roundOutside)
	}

	return &WrapSylinderSDF{
		SDF2:          sdf2d,
		width:         width,
		outerCylinder: outerCylinder,
		innerCylinder: innerCylinder,
		innerMax:      innerMax,
		outerMax:      outerMax,
	}
}

func (s *WrapSylinderSDF) BoundingBox() sdf.Box3 {
	return s.outerCylinder.BoundingBox()
}

func (s *WrapSylinderSDF) Evaluate(p v3.Vec) float64 {
	var dist float64
	// Convert the 3D point to a 2D point to sample the distance from the 2D SDF.
	p2 := v2.Vec{Y: p.Z}

	angle := math.Atan2(p.Y, p.X)
	angle += math.Pi / 2
	p2.X = s.width * angle / math.Pi

	dist = s.SDF2.Evaluate(p2)
	// Constrain to inside outer cylinder
	dist = s.outerMax(dist, s.outerCylinder.Evaluate(p))
	// Subtract inner cylinder
	dist = s.innerMax(dist, -s.innerCylinder.Evaluate(p))
	return dist
}
