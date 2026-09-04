package noise

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateSolidNoise2DPNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const seed = 1337

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	sn := NewSolidNoise(seed)

	fileName := fmt.Sprintf("%s/solid_2d.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return sn.Noise2D(x, y)
	})

	if err != nil {
		t.Fatalf("Render failed for Solid 2D: %v", err)
	}

	t.Logf("Generated: %s", fileName)
}

func TestGenerateSolidNoise3DSlicePNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const (
		seed   = 42
		zSlice = 2.5
	)

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	sn := NewSolidNoise(seed)

	fileName := fmt.Sprintf("%s/solid_3d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return sn.Noise3D(x, y, zSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Solid 3D slice: %v", err)
	}

	t.Logf("Generated 3D slice: %s", fileName)
}

func TestGenerateSolidNoise4DSlicePNG(t *testing.T) {
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

	sn := NewSolidNoise(seed)

	fileName := fmt.Sprintf("%s/solid_4d_slice.png", outDir)

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return sn.Noise4D(x, y, zSlice, wSlice)
	})

	if err != nil {
		t.Fatalf("Render failed for Solid 4D slice: %v", err)
	}

	t.Logf("Generated 4D slice: %s", fileName)
}
