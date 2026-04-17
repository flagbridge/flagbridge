package crypto

import (
	"strings"
	"testing"
)

func TestDeriveKey(t *testing.T) {
	key := DeriveKey("my-secret-passphrase")
	if len(key) != 32 {
		t.Errorf("expected 32-byte key, got %d", len(key))
	}

	key2 := DeriveKey("my-secret-passphrase")
	if string(key) != string(key2) {
		t.Error("same passphrase should produce same key")
	}

	key3 := DeriveKey("different-passphrase")
	if string(key) == string(key3) {
		t.Error("different passphrases should produce different keys")
	}
}

func TestEncryptDecrypt(t *testing.T) {
	key := DeriveKey("test-encryption-key")
	plaintext := "sk-ant-api03-very-secret-key-1234567890"

	encrypted, err := Encrypt(plaintext, key)
	if err != nil {
		t.Fatalf("encrypt failed: %v", err)
	}

	if encrypted == plaintext {
		t.Error("encrypted should differ from plaintext")
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("decrypt failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("expected %q, got %q", plaintext, decrypted)
	}
}

func TestEncryptProducesDifferentCiphertexts(t *testing.T) {
	key := DeriveKey("test-key")
	plaintext := "same-input"

	enc1, _ := Encrypt(plaintext, key)
	enc2, _ := Encrypt(plaintext, key)

	if enc1 == enc2 {
		t.Error("encrypting the same plaintext twice should produce different ciphertexts (random nonce)")
	}
}

func TestDecryptWrongKey(t *testing.T) {
	key1 := DeriveKey("correct-key")
	key2 := DeriveKey("wrong-key")

	encrypted, _ := Encrypt("secret", key1)

	_, err := Decrypt(encrypted, key2)
	if err == nil {
		t.Error("decrypting with wrong key should fail")
	}
}

func TestDecryptInvalidInput(t *testing.T) {
	key := DeriveKey("test-key")

	_, err := Decrypt("not-valid-base64!!!", key)
	if err == nil {
		t.Error("decrypting invalid base64 should fail")
	}

	_, err = Decrypt("YQ==", key)
	if err == nil {
		t.Error("decrypting too-short ciphertext should fail")
	}
}

func TestEncryptEmptyString(t *testing.T) {
	key := DeriveKey("test-key")

	encrypted, err := Encrypt("", key)
	if err != nil {
		t.Fatalf("encrypting empty string should work: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("decrypting empty string should work: %v", err)
	}

	if decrypted != "" {
		t.Errorf("expected empty string, got %q", decrypted)
	}
}

func TestEncryptLongKey(t *testing.T) {
	key := DeriveKey("test-key")
	longKey := strings.Repeat("x", 10000)

	encrypted, err := Encrypt(longKey, key)
	if err != nil {
		t.Fatalf("encrypting long string should work: %v", err)
	}

	decrypted, err := Decrypt(encrypted, key)
	if err != nil {
		t.Fatalf("decrypting long string should work: %v", err)
	}

	if decrypted != longKey {
		t.Error("round-trip failed for long string")
	}
}
