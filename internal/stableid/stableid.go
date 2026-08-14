package stableid

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
)

func From(prefix string, parts ...string) string {
	digest := sha256.Sum256([]byte(strings.Join(parts, "\n")))
	hexDigest := hex.EncodeToString(digest[:])
	return prefix + "_" + hexDigest[:24]
}
