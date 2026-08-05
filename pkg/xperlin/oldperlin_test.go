package xperlin

import (
    "testing"
)

const (
    seed = 123
)

func Test_OldPerlinNoise1D(t *testing.T) {
    expected := 0.0
    p := NewOldPerlin(2, 2, 3, seed)
    noise := p.Noise1D(10)
    if noise != expected {
        t.Fail()
        t.Logf("Wrong node result: given: %f, expected: %f", noise, expected)
    }
}

func Test_OldPerlinNoise2D(t *testing.T) {
    expected := 0.0
    p := NewOldPerlin(2, 2, 3, seed)
    noise := p.Noise2D(10, 10)
    if noise != expected {
        t.Fail()
        t.Logf("Wrong node result: given: %f, expected: %f", noise, expected)
    }
}

func Test_OldPerlinNoise3D(t *testing.T) {
    expected := 0.0
    p := NewOldPerlin(2, 2, 3, seed)
    noise := p.Noise3D(10, 10, 10)
    if noise != expected {
        t.Fail()
        t.Logf("Wrong node result: given: %f, expected: %f", noise, expected)
    }
}

func Benchmark_OldPerlinNoise1D(b *testing.B) {
    p := NewOldPerlin(2, 2, 3, seed)
    for n := 0; n < b.N; n++ {
        p.Noise1D(10)
    }
}

func Benchmark_OldPerlinNoise2D(b *testing.B) {
    p := NewOldPerlin(2, 2, 3, seed)
    for n := 0; n < b.N; n++ {
        p.Noise2D(10, 10)
    }
}
func Benchmark_OldPerlinNoise3D(b *testing.B) {
    p := NewOldPerlin(2, 2, 3, seed)
    for n := 0; n < b.N; n++ {
        p.Noise3D(10, 10, 10)
    }
}
