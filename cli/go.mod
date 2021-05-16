module github.com/mksign/nsd-hack/cli

go 1.16

replace github.com/KMSIGN/nsd-hack/go-file-handler => ../go-file-handler

require (
	github.com/KMSIGN/nsd-hack/go-file-handler v0.0.0-20210516064302-884ae121d8c0
	github.com/cpuguy83/go-md2man/v2 v2.0.0 // indirect
	github.com/russross/blackfriday/v2 v2.1.0 // indirect
	github.com/urfave/cli/v2 v2.3.0
)
