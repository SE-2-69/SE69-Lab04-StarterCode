package utils

import (
	"crypto/rand"
	"encoding/hex"
)

// random, unguessable id (32 bytes) for session cookies
func RandomID() string {
	b := make([]byte, 32)
	rand.Read(b)
	return hex.EncodeToString(b)
}
