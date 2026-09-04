package noise

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateValueNoise2DPNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const seed = 1337

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	vn := NewValueNoise(seed)

	fileName := fmt.Sprintf("%s/value_2d.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return vn.Noise2D(x, y)
	})

	if err != nil {
		t.Fatalf("Render failed for Value 2D: %v", err)
	}

	t.Logf("Generated: %s", fileName)
}

func TestGenerateValueNoise3DSlicePNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const (
		seed   = 42
		zSlice = 2.5
	)

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	vn := NewValueNoise(seed)

	fileName := fmt.Sprintf("%s/value_3d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return vn.Noise3D(x, y, zSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Value 3D slice: %v", err)
	}

	t.Logf("Generated 3D slice: %s", fileName)
}

func TestGenerateValueNoise4DSlicePNG(t *testing.T) {
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

	vn := NewValueNoise(seed)

	fileName := fmt.Sprintf("%s/value_4d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return vn.Noise4D(x, y, zSlice, wSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Value 4D slice: %v", err)
	}

	t.Logf("Generated 4D slice: %s", fileName)
}
