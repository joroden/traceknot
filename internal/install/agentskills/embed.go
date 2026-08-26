package agentskills

import (
	"embed"
	"io/fs"
)

//go:embed assets/traceknot-analyze
var assets embed.FS

func SkillFiles() (fs.FS, error) {
	return fs.Sub(assets, "assets/traceknot-analyze")
}
