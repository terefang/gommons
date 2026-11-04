package blowfish

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
)

// TestNewCipher validates key bounds and error handling.
func TestNewCipher(t *testing.T) {
	// Empty key should fail
	_, err := NewCipher([]byte{})
	if err == nil {
		t.Error("Expected error for empty key, got nil")
	}

	// Key longer than 56 bytes (448 bits) should fail
	_, err = NewCipher(make([]byte, 57))
	if err == nil {
		t.Error("Expected error for key size > 56 bytes, got nil")
	}

	// Valid key should succeed
	block, err := NewCipher([]byte("valid-key-length"))
	if err != nil {
		t.Errorf("Unexpected error for valid key: %v", err)
	}
	if block.BlockSize() != BlockSize {
		t.Errorf("Expected block size %d, got %d", BlockSize, block.BlockSize())
	}
}

// TestEncryptionRoundTrip ensures data encrypts and decrypts correctly.
func TestEncryptionRoundTrip(t *testing.T) {
	key := []byte("my-secure-key-12345")
	block, err := NewCipher(key)
	if err != nil {
		t.Fatalf("Failed to create cipher: %v", err)
	}

	plaintext := []byte("ABCDEFGH") // Exactly 8 bytes
	ciphertext := make([]byte, 8)
	decrypted := make([]byte, 8)

	block.Encrypt(ciphertext, plaintext)
	block.Decrypt(decrypted, ciphertext)

	for i := range plaintext {
		if plaintext[i] != decrypted[i] {
			t.Errorf("Mismatch at index %d: expected %02x, got %02x", i, plaintext[i], decrypted[i])
		}
	}
}

// TestCCompatibilityVector mirrors the exact self-test parameters from the original C code
// to verify 100% binary compatibility with the legacy module.
func TestCCompatibilityVector(t *testing.T) {
	keySize := 56
	keyword := make([]byte, keySize)
	for j := 0; j < keySize; j++ {
		keyword[j] = byte((j*2 + 10) % 256)
	}

	blockSize := 8
	plaintext := make([]byte, blockSize)
	for j := 0; j < blockSize; j++ {
		plaintext[j] = byte(j % 256)
	}

	block, err := NewCipher(keyword)
	if err != nil {
		t.Fatalf("Failed to initialize cipher with C test vector key: %v", err)
	}

	ciphertext := make([]byte, blockSize)
	block.Encrypt(ciphertext, plaintext)

	cipherHex := hex.EncodeToString(ciphertext)
	expectedHex := "c8c033bc57874d74"

	if cipherHex != expectedHex {
		t.Errorf("C compatibility vector failed.\nExpected: %s\nGot:      %s", expectedHex, cipherHex)
	}

	// Verify that decrypting the ciphertext returns the original plaintext
	decrypted := make([]byte, blockSize)
	block.Decrypt(decrypted, ciphertext)

	for i := range plaintext {
		if plaintext[i] != decrypted[i] {
			t.Errorf("Decryption mismatch at byte %d: expected %02x, got %02x", i, plaintext[i], decrypted[i])
		}
	}
}

// TestBlowfishVectors runs standard official test vectors for Blowfish.
func TestBlowfishVectors(t *testing.T) {
	vectors := []struct {
		key        string
		plaintext  string
		ciphertext string
	}{
		{"0000000000000000", "0000000000000000", "4EF997456198DD78"},
		{"FFFFFFFFFFFFFFFF", "FFFFFFFFFFFFFFFF", "51866FD5B85ECB8A"},
		{"3000000000000000", "1000000000000001", "7D856F9A613063F2"},
		{"1111111111111111", "1111111111111111", "2466DD878B963C9D"},
		{"0123456789ABCDEF", "1111111111111111", "61F9C3802281B096"},
	}

	for i, v := range vectors {
		t.Run(fmt.Sprintf("Vector_%d", i+1), func(t *testing.T) {
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

			// Test Encryption
			dst := make([]byte, 8)
			block.Encrypt(dst, pt)
			gotCt := hex.EncodeToString(dst)

			if !strings.EqualFold(gotCt, v.ciphertext) {
				t.Errorf("Encryption mismatch:\nExpected: %s\nGot:      %s", v.ciphertext, gotCt)
			}

			// Test Decryption Round-trip
			decrypted := make([]byte, 8)
			block.Decrypt(decrypted, dst)

			if hex.EncodeToString(decrypted) != hex.EncodeToString(pt) {
				t.Errorf("Decryption mismatch:\nExpected: %x\nGot:      %x", pt, decrypted)
			}
		})
	}
}
