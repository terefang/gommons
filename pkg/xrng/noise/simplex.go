package noise

import (
	"math/rand"

	"github.com/terefang/gommons/pkg/xmath"
)

// Simplex generates deterministic 1D, 2D, and 3D Simplex noise.
type Simplex struct {
	perm [512]uint8
}

// NewSimplex creates a Simplex noise generator from a seed.
func NewSimplex(seed int64) *Simplex {
	return NewSimplexWithSource(rand.NewSource(seed))
}

func NewSimplexWithSource(source rand.Source) *Simplex {
	var p [256]uint8
	for i, _ := range p {
		p[i] = uint8(i)
	}
	r := rand.New(source)
	for i := 255; i > 0; i-- {
		j := r.Intn(i + 1)
		p[i], p[j] = p[j], p[i]
	}

	s := &Simplex{}
	for i := 0; i < 512; i++ {
		s.perm[i] = p[i&0xff]
	}
	return s
}

// Noise1D returns 1D Simplex noise scaled to approximately [-1, 1].
func (s *Simplex) Noise1D(x float64) float64 {
	i0 := xmath.FastFloor(x)
	i1 := i0 + 1

	x0 := x - float64(i0)
	x1 := x0 - 1.0

	var n0, n1 float64

	// Left corner contribution
	if t0 := 1.0 - x0*x0; t0 > 0 {
		t0 *= t0
		hash := s.perm[uint8(i0)]
		n0 = t0 * t0 * grad1(hash, x0)
	}

	// Right corner contribution
	if t1 := 1.0 - x1*x1; t1 > 0 {
		t1 *= t1
		hash := s.perm[uint8(i1)]
		n1 = t1 * t1 * grad1(hash, x1)
	}

	return 0.395 * (n0 + n1)
}

func grad1(hash uint8, x float64) float64 {
	grad := float64(1 + (hash & 7))
	if (hash & 8) != 0 {
		grad = -grad
	}
	return grad * x
}

// Noise2D returns 2D Simplex noise scaled to approximately [-1, 1].
func (s *Simplex) Noise2D(x, y float64) float64 {
	const (
		F2 = 0.3660254037844386 // 0.5 * (sqrt(3) - 1)
		G2 = 0.2113248654051871 // (3 - sqrt(3)) / 6
	)

	skew := (x + y) * F2
	i := xmath.FastFloor(x + skew)
	j := xmath.FastFloor(y + skew)

	t := float64(i+j) * G2
	x0 := x - (float64(i) - t)
	y0 := y - (float64(j) - t)

	var i1, j1 int
	if x0 > y0 {
		i1, j1 = 1, 0
	} else {
		i1, j1 = 0, 1
	}

	x1 := x0 - float64(i1) + G2
	y1 := y0 - float64(j1) + G2
	x2 := x0 - 1.0 + (2.0 * G2)
	y2 := y0 - 1.0 + (2.0 * G2)

	ii := uint8(i)
	jj := uint8(j)

	var n0, n1, n2 float64

	if t0 := 0.5 - x0*x0 - y0*y0; t0 > 0 {
		t0 *= t0
		hash := s.perm[ii+s.perm[jj]]
		n0 = t0 * t0 * grad2Inline(hash, x0, y0)
	}

	if t1 := 0.5 - x1*x1 - y1*y1; t1 > 0 {
		t1 *= t1
		hash := s.perm[ii+uint8(i1)+s.perm[jj+uint8(j1)]]
		n1 = t1 * t1 * grad2Inline(hash, x1, y1)
	}

	if t2 := 0.5 - x2*x2 - y2*y2; t2 > 0 {
		t2 *= t2
		hash := s.perm[ii+1+s.perm[jj+1]]
		n2 = t2 * t2 * grad2Inline(hash, x2, y2)
	}

	return 70.0 * (n0 + n1 + n2)
}

func grad2Inline(hash uint8, x, y float64) float64 {
	h := hash & 7
	u := x
	if h >= 4 {
		u = y
	}
	v := y
	if h >= 4 {
		v = x
	}
	if (h & 1) != 0 {
		u = -u
	}
	if (h & 2) != 0 {
		v = -v
	}
	if h == 4 || h == 5 {
		return u
	}
	if h == 6 || h == 7 {
		return v
	}
	return u + v
}

