#!/bin/sh
rm go.mod
rm go.sum

go mod init github.com/tnek/notes-site
go mod tidy
go build
