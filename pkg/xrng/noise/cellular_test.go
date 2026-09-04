package noise

import (
	"fmt"
	"os"
	"testing"
)

func TestGenerateCellularNoisePNGs(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const seed = 1337

	returnTypes := []struct {
		name       string
		returnType CellularReturnType
	}{
		{"DISTANCE", DISTANCE},
		{"DISTANCE_2", DISTANCE_2},
		{"DISTANCE_2_ADD", DISTANCE_2_ADD},
		{"DISTANCE_2_SUB", DISTANCE_2_SUB},
		{"DISTANCE_2_MUL", DISTANCE_2_MUL},
		{"DISTANCE_2_DIV", DISTANCE_2_DIV},
		{"NOISE_LOOKUP", NOISE_LOOKUP},
		{"CELL_VALUE", CELL_VALUE},
	}

	outDir := "test_output"
	if err := os.MkdirAll(outDir, 0755); err != nil {
		t.Fatalf("Failed to create output directory: %v", err)
	}

	for _, rt := range returnTypes {
		t.Run(rt.name, func(t *testing.T) {
			cn := NewCellularNoise(seed)
			cn.Variant = EUCLIDEAN
			cn.ReturnType = rt.returnType

			fileName := fmt.Sprintf("%s/cellular_2d_%s.png", outDir, rt.name)

			err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
				return cn.Noise2D(x, y)
			})

			if err != nil {
				t.Fatalf("Render failed for %s: %v", rt.name, err)
			}

			t.Logf("Generated: %s", fileName)
		})
	}
}

func TestGenerateCellular3DSlicePNG(t *testing.T) {
	renderCfg := NewDefaultRenderConfig()
	const (
		seed   = 42
		zSlice = 2.5
	)

	cn := NewCellularNoise(seed)
	cn.Variant = EUCLIDEAN
	cn.ReturnType = DISTANCE_2_SUB

	fileName := "test_output/cellular_3d_DISTANCE_2_SUB.png"

	err := RenderNoiseToPNG(fileName, renderCfg, func(x, y float64) float64 {
		return cn.Noise3D(x, y, zSlice)
	})

	if err != nil {
		t.Fatalf("Render failed: %v", err)
	}

	t.Logf("Generated 3D slice: %s", fileName)
}