// Noise3D returns 3D Simplex noise scaled to approximately [-1, 1].
func (s *Simplex) Noise3D(x, y, z float64) float64 {
	const (
		F3 = 1.0 / 3.0
		G3 = 1.0 / 6.0
	)

	// Skew space to find tetrahedral grid cell origin
	skew := (x + y + z) * F3
	i := xmath.FastFloor(x + skew)
	j := xmath.FastFloor(y + skew)
	k := xmath.FastFloor(z + skew)

	// Unskew origin back to (x, y, z) space
	t := float64(i+j+k) * G3
	x0 := x - (float64(i) - t)
	y0 := y - (float64(j) - t)
	z0 := z - (float64(k) - t)

	// Determine which tetrahedron simplex corner traversal path to take
	var i1, j1, k1 int
	var i2, j2, k2 int

	if x0 >= y0 {
		if y0 >= z0 {
			i1, j1, k1 = 1, 0, 0
			i2, j2, k2 = 1, 1, 0
		} else if x0 >= z0 {
			i1, j1, k1 = 1, 0, 0
			i2, j2, k2 = 1, 0, 1
		} else {
			i1, j1, k1 = 0, 0, 1
			i2, j2, k2 = 1, 0, 1
		}
	} else {
		if y0 < z0 {
			i1, j1, k1 = 0, 0, 1
			i2, j2, k2 = 0, 1, 1
		} else if x0 < z0 {
			i1, j1, k1 = 0, 1, 0
			i2, j2, k2 = 0, 1, 1
		} else {
			i1, j1, k1 = 0, 1, 0
			i2, j2, k2 = 1, 1, 0
		}
	}

	// Calculate unskewed relative positions for the 4 simplex corners
	x1 := x0 - float64(i1) + G3
	y1 := y0 - float64(j1) + G3
	z1 := z0 - float64(k1) + G3

	x2 := x0 - float64(i2) + 2.0*G3
	y2 := y0 - float64(j2) + 2.0*G3
	z2 := z0 - float64(k2) + 2.0*G3

	x3 := x0 - 1.0 + 3.0*G3
	y3 := y0 - 1.0 + 3.0*G3
	z3 := z0 - 1.0 + 3.0*G3

	ii := uint8(i)
	jj := uint8(j)
	kk := uint8(k)

	var n0, n1, n2, n3 float64

	// Corner 0 contribution
	if t0 := 0.6 - x0*x0 - y0*y0 - z0*z0; t0 > 0 {
		t0 *= t0
		hash := s.perm[ii+s.perm[jj+s.perm[kk]]]
		n0 = t0 * t0 * grad3Inline(hash, x0, y0, z0)
	}

	// Corner 1 contribution
	if t1 := 0.6 - x1*x1 - y1*y1 - z1*z1; t1 > 0 {
		t1 *= t1
		hash := s.perm[ii+uint8(i1)+s.perm[jj+uint8(j1)+s.perm[kk+uint8(k1)]]]
		n1 = t1 * t1 * grad3Inline(hash, x1, y1, z1)
	}

	// Corner 2 contribution
	if t2 := 0.6 - x2*x2 - y2*y2 - z2*z2; t2 > 0 {
		t2 *= t2
		hash := s.perm[ii+uint8(i2)+s.perm[jj+uint8(j2)+s.perm[kk+uint8(k2)]]]
		n2 = t2 * t2 * grad3Inline(hash, x2, y2, z2)
	}

	// Corner 3 contribution
	if t3 := 0.6 - x3*x3 - y3*y3 - z3*z3; t3 > 0 {
		t3 *= t3
		hash := s.perm[ii+1+s.perm[jj+1+s.perm[kk+1]]]
		n3 = t3 * t3 * grad3Inline(hash, x3, y3, z3)
	}

	return 32.0 * (n0 + n1 + n2 + n3)
}

// Inlined 3D gradient vector evaluation over 12 unit edge directions
func grad3Inline(hash uint8, x, y, z float64) float64 {
	h := hash & 15
	u := x
	if h >= 8 {
		u = y
	}
	v := y
	if h >= 4 {
		if h == 12 || h == 14 {
			v = x
		} else {
			v = z
		}
	}
	if (h & 1) != 0 {
		u = -u
	}
	if (h & 2) != 0 {
		v = -v
	}
	return u + v
}

