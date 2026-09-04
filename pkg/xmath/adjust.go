package xmath

func FastFloor(x float64) int {
	xi := int(x)
	if x < float64(xi) {
		return xi - 1
	}
	return xi
}

func Quintic(t float64) float64 {
	return t * t * t * (t*(t*6.-15.) + 10.)
}

func QuinticField(t []float64) []float64 {
	r := make([]float64, len(t))
	for i := 0; i < len(t); i++ {
		r[i] = t[i] * t[i] * t[i] * (t[i]*(t[i]*6.-15.) + 10.)
	}
	return r
}

func Scurve(t float64) float64 {
	return t * t * (3. - 2.*t)
}

func ScurveField(t []float64) []float64 {
	r := make([]float64, len(t))
	for i := 0; i < len(t); i++ {
		r[i] = t[i] * t[i] * (3. - 2.*t[i])
	}
	return r
}

func Lerp(t, a, b float64) float64 {
	return a + t*(b-a)
}

func LerpField(t []float64, a, b float64) []float64 {
	r := make([]float64, len(t))
	for i := 0; i < len(t); i++ {
		r[i] = a + t[i]*(b-a)
	}
	return r
}

// SCurveContrast applies a smooth, piecewise quadratic S-curve contrast
// transformation to a noise value in the range [-1.0, 1.0].
func SCurveContrast(result float64) float64 {
	x := result + 1.0
	if x <= 1.0 {
		return (x * x) - 1.0
	}
	dx := x - 2.0
	return -(dx * dx) + 1.0
}
