package ui

import "embed"

//go:embed "templates"
var Files embed.FS

//go:embed "assets"
var StaticFiles embed.FS
