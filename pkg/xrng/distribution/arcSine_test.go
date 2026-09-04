package distribution

import (
	"math"
	"math/rand/v2"
	"testing"
)

func TestArcsineDistribution(t *testing.T) {
	gen := rand.NewPCG(42, 108) // Seeded generator for reproducible tests
	alpha := 0.0
	beta := 10.0

	// Test Construction & Parameter Validation
	dist, err := NewArcsineDistribution(gen, alpha, beta)
	if err != nil {
		t.Fatalf("Failed to instantiate ArcsineDistribution: %v", err)
	}

	t.Run("Getters & Identifiers", func(t *testing.T) {
		if tag := dist.GetTag(); tag != "Arcsine" {
			t.Errorf("Expected tag 'Arcsine', got '%s'", tag)
		}
		if a := dist.GetParameterA(); a != alpha {
			t.Errorf("Expected Parameter A (Alpha) = %f, got %f", alpha, a)
		}
		if b := dist.GetParameterB(); b != beta {
			t.Errorf("Expected Parameter B (Beta) = %f, got %f", beta, b)
		}
		if c := dist.GetParameterC(); !math.IsNaN(c) {
			t.Errorf("Expected Parameter C to be NaN, got %f", c)
		}
	})

	t.Run("Statistical Properties", func(t *testing.T) {
		expectedMean := 0.5 * (alpha + beta)
		if mean := dist.GetMean(); mean != expectedMean {
			t.Errorf("Expected Mean = %f, got %f", expectedMean, mean)
		}

		if median := dist.GetMedian(); median != expectedMean {
			t.Errorf("Expected Median = %f, got %f", expectedMean, median)
		}

		expectedVariance := (beta - alpha) * (beta - alpha) * 0.125
		if v := dist.GetVariance(); v != expectedVariance {
			t.Errorf("Expected Variance = %f, got %f", expectedVariance, v)
		}

		modes := dist.GetMode()
		if len(modes) != 2 || modes[0] != alpha || modes[1] != beta {
			t.Errorf("Expected Mode = [%f, %f], got %v", alpha, beta, modes)
		}

		if min := dist.GetMinimum(); min != alpha {
			t.Errorf("Expected Minimum = %f, got %f", alpha, min)
		}

		if max := dist.GetMaximum(); max != beta {
			t.Errorf("Expected Maximum = %f, got %f", beta, max)
		}
	})

	t.Run("Sampling Bounds Verification", func(t *testing.T) {
		samples := 10_000
		for i := 0; i < samples; i++ {
			sample := dist.NextDouble()
			if sample < alpha || sample > beta {
				t.Fatalf("Sample %f out of bounds [%f, %f]", sample, alpha, beta)
			}
		}
	})

	t.Run("Invalid Parameters", func(t *testing.T) {
		// Alpha cannot be greater than or equal to Beta
		_, err := NewArcsineDistribution(gen, 10.0, 5.0)
		if err == nil {
			t.Error("Expected error when alpha >= beta, got nil")
		}

		if dist.SetParameters(5.0, 5.0, 0.0) {
			t.Error("SetParameters should return false when alpha == beta")
		}
	})

	t.Run("Copy & Serialization", func(t *testing.T) {
		copyDist := dist.Copy()
		if copyDist.GetTag() != dist.GetTag() {
			t.Errorf("Copied distribution tag mismatch")
		}

		serialized := dist.StringSerialize()
		expectedSerialized := "Arcsine:0.000000:10.000000"
		if serialized != expectedSerialized {
			t.Errorf("Expected serialized string '%s', got '%s'", expectedSerialized, serialized)
		}

		deserialized, err := dist.StringDeserialize(serialized)
		if err != nil {
			t.Fatalf("Deserialization failed: %v", err)
		}
		if deserialized.GetParameterA() != alpha || deserialized.GetParameterB() != beta {
			t.Errorf("Deserialized parameters mismatch original")
		}
	})
}
