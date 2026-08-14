package agentenv

import (
	"fmt"

	"embed"
)

//go:embed assets
var assets embed.FS

func loadSnippet(path string) ([]byte, error) {
	content, err := assets.ReadFile("assets/" + path)
	if err != nil {
		return nil, fmt.Errorf("load %s: %w", path, err)
	}
	return content, nil
}
