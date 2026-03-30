package handler

import (
	"crypto/rand"
	"encoding/hex"
)

// generateID produces a short random hex ID (8 bytes = 16 hex chars).
func generateID() string {
	b := make([]byte, 8)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}
