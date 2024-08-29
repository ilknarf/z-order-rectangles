package main

import (
	"math"
)

const (
	maxDimBits = 16            // 64 bits / 4 dimensions
	maxDimVal  = (1 << 16) - 1 // calculate max value of dimension
)

// rough implementation of a 4D  z-order curve (x0, x1, y0, y1)
type RectHash struct {
	// the 4d rectangle is encoded by interleaving bits of each dimension
	Val uint64
}

func (r *RectHash) X0() uint {
	// it's the first bit
	return lookupVal(r.Val, 3)
}

func (r *RectHash) SetX0(val uint) {
	r.Val = setVal(r.Val, val, 3)
}

func (r *RectHash) X1() uint {
	return lookupVal(r.Val, 2)
}

func (r *RectHash) SetX1(val uint) {
	r.Val = setVal(r.Val, val, 2)
}

func (r *RectHash) Y0() uint {
	return lookupVal(r.Val, 1)
}

func (r *RectHash) SetY0(val uint) {
	r.Val = setVal(r.Val, val, 1)
}

func (r *RectHash) Y1() uint {
	return lookupVal(r.Val, 0)
}

func (r *RectHash) SetY1(val uint) {
	r.Val = setVal(r.Val, val, 0)
}

func lookupVal(val uint64, start int) uint {
	var res uint

	for i := 0; i < maxDimBits; i += 1 {
		queryMask := uint64(1 << (i*4 + start))

		if queryMask&val > 0 {
			res |= 1 << i
		}
	}

	return res
}

// assumes that it's within bounds, (negatives will get messed up and larger numbers truncated)
func setVal(encodedVal uint64, varVal uint, start int) uint64 {
	res := encodedVal

	// clear set var using mask
	var clearMask uint64 = math.MaxUint64

	for i := 0; i < maxDimBits; i += 1 {
		// set all bits of what we want to set to 0
		clearMask ^= 1 << (i*4 + start)
	}

	res &= clearMask

	for i := 0; i < maxDimBits; i += 1 {
		if varVal&(1<<i) > 0 {
			res |= 1 << (i*4 + start)
		}
	}

	return res
}
