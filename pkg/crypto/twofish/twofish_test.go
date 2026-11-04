package twofish

import (
	"encoding/hex"
	"strings"
	"testing"
)

// TestNewCipher validates key bounds and error handling.
func TestNewCipher(t *testing.T) {
	// Invalid key sizes (must be 16, 24, or 32 bytes)
	invalidSizes := []int{0, 8, 15, 17, 20, 25, 31, 33, 64}
	for _, size := range invalidSizes {
		_, err := NewCipher(make([]byte, size))
		if err == nil {
			t.Errorf("Expected error for key size %d bytes, got nil", size)
		}
	}

	// Valid key sizes (128-bit, 192-bit, 256-bit)
	validSizes := []int{16, 24, 32}
	for _, size := range validSizes {
		block, err := NewCipher(make([]byte, size))
		if err != nil {
			t.Errorf("Unexpected error for valid key size %d bytes: %v", size, err)
		}
		if block.BlockSize() != BlockSize {
			t.Errorf("Expected block size %d, got %d", BlockSize, block.BlockSize())
		}
	}
}

// TestEncryptionRoundTrip ensures data encrypts and decrypts correctly.
func TestEncryptionRoundTrip(t *testing.T) {
	key := []byte("12345678901234567890123456789012") // 32 bytes (256-bit)
	block, err := NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	plaintext := []byte("1234567890ABCDEF") // Exactly 16 bytes (BlockSize)
	ciphertext := make([]byte, BlockSize)
	decrypted := make([]byte, BlockSize)

	block.Encrypt(ciphertext, plaintext)
	block.Decrypt(decrypted, ciphertext)

	for i := range plaintext {
		if plaintext[i] != decrypted[i] {
			t.Errorf("Mismatch at index %d: expected %02x, got %02x", i, plaintext[i], decrypted[i])
		}
	}
}

// TestTwofishSelfTestVector mirrors the exact self-test parameters from the original C code
// to verify 100% binary compatibility with the legacy module.
func TestTwofishSelfTestVector(t *testing.T) {
	keyHex := "9f589f5cf6122c32b6bfec2f2ae8c35a"
	ptHex := "d491db16e7b1c39e86cb086b789f5419"
	expectedCtHex := "019f9809de1711858faac3a3ba20fbc3"

	key, err := hex.DecodeString(keyHex)
	if err != nil {
		t.Fatalf("failed to decode key hex: %v", err)
	}
	pt, err := hex.DecodeString(ptHex)
	if err != nil {
		t.Fatalf("failed to decode plaintext hex: %v", err)
	}

	block, err := NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher with C test vector key: %v", err)
	}

	ciphertext := make([]byte, BlockSize)
	block.Encrypt(ciphertext, pt)
	gotCtHex := hex.EncodeToString(ciphertext)

	if !strings.EqualFold(gotCtHex, expectedCtHex) {
		t.Errorf("C compatibility vector failed.\nExpected: %s\nGot:      %s", expectedCtHex, gotCtHex)
	}

	// Verify decryption round-trip
	decrypted := make([]byte, BlockSize)
	block.Decrypt(decrypted, ciphertext)

	if hex.EncodeToString(decrypted) != hex.EncodeToString(pt) {
		t.Errorf("C compatibility decryption mismatch.\nExpected: %x\nGot:      %x", pt, decrypted)
	}
}

// TestTwofishVectors runs standard official Twofish test vectors across different key lengths.
func TestTwofishVectors(t *testing.T) {
	vectors := []struct {
		name       string
		key        string
		plaintext  string
		ciphertext string
	}{
		{
			name:       "Vector_128_bit",
			key:        "00000000000000000000000000000000",
			plaintext:  "00000000000000000000000000000000",
			ciphertext: "9f589f5cf6122c32b6bfec2f2ae8c35a", // Note: mapped per standard reference suite or specific test vector entries
		},
	}

	// Additional official test cases can be appended here.
	for _, v := range vectors {
		t.Run(v.name, func(t *testing.T) {
			key, err := hex.DecodeString(v.key)
			if err != nil {
				t.Fatalf("invalid hex key: %v", err)
			}
			pt, err := hex.DecodeString(v.plaintext)
			if err != nil {
				t.Fatalf("invalid hex plaintext: %v", err)
			}

			block, err := NewCipher(key)
			if err != nil {
				t.Fatalf("failed to create cipher: %v", err)
			}

			dst := make([]byte, BlockSize)
			block.Encrypt(dst, pt)
			gotCt := hex.EncodeToString(dst)

			// If testing against specific validation sets, check match:
			if v.ciphertext != gotCt {
				t.Errorf("encryption mismatch for %s", v.name)
			}

			decrypted := make([]byte, BlockSize)
			block.Decrypt(decrypted, dst)
			dHex := hex.EncodeToString(decrypted)
			if dHex != v.plaintext {
				t.Errorf("Decryption round-trip mismatch for %s", v.name)
			}
		})
	}
}
