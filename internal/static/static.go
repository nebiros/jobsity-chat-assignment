package static

import "embed"

var (
	//go:embed all:*.gohtml
	FS embed.FS
)
