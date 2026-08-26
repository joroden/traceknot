package content

import (
	"fmt"
	"regexp"
)

var base64Token = regexp.MustCompile(`[A-Za-z0-9+/]{200,}={0,2}`)

func StripBinaryBlobs(text string) string {
	return base64Token.ReplaceAllStringFunc(text, func(match string) string {
		return fmt.Sprintf("[binary data omitted: ~%d bytes, base64]", len(match))
	})
}
