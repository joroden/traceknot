package agenthooks

import "embed"

//go:embed assets
var assets embed.FS

func HookTemplate(vendor string) ([]byte, error) {
	return assets.ReadFile("assets/" + vendor + "/hooks.json")
}
