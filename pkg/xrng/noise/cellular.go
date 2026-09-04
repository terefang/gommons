package noise

import (
	"math"
	"math/rand"

	"github.com/terefang/gommons/pkg/xmath"
)

// CellularNoiseVariant specifies the distance metric algorithm.
type CellularNoiseVariant int

const (
	EUCLIDEAN CellularNoiseVariant = iota
	MANHATTAN
	NATURAL
)

// CellularReturnType specifies the mathematical transformation applied to nearest feature distances.
type CellularReturnType int

const (
	DISTANCE CellularReturnType = iota
	DISTANCE_2
	DISTANCE_2_ADD
	DISTANCE_2_SUB
	DISTANCE_2_MUL
	DISTANCE_2_DIV
	NOISE_LOOKUP
	CELL_VALUE
)

// CellularNoise implements multi-dimensional Worley/Cellular Noise.
type CellularNoise struct {
	perm       [512]uint8
	Variant    CellularNoiseVariant
	ReturnType CellularReturnType
}

// NewCellularNoise creates a new CellularNoise generator initialized with a seed and default configurations.
func NewCellularNoise(seed int64) *CellularNoise {
	return NewCellularNoiseWithSource(rand.NewSource(seed))
}

// NewCellularNoiseWithSource creates a new CellularNoise generator using a rand.Source.
func NewCellularNoiseWithSource(source rand.Source) *CellularNoise {
	var p [256]uint8
	for i := range p {
		p[i] = uint8(i)
	}
	r := rand.New(source)
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}

	cn := &CellularNoise{
		Variant:    EUCLIDEAN,
		ReturnType: DISTANCE,
	}
	for i := 0; i < 512; i++ {
		cn.perm[i] = p[i&0xff]
	}
	return cn
}

func featureOffset(hash uint8) float64 {
	return float64(hash) / 256.0
}

func (cn *CellularNoise) calcDist1D(dx float64) float64 {
	return math.Abs(dx)
}

func (cn *CellularNoise) calcDist2D(dx, dy float64) float64 {
	adx, ady := math.Abs(dx), math.Abs(dy)
	switch cn.Variant {
	case MANHATTAN:
		return adx + ady
	case NATURAL:
		return (math.Sqrt(dx*dx+dy*dy) + adx + ady) * 0.5
	case EUCLIDEAN:
		fallthrough
	default:
		return math.Sqrt(dx*dx + dy*dy)
	}
}

func (cn *CellularNoise) calcDist3D(dx, dy, dz float64) float64 {
	adx, ady, adz := math.Abs(dx), math.Abs(dy), math.Abs(dz)
	switch cn.Variant {
	case MANHATTAN:
		return adx + ady + adz
	case NATURAL:
		return (math.Sqrt(dx*dx+dy*dy+dz*dz) + adx + ady + adz) * 0.5
	case EUCLIDEAN:
		fallthrough
	default:
		return math.Sqrt(dx*dx + dy*dy + dz*dz)
	}
}

func (cn *CellularNoise) calcDist4D(dx, dy, dz, dw float64) float64 {
	adx, ady, adz, adw := math.Abs(dx), math.Abs(dy), math.Abs(dz), math.Abs(dw)
	switch cn.Variant {
	case MANHATTAN:
		return adx + ady + adz + adw
	case NATURAL:
		return (math.Sqrt(dx*dx+dy*dy+dz*dz+dw*dw) + adx + ady + adz + adw) * 0.5
	case EUCLIDEAN:
		fallthrough
	default:
		return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw)
	}
}

func (cn *CellularNoise) calcDist5D(dx, dy, dz, dw, dv float64) float64 {
	adx, ady, adz, adw, adv := math.Abs(dx), math.Abs(dy), math.Abs(dz), math.Abs(dw), math.Abs(dv)
	switch cn.Variant {
	case MANHATTAN:
		return adx + ady + adz + adw + adv
	case NATURAL:
		return (math.Sqrt(dx*dx+dy*dy+dz*dz+dw*dw+dv*dv) + adx + ady + adz + adw + adv) * 0.5
	case EUCLIDEAN:
		fallthrough
	default:
		return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw + dv*dv)
	}
}

func (cn *CellularNoise) calcDist6D(dx, dy, dz, dw, dv, du float64) float64 {
	adx, ady, adz, adw, adv, adu := math.Abs(dx), math.Abs(dy), math.Abs(dz), math.Abs(dw), math.Abs(dv), math.Abs(du)
	switch cn.Variant {
	case MANHATTAN:
		return adx + ady + adz + adw + adv + adu
	case NATURAL:
		return (math.Sqrt(dx*dx+dy*dy+dz*dz+dw*dw+dv*dv+du*du) + adx + ady + adz + adw + adv + adu) * 0.5
	case EUCLIDEAN:
		fallthrough
	default:
		return math.Sqrt(dx*dx + dy*dy + dz*dz + dw*dw + dv*dv + du*du)
	}
}