// Noise4D returns 4D Simplex noise scaled to approximately [-1, 1].
func (s *Simplex) Noise4D(x, y, z, w float64) float64 {
	const (
		F4 = 0.309016994374947451 // (sqrt(5) - 1) / 4
		G4 = 0.138196601125010515 // (5 - sqrt(5)) / 20
	)

	// Skew space to find 4D hypercube cell origin
	skew := (x + y + z + w) * F4
	i := xmath.FastFloor(x + skew)
	j := xmath.FastFloor(y + skew)
	k := xmath.FastFloor(z + skew)
	l := xmath.FastFloor(w + skew)

	// Unskew origin back to (x, y, z, w) space
	t := float64(i+j+k+l) * G4
	x0 := x - (float64(i) - t)
	y0 := y - (float64(j) - t)
	z0 := z - (float64(k) - t)
	w0 := w - (float64(l) - t)

	// Bitwise magnitude sorting to determine pentatope corner traversal order
	c1 := 0
	if x0 > y0 {
		c1 |= 1
	}
	if x0 > z0 {
		c1 |= 2
	}
	if x0 > w0 {
		c1 |= 4
	}
	if y0 > z0 {
		c1 |= 8
	}
	if y0 > w0 {
		c1 |= 16
	}
	if z0 > w0 {
		c1 |= 32
	}

	// Corner rank offsets derived from 4D simplex edge comparisons
	var i1, j1, k1, l1 int
	var i2, j2, k2, l2 int
	var i3, j3, k3, l3 int

	// Determine traversal offsets for 5 corners based on magnitude order
	if c1&1 != 0 {
		i1++
	} else {
		j1++
	}
	if c1&2 != 0 {
		i1++
	} else {
		k1++
	}
	if c1&4 != 0 {
		i1++
	} else {
		l1++
	}
	if c1&8 != 0 {
		j1++
	} else {
		k1++
	}
	if c1&16 != 0 {
		j1++
	} else {
		l1++
	}
	if c1&32 != 0 {
		k1++
	} else {
		l1++
	}

	if i1 >= 3 {
		i1 = 1
	} else {
		i1 = 0
	}
	if j1 >= 3 {
		j1 = 1
	} else {
		j1 = 0
	}
	if k1 >= 3 {
		k1 = 1
	} else {
		k1 = 0
	}
	if l1 >= 3 {
		l1 = 1
	} else {
		l1 = 0
	}

	// Corner rank 2
	i2, j2, k2, l2 = i1, j1, k1, l1
	if c1&1 != 0 {
		i2++
	} else {
		j2++
	}
	if c1&2 != 0 {
		i2++
	} else {
		k2++
	}
	if c1&4 != 0 {
		i2++
	} else {
		l2++
	}
	if c1&8 != 0 {
		j2++
	} else {
		k2++
	}
	if c1&16 != 0 {
		j2++
	} else {
		l2++
	}
	if c1&32 != 0 {
		k2++
	} else {
		l2++
	}

	if i2 >= 2 {
		i2 = 1
	} else {
		i2 = 0
	}
	if j2 >= 2 {
		j2 = 1
	} else {
		j2 = 0
	}
	if k2 >= 2 {
		k2 = 1
	} else {
		k2 = 0
	}
	if l2 >= 2 {
		l2 = 1
	} else {
		l2 = 0
	}

	// Corner rank 3
	i3, j3, k3, l3 = i2, j2, k2, l2
	if c1&1 != 0 {
		i3++
	} else {
		j3++
	}
	if c1&2 != 0 {
		i3++
	} else {
		k3++
	}
	if c1&4 != 0 {
		i3++
	} else {
		l3++
	}
	if c1&8 != 0 {
		j3++
	} else {
		k3++
	}
	if c1&16 != 0 {
		j3++
	} else {
		l3++
	}
	if c1&32 != 0 {
		k3++
	} else {
		l3++
	}

	if i3 >= 1 {
		i3 = 1
	} else {
		i3 = 0
	}
	if j3 >= 1 {
		j3 = 1
	} else {
		j3 = 0
	}
	if k3 >= 1 {
		k3 = 1
	} else {
		k3 = 0
	}
	if l3 >= 1 {
		l3 = 1
	} else {
		l3 = 0
	}

	// Unskewed relative coordinates for the 5 corners
	x1 := x0 - float64(i1) + G4
	y1 := y0 - float64(j1) + G4
	z1 := z0 - float64(k1) + G4
	w1 := w0 - float64(l1) + G4

	x2 := x0 - float64(i2) + 2.0*G4
	y2 := y0 - float64(j2) + 2.0*G4
	z2 := z0 - float64(k2) + 2.0*G4
	w2 := w0 - float64(l2) + 2.0*G4

	x3 := x0 - float64(i3) + 3.0*G4
	y3 := y0 - float64(j3) + 3.0*G4
	z3 := z0 - float64(k3) + 3.0*G4
	w3 := w0 - float64(l3) + 3.0*G4

	x4 := x0 - 1.0 + 4.0*G4
	y4 := y0 - 1.0 + 4.0*G4
	z4 := z0 - 1.0 + 4.0*G4
	w4 := w0 - 1.0 + 4.0*G4

	ii := uint8(i)
	jj := uint8(j)
	kk := uint8(k)
	ll := uint8(l)

	var n0, n1, n2, n3, n4 float64

	// Corner 0 contribution
	if t0 := 0.6 - x0*x0 - y0*y0 - z0*z0 - w0*w0; t0 > 0 {
		t0 *= t0
		hash := s.perm[ii+s.perm[jj+s.perm[kk+s.perm[ll]]]]
		n0 = t0 * t0 * grad4Inline(hash, x0, y0, z0, w0)
	}

	// Corner 1 contribution
	if t1 := 0.6 - x1*x1 - y1*y1 - z1*z1 - w1*w1; t1 > 0 {
		t1 *= t1
		hash := s.perm[ii+uint8(i1)+s.perm[jj+uint8(j1)+s.perm[kk+uint8(k1)+s.perm[ll+uint8(l1)]]]]
		n1 = t1 * t1 * grad4Inline(hash, x1, y1, z1, w1)
	}

	// Corner 2 contribution
	if t2 := 0.6 - x2*x2 - y2*y2 - z2*z2 - w2*w2; t2 > 0 {
		t2 *= t2
		hash := s.perm[ii+uint8(i2)+s.perm[jj+uint8(j2)+s.perm[kk+uint8(k2)+s.perm[ll+uint8(l2)]]]]
		n2 = t2 * t2 * grad4Inline(hash, x2, y2, z2, w2)
	}

	// Corner 3 contribution
	if t3 := 0.6 - x3*x3 - y3*y3 - z3*z3 - w3*w3; t3 > 0 {
		t3 *= t3
		hash := s.perm[ii+uint8(i3)+s.perm[jj+uint8(j3)+s.perm[kk+uint8(k3)+s.perm[ll+uint8(l3)]]]]
		n3 = t3 * t3 * grad4Inline(hash, x3, y3, z3, w3)
	}

	// Corner 4 contribution
	if t4 := 0.6 - x4*x4 - y4*y4 - z4*z4 - w4*w4; t4 > 0 {
		t4 *= t4
		hash := s.perm[ii+1+s.perm[jj+1+s.perm[kk+1+s.perm[ll+1]]]]
		n4 = t4 * t4 * grad4Inline(hash, x4, y4, z4, w4)
	}

	return 27.0 * (n0 + n1 + n2 + n3 + n4)
}

