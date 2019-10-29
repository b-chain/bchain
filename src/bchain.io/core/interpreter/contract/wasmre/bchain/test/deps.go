package test

//go:generate go-bindata -nometadata -pkg test -o bindata.go ../bchain.wasm
//go:generate gofmt -w -s bindata.go