func (cn *CellularNoise) evaluateResult(f1, f2 float64, winningHash uint8, scale float64) float64 {
	switch cn.ReturnType {
	case DISTANCE_2:
		return (f2 * scale) - 1.0
	case DISTANCE_2_ADD:
		return ((f1 + f2) * scale * 0.5) - 1.0
	case DISTANCE_2_SUB:
		// Uses f2 - f1 to represent boundary edge distance from the secondary node
		return ((f2 - f1) * scale) - 1.0
	case DISTANCE_2_MUL:
		return (f1 * f2 * scale) - 1.0
	case DISTANCE_2_DIV:
		if f2 == 0 {
			return -1.0
		}
		return ((f1 / f2) * 2.0) - 1.0
	case NOISE_LOOKUP:
		return (float64(cn.perm[winningHash]) / 127.5) - 1.0
	case CELL_VALUE:
		return (float64(winningHash) / 127.5) - 1.0
	case DISTANCE:
		fallthrough
	default:
		return (f1 * scale) - 1.0
	}
}

// Noise1D returns 1D Cellular Noise scaled to [-1, 1].
func (cn *CellularNoise) Noise1D(x float64) float64 {
	cell := xmath.FastFloor(x)
	f1, f2 := 1e9, 1e9
	var winningHash uint8

	for i := -1; i <= 1; i++ {
		neighborCell := cell + i
		hash := cn.perm[uint8(neighborCell)]
		featurePoint := float64(neighborCell) + featureOffset(hash)

		dist := cn.calcDist1D(x - featurePoint)

		if dist < f1 {
			f2 = f1
			f1 = dist
			winningHash = hash
		} else if dist < f2 {
			f2 = dist
		}
	}

	return cn.evaluateResult(f1, f2, winningHash, 2.0)
}

// Noise2D returns 2D Cellular Noise scaled to [-1, 1].
func (cn *CellularNoise) Noise2D(x, y float64) float64 {
	cellX, cellY := xmath.FastFloor(x), xmath.FastFloor(y)
	f1, f2 := 1e9, 1e9
	var winningHash uint8

	for ix := -1; ix <= 1; ix++ {
		for iy := -1; iy <= 1; iy++ {
			nx, ny := cellX+ix, cellY+iy

			hashX := cn.perm[uint8(nx)+cn.perm[uint8(ny)]]
			hashY := cn.perm[uint8(ny)+cn.perm[uint8(nx)]]

			featureX := float64(nx) + featureOffset(hashX)
			featureY := float64(ny) + featureOffset(hashY)

			dist := cn.calcDist2D(x-featureX, y-featureY)

			if dist < f1 {
				f2 = f1
				f1 = dist
				winningHash = hashX
			} else if dist < f2 {
				f2 = dist
			}
		}
	}

	return cn.evaluateResult(f1, f2, winningHash, 1.4142)
}

// Noise3D returns 3D Cellular Noise scaled to [-1, 1].
func (cn *CellularNoise) Noise3D(x, y, z float64) float64 {
	cellX, cellY, cellZ := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z)
	f1, f2 := 1e9, 1e9
	var winningHash uint8

	for ix := -1; ix <= 1; ix++ {
		for iy := -1; iy <= 1; iy++ {
			for iz := -1; iz <= 1; iz++ {
				nx, ny, nz := cellX+ix, cellY+iy, cellZ+iz

				hZ := cn.perm[uint8(nz)]
				hY := cn.perm[uint8(ny)+hZ]

				hashX := cn.perm[uint8(nx)+hY]
				hashY := cn.perm[uint8(ny)+hashX]
				hashZ := cn.perm[uint8(nz)+hashY]

				featureX := float64(nx) + featureOffset(hashX)
				featureY := float64(ny) + featureOffset(hashY)
				featureZ := float64(nz) + featureOffset(hashZ)

				dist := cn.calcDist3D(x-featureX, y-featureY, z-featureZ)

				if dist < f1 {
					f2 = f1
					f1 = dist
					winningHash = hashX
				} else if dist < f2 {
					f2 = dist
				}
			}
		}
	}

	return cn.evaluateResult(f1, f2, winningHash, 1.1547)
}

// Noise4D returns 4D Cellular Noise scaled to [-1, 1].
func (cn *CellularNoise) Noise4D(x, y, z, w float64) float64 {
	cellX, cellY, cellZ, cellW := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w)
	f1, f2 := 1e9, 1e9
	var winningHash uint8

	for ix := -1; ix <= 1; ix++ {
		for iy := -1; iy <= 1; iy++ {
			for iz := -1; iz <= 1; iz++ {
				for iw := -1; iw <= 1; iw++ {
					nx, ny, nz, nw := cellX+ix, cellY+iy, cellZ+iz, cellW+iw

					hW := cn.perm[uint8(nw)]
					hZ := cn.perm[uint8(nz)+hW]
					hY := cn.perm[uint8(ny)+hZ]

					hashX := cn.perm[uint8(nx)+hY]
					hashY := cn.perm[uint8(ny)+hashX]
					hashZ := cn.perm[uint8(nz)+hashY]
					hashW := cn.perm[uint8(nw)+hashZ]

					featureX := float64(nx) + featureOffset(hashX)
					featureY := float64(ny) + featureOffset(hashY)
					featureZ := float64(nz) + featureOffset(hashZ)
					featureW := float64(nw) + featureOffset(hashW)

					dist := cn.calcDist4D(x-featureX, y-featureY, z-featureZ, w-featureW)

					if dist < f1 {
						f2 = f1
						f1 = dist
						winningHash = hashX
					} else if dist < f2 {
						f2 = dist
					}
				}
			}
		}
	}

	return cn.evaluateResult(f1, f2, winningHash, 1.0)
}

