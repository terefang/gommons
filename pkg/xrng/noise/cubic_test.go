package noise

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateCubicNoise2DPNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const seed = 1337

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	cn := NewCubicNoise(seed)

	fileName := fmt.Sprintf("%s/cubic_2d.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return cn.Noise2D(x, y)
	})

	if err != nil {
		t.Fatalf("Render failed for Cubic 2D: %v", err)
	}

	t.Logf("Generated: %s", fileName)
}

func TestGenerateCubicNoise3DSlicePNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const (
		seed   = 42
		zSlice = 2.5
	)

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	cn := NewCubicNoise(seed)

	fileName := fmt.Sprintf("%s/cubic_3d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return cn.Noise3D(x, y, zSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Cubic 3D slice: %v", err)
	}

	t.Logf("Generated 3D slice: %s", fileName)
}

func TestGenerateCubicNoise4DSlicePNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const (
		seed   = 101
		zSlice = 1.0
		wSlice = 0.5
	)

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	cn := NewCubicNoise(seed)

	fileName := fmt.Sprintf("%s/cubic_4d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return cn.Noise4D(x, y, zSlice, wSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Cubic 4D slice: %v", err)
	}

	t.Logf("Generated 4D slice: %s", fileName)
}
