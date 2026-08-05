package xmath

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
