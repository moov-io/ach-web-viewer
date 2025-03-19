package webui

import (
	"embed"
)

//go:embed static/*.css static/*.js *.html.tmpl
var WebRoot embed.FS
