package deps


//go:generate loadWasmFile ../pledge/pledge.wasm pledge.json
//go:generate loadWasmFile ../system/system.wasm system.json
//go:generate loadWasmFile ../bchain/bchain.wasm bchain.json

//go:generate go-bindata -nometadata -pkg deps -o bindata.go pledge.json system.json  bchain.json
