package achwebviewer

import (
	"embed"
)

//go:embed configs/config.default.yml
var ConfigDefaults embed.FS
