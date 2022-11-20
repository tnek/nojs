package assets

import (
	"embed"

	"github.com/benbjohnson/hashfs"
)

//go:embed static/*
var Static embed.FS
var StaticFS = hashfs.NewFS(Static)
var StaticHTTPServ = hashfs.FileServer(StaticFS)
