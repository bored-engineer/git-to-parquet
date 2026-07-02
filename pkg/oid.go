package git

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
)

// OID represents a Git object ID, which is a SHA-1 hash.
// I don't have any plans to support SHA-256 unless it sees wide adoption.
type OID [sha1.Size]byte

// String implements fmt.Stringer
func (o OID) String() string {
	return hex.EncodeToString(o[:])
}

// UnmarshalText implements the encoding.TextUnmarshaler interface for a hex-encoded OID.
func (o *OID) UnmarshalText(text []byte) error {
	if expected := hex.EncodedLen(sha1.Size); len(text) != expected {
		return fmt.Errorf("invalid OID length: expected %d bytes, got %d", expected, len(text))
	}
	if _, err := hex.Decode(o[:], text); err != nil {
		return fmt.Errorf("hex.Decode failed: %w", err)
	}
	return nil
}
