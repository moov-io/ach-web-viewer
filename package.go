// generated-from:00b592a39a8330994d5ea8583e2b41307b36cda123cfeb87c268cbe9252aedc4 DO NOT REMOVE, DO UPDATE

// stub to get pkger to work
package achwebviewer

import (
	"github.com/markbates/pkger"
)

// Add in all includes that pkger should embed into the application here
var _ = pkger.Include("/configs/config.default.yml")
var _ = pkger.Include("/migrations/")

// Load our HTML templates
var _ = pkger.Include("/webui/style.css")
var _ = pkger.Include("/webui/index.html.tpl")
var _ = pkger.Include("/webui/file.html.tpl")