// Noise5D returns 5D Cellular Noise scaled to [-1, 1].
func (cn *CellularNoise) Noise5D(x, y, z, w, v float64) float64 {
	cellX, cellY, cellZ, cellW, cellV := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v)
	f1, f2 := 1e9, 1e9
	var winningHash uint8

	for ix := -1; ix <= 1; ix++ {
		for iy := -1; iy <= 1; iy++ {
			for iz := -1; iz <= 1; iz++ {
				for iw := -1; iw <= 1; iw++ {
					for iv := -1; iv <= 1; iv++ {
						nx, ny, nz, nw, nv := cellX+ix, cellY+iy, cellZ+iz, cellW+iw, cellV+iv

						hV := cn.perm[uint8(nv)]
						hW := cn.perm[uint8(nw)+hV]
						hZ := cn.perm[uint8(nz)+hW]
						hY := cn.perm[uint8(ny)+hZ]

						hashX := cn.perm[uint8(nx)+hY]
						hashY := cn.perm[uint8(ny)+hashX]
						hashZ := cn.perm[uint8(nz)+hashY]
						hashW := cn.perm[uint8(nw)+hashZ]
						hashV := cn.perm[uint8(nv)+hashW]

						featureX := float64(nx) + featureOffset(hashX)
						featureY := float64(ny) + featureOffset(hashY)
						featureZ := float64(nz) + featureOffset(hashZ)
						featureW := float64(nw) + featureOffset(hashW)
						featureV := float64(nv) + featureOffset(hashV)

						dist := cn.calcDist5D(x-featureX, y-featureY, z-featureZ, w-featureW, v-featureV)

						if dist < f1 {
							f2 = f1
							f1 = dist
							winningHash = hashX
						} else if dist < f2 {
							f2 = dist
						}
					}
				}
			}
		}
	}

	return cn.evaluateResult(f1, f2, winningHash, 0.89)
}

// Noise6D returns 6D Cellular Noise scaled to [-1, 1].
func (cn *CellularNoise) Noise6D(x, y, z, w, v, u float64) float64 {
	cellX, cellY, cellZ, cellW, cellV, cellU := xmath.FastFloor(x), xmath.FastFloor(y), xmath.FastFloor(z), xmath.FastFloor(w), xmath.FastFloor(v), xmath.FastFloor(u)
	f1, f2 := 1e9, 1e9
	var winningHash uint8

	for ix := -1; ix <= 1; ix++ {
		for iy := -1; iy <= 1; iy++ {
			for iz := -1; iz <= 1; iz++ {
				for iw := -1; iw <= 1; iw++ {
					for iv := -1; iv <= 1; iv++ {
						for iu := -1; iu <= 1; iu++ {
							nx, ny, nz, nw, nv, nu := cellX+ix, cellY+iy, cellZ+iz, cellW+iw, cellV+iv, cellU+iu

							hU := cn.perm[uint8(nu)]
							hV := cn.perm[uint8(nv)+hU]
							hW := cn.perm[uint8(nw)+hV]
							hZ := cn.perm[uint8(nz)+hW]
							hY := cn.perm[uint8(ny)+hZ]

							hashX := cn.perm[uint8(nx)+hY]
							hashY := cn.perm[uint8(ny)+hashX]
							hashZ := cn.perm[uint8(nz)+hashY]
							hashW := cn.perm[uint8(nw)+hashZ]
							hashV := cn.perm[uint8(nv)+hashW]
							hashU := cn.perm[uint8(nu)+hashV]

							featureX := float64(nx) + featureOffset(hashX)
							featureY := float64(ny) + featureOffset(hashY)
							featureZ := float64(nz) + featureOffset(hashZ)
							featureW := float64(nw) + featureOffset(hashW)
							featureV := float64(nv) + featureOffset(hashV)
							featureU := float64(nu) + featureOffset(hashU)

							dist := cn.calcDist6D(x-featureX, y-featureY, z-featureZ, w-featureW, v-featureV, u-featureU)

							if dist < f1 {
								f2 = f1
								f1 = dist
								winningHash = hashX
							} else if dist < f2 {
								f2 = dist
							}
						}
					}
				}
			}
		}
	}

	return cn.evaluateResult(f1, f2, winningHash, 0.81)
}
