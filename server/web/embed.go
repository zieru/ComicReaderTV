package web

import "embed"

//go:embed *.html *.js
var Files embed.FS
