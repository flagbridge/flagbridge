package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// GenerateAPIKey creates a new API key with the format: fb_sk_{scope}_{random_32_chars}
// Returns the full key (shown once) and the SHA-256 hash (stored in DB).
func GenerateAPIKey(scope string) (fullKey string, hash string, prefix string, err error) {
	randomBytes := make([]byte, 16)
	if _, err := rand.Read(randomBytes); err != nil {
		return "", "", "", fmt.Errorf("generating random bytes: %w", err)
	}

	randomHex := hex.EncodeToString(randomBytes)
	fullKey = fmt.Sprintf("fb_sk_%s_%s", scope, randomHex)
	prefix = fullKey[:20]

	h := sha256.Sum256([]byte(fullKey))
	hash = hex.EncodeToString(h[:])

	return fullKey, hash, prefix, nil
}

// HashAPIKey computes the SHA-256 hash of an API key for lookup.
func HashAPIKey(key string) string {
	h := sha256.Sum256([]byte(key))
	return hex.EncodeToString(h[:])
}
