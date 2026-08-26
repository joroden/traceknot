package content

import "crypto/sha256"

const minDedupLen = 200

type Deduper struct {
	seen map[[32]byte]string
}

func NewDeduper() *Deduper {
	return &Deduper{seen: make(map[[32]byte]string)}
}

func (d *Deduper) Check(text, label string) (firstLabel string, isDuplicate bool) {
	if len(text) < minDedupLen {
		return "", false
	}
	sum := sha256.Sum256([]byte(text))
	if existing, found := d.seen[sum]; found {
		return existing, true
	}
	d.seen[sum] = label
	return "", false
}
