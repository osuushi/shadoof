# ShaDooF

ShaDooF is a hodgepodge library of helpers for use with
[deadsy/sdfx](https://github.com/deadsy/sdfx) for programmatic CAD via signed
distance fields.

These utilities are intended to be pragmatic tools for quick hacking, which
means that they often come with caveats, and these caveats are documented as
thoroughly as possible.

In particular, some of the SDFs in this library do not produce correct signs,
but incorrect magnitudes when evaluated. This means that they will render
correctly in isolation, but some more complex operations like
[sdf.Shell3D](https://pkg.go.dev/github.com/deadsy/sdfx/sdf#Shell3D) will not
produce exactly correct results. This is the same caveat that applies to
[sdf.Transform3D](https://pkg.go.dev/github.com/deadsy/sdfx/sdf#Transform3D) for
scale transforms.

## Transforms

transform.go provides a number of convenience functions for positioning SDFs.
Most of these deal with centering and aligning two SDFs based on their bounding
boxes. It is important to note that these functions will only work correctly if
the bounding boxes are minimal, and that this is not guaranteed for all SDFs.

Bounding boxes are fundamentally a performance facility, not a measurement
facility. For example, rotations by non-integer multiples of π/2 will always
cause an SDFs bounding box to expand, even if the _minimal_ bounding box is
unchanged. Even a sphere, which has rotational symmetry in every direction will
see its bounding box expand under rotation. In these cases, the alignment
functions will not behave as you expect.

Nonetheless, as a pragmatic tool, alignment via bounding boxes is extremely
convenient and useful, so long as you understand its limitations.