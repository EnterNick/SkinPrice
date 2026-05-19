package crypto

import (
	"encoding/base64"
	"testing"
)

func TestTokenCipherRoundTrip(t *testing.T) {
	cipher, err := NewTokenCipher([]byte("12345678901234567890123456789012"))
	if err != nil {
		t.Fatalf("NewTokenCipher() error = %v", err)
	}
	orig := "super-secret-token"
	encrypted, err := cipher.Encrypt(orig)
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	decrypted, err := cipher.Decrypt(encrypted)
	if err != nil {
		t.Fatalf("Decrypt() error = %v", err)
	}
	if decrypted != orig {
		t.Fatalf("Decrypt() = %q, want %q", decrypted, orig)
	}
}

func TestTokenCipherInvalidKey(t *testing.T) {
	if _, err := NewTokenCipher([]byte("short")); err == nil {
		t.Fatal("NewTokenCipher() expected error for invalid key length")
	}
}

func TestTokenCipherFromBase64(t *testing.T) {
	key := base64.StdEncoding.EncodeToString([]byte("12345678901234567890123456789012"))
	if _, err := NewTokenCipherFromBase64(key); err != nil {
		t.Fatalf("NewTokenCipherFromBase64() error = %v", err)
	}
}

func TestTokenCipherCorruptedCiphertext(t *testing.T) {
	cipher, err := NewTokenCipher([]byte("12345678901234567890123456789012"))
	if err != nil {
		t.Fatalf("NewTokenCipher() error = %v", err)
	}
	encrypted, err := cipher.Encrypt("payload")
	if err != nil {
		t.Fatalf("Encrypt() error = %v", err)
	}
	corrupted := encrypted[:len(encrypted)-2] + "zz"
	if _, err = cipher.Decrypt(corrupted); err == nil {
		t.Fatal("Decrypt() expected error for corrupted ciphertext")
	}
}