// Inlined 4D gradient evaluator selecting from 32 unit edge directions in 4D hypercube
func grad4Inline(hash uint8, x, y, z, w float64) float64 {
	h := hash & 31
	u := x
	if h >= 16 {
		u = y
	}
	v := y
	if h >= 8 {
		if h >= 24 {
			v = z
		} else {
			v = w
		}
	}
	s := z
	if h >= 4 {
		if h >= 20 {
			s = w
		} else {
			s = x
		}
	}

	if (h & 1) != 0 {
		u = -u
	}
	if (h & 2) != 0 {
		v = -v
	}
	if (h & 4) != 0 {
		s = -s
	}

	return u + v + s
}

// FBM generates fractal Brownian motion from 2D Simplex noise.
func (s *Simplex) FBM(x, y float64, octaves int, lacunarity, persistence float64) float64 {
	if octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < octaves; i++ {
		value += s.Noise2D(x*frequency, y*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= lacunarity
		amplitude *= persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

// FBM3D generates fractal Brownian motion from 3D Simplex noise.
func (s *Simplex) FBM3D(x, y, z float64, octaves int, lacunarity, persistence float64) float64 {
	if octaves <= 0 {
		return 0
	}

	value := 0.0
	amplitude := 1.0
	frequency := 1.0
	maxAmplitude := 0.0

	for i := 0; i < octaves; i++ {
		value += s.Noise3D(x*frequency, y*frequency, z*frequency) * amplitude
		maxAmplitude += amplitude
		frequency *= lacunarity
		amplitude *= persistence
	}

	if maxAmplitude == 0 {
		return 0
	}
	return value / maxAmplitude
}

// Normalize converts values from approximately [-1, 1] into [0, 1].
func Normalize(value float64) float64 {
	return (value + 1.0) * 0.5
}

// Clamp01 clamps a value to [0, 1].
func Clamp01(value float64) float64 {
	if value < 0 {
		return 0
	}
	if value > 1 {
		return 1
	}
	return value
}

func (s *Simplex) Noise5D(x, y, z, w, v float64) float64 {
	return 0
}

func (s *Simplex) Noise6D(x, y, z, w, v, u float64) float64 {
	return 0
}
