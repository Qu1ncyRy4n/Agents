// Package mint allocates collision-free proquint handles within a repository.
//
// Copied and adapted from cdint-grid tools/mint-handle at commit
// b73df41f3e43ecb38f2133a5ad8ad91cc3cfceb3 with user-authorized internal
// transfer, 2026-09-10.
package mint

import (
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"time"
)

// Mint returns a fresh handle of width 1 or 2 that is not present in corpus.
// It draws entropy from time.Now().UnixNano(), folds through SHA-256, and
// retries on collision. A non-zero seed supplies the first attempt.
//
// In dry-run mode no retries are performed; a collision is reported instead.
func Mint(width int, corpus map[string]string, seed int64, dryRun bool) (string, error) {
	const maxAttempts = 1_000_000
	for attempt := 0; attempt < maxAttempts; attempt++ {
		var ns int64
		if attempt == 0 && seed != 0 {
			ns = seed
		} else {
			ns = time.Now().UnixNano()
		}
		var buf [8]byte
		binary.BigEndian.PutUint64(buf[:], uint64(ns))
		sum := sha256.Sum256(buf[:])

		var candidate string
		switch width {
		case 1:
			candidate = proquint1FromBytes(sum[:2])
		case 2:
			candidate = proquint2FromBytes(sum[:4])
		default:
			return "", fmt.Errorf("-w must be 1 or 2, got %d", width)
		}

		if _, taken := corpus[candidate]; !taken {
			return candidate, nil
		}
		if dryRun {
			return "", fmt.Errorf("dry-run: first attempt %q collides with corpus", candidate)
		}

		// Advance coarse clocks enough that the next seed differs.
		time.Sleep(time.Microsecond)
	}
	return "", fmt.Errorf("exhausted %d attempts; corpus may be saturated (%d handles in use)", maxAttempts, len(corpus))
}

const (
	proquintCons = "bdfghjklmnprstvz"
	proquintVows = "aiou"
)

// uint16ToProquint encodes one 16-bit value as a CVCVC proquint. The bit
// layout follows Wilkerson's original proquint paper exactly.
func uint16ToProquint(n uint16) string {
	return string([]byte{
		proquintCons[(n>>12)&0x0f],
		proquintVows[(n>>10)&0x03],
		proquintCons[(n>>6)&0x0f],
		proquintVows[(n>>4)&0x03],
		proquintCons[n&0x0f],
	})
}

func proquint1FromBytes(b []byte) string {
	return uint16ToProquint(uint16(b[0])<<8 | uint16(b[1]))
}

func proquint2FromBytes(b []byte) string {
	return proquint1FromBytes(b[:2]) + "-" + proquint1FromBytes(b[2:4])
}
