package noise

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateSimplexNoise2DPNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const seed = 1337

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	sn := NewSimplex(seed)

	fileName := fmt.Sprintf("%s/simplex_2d.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return sn.Noise2D(x, y)
	})

	if err != nil {
		t.Fatalf("Render failed for Simplex 2D: %v", err)
	}

	t.Logf("Generated: %s", fileName)
}

func TestGenerateSimplex3DSlicePNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const (
		seed   = 42
		zSlice = 2.5
	)

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	sn := NewSimplex(seed)

	fileName := fmt.Sprintf("%s/simplex_3d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return sn.Noise3D(x, y, zSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Simplex 3D slice: %v", err)
	}

	t.Logf("Generated 3D slice: %s", fileName)
}

func TestGenerateSimplex4DSlicePNG(t *testing.T) {
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

	sn := NewSimplex(seed)

	fileName := fmt.Sprintf("%s/simplex_4d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return sn.Noise4D(x, y, zSlice, wSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Simplex 4D slice: %v", err)
	}

	t.Logf("Generated 4D slice: %s", fileName)
}
