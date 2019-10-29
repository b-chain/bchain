// Code generated by go-bindata.
// sources:
// ../pledge.wasm
// ../../bchain/bchain.wasm
// ../../system/system.wasm
// DO NOT EDIT!

package test

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func bindataRead(data []byte, name string) ([]byte, error) {
	gz, err := gzip.NewReader(bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}

	var buf bytes.Buffer
	_, err = io.Copy(&buf, gz)
	clErr := gz.Close()

	if err != nil {
		return nil, fmt.Errorf("Read %q: %v", name, err)
	}
	if clErr != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

type asset struct {
	bytes []byte
	info  os.FileInfo
}

type bindataFileInfo struct {
	name    string
	size    int64
	mode    os.FileMode
	modTime time.Time
}

func (fi bindataFileInfo) Name() string {
	return fi.name
}
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}
func (fi bindataFileInfo) IsDir() bool {
	return false
}
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _PledgeWasm = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xa4\x54\xbf\x6f\xdb\x46\x14\x7e\x77\x47\x53\x94\x69\xb9\xd7\x5f\x80\xd1\x16\xe8\xab\xbb\xb4\xf0\x22\x59\xf6\xc9\xed\xe4\x2b\xe0\xb5\x35\x8a\x02\x2d\xbc\x48\x94\x78\xb2\x45\x53\xa4\x41\x9e\x5b\x1b\x68\x4d\x19\x2d\x0a\xb5\x93\x87\x0c\x19\x12\xc0\x7b\x32\xe4\x4f\xf0\x92\x29\x4b\xc6\x8c\xfe\x0b\x02\xff\x01\x41\x10\x1c\x49\x09\xb2\x35\x24\x40\x00\x41\x7a\x3f\xbe\xf7\xdd\xfb\x3e\x1e\x05\x5e\x3a\xac\x01\x00\x79\x3c\x1a\x8d\xc0\xe9\x90\x0c\x3a\x24\x23\x59\x87\x65\x59\x06\xf9\x37\xc9\x3a\xd4\xc4\x56\x96\x97\x68\x76\x06\x1d\x72\x06\xf4\xe9\xf9\x68\x04\x2e\x53\xd1\xef\x1f\x7b\x3d\x3d\x88\xa3\x76\xcf\x0b\xc3\x5f\x07\xfa\x60\xd7\x4b\x3c\xa0\xa6\x53\x2b\x3b\xa9\x8a\x7c\x95\x00\x98\x9a\xed\xa5\xa9\x4a\x34\x58\x26\xe1\xbd\x38\xd2\x89\xd7\xd3\x6d\xcf\xf7\x13\x95\xa6\x05\xe6\xd3\x69\x79\x9e\xd3\xf6\xbb\xed\x7d\xa5\x81\x4d\x92\x54\x69\x58\x30\x49\x35\x55\xfa\x67\x95\x1e\x87\x25\x79\x35\xd5\xc9\x7a\x18\xff\x31\x39\xd9\x49\x75\x12\xc4\x83\x68\x32\x9b\xea\x24\x54\x11\x10\xf6\x5f\x2e\xde\xb6\x2d\xab\x52\x01\xb0\xfe\x31\x29\x39\x02\x58\xf8\x3b\x8f\x80\x54\x1e\x19\xb1\x55\x7b\xa8\x86\x71\x72\x4a\x61\xa5\xbd\xf7\x63\xa3\x7e\x14\x2a\x7f\x5f\xed\xc6\x71\x28\x8a\x70\xe7\x14\xdc\xb9\x56\xa2\x7c\xa5\x86\x3b\xa7\xb0\xf4\xf9\x9d\xd6\x56\x11\xfe\xd4\xdf\xd9\xed\x41\xed\xcb\x3b\xdd\x46\x63\xda\x3e\xd1\x06\xb1\x6c\x17\x05\xf8\xc0\x2e\x38\x81\x3b\x13\x08\x7c\xe8\xce\xa0\xe1\xa3\xea\xb9\xd9\x1c\x16\x47\xff\x1a\x65\xaf\xcd\xf6\x84\x66\x74\x5b\x82\x84\x6f\xa8\x25\x2f\xe8\xe1\x2a\x13\xd4\x42\x86\xa4\xc5\xc6\x04\x99\x7c\x0e\x01\x37\xbf\x18\x70\x56\xa6\x0e\xb2\x1f\xa0\xc5\x38\x32\x79\x41\x02\xc9\x8b\x72\x91\xf1\xea\x2a\x45\x26\x79\xc0\x17\x3e\x43\x86\xec\x5b\xc6\x91\xfc\x39\x05\x63\xd1\x94\x0e\xb7\x91\x49\x4b\xd0\x0b\xc3\xbd\x26\xe8\x7d\x62\xd0\x12\x03\x41\xef\x99\x92\x23\xe8\x4a\x5e\x19\x93\x40\xd0\x4f\x90\x49\x22\xa8\x21\xe1\xf9\x7a\x66\x14\x24\xe7\x8b\x5f\x19\x42\x10\x74\x09\x19\x52\x69\xb0\x8e\x19\x13\xf4\x8b\x7c\x9a\x07\x82\x3e\x20\xf2\x0a\x24\x96\xfb\x81\x04\x13\xd1\x40\x50\xcb\x75\x5f\xcd\x59\x70\x4d\x66\x2d\xf8\xad\x50\x5e\x1a\x30\x23\xfc\x1a\x0a\xe1\x18\x14\xc9\x6d\xdd\xb9\xea\x3d\x79\x03\x9c\x4e\x4d\xf8\x6b\x3a\x38\x67\x42\xa7\xf0\xe0\x60\x6a\x81\x3f\xeb\xc0\x0b\xb8\xed\x80\x91\x8e\x4c\xd6\x05\xb5\x64\xfd\x9d\x1c\x08\x27\x06\x98\x45\xad\xdc\x80\x6b\x52\x18\xf0\x2c\xbf\xca\x64\xf6\x0e\x90\xc3\x55\x6a\x0c\xa0\x46\xeb\x98\x20\xe1\x0e\x52\xc9\x91\x20\xe5\x55\xa4\xe6\x89\x18\x91\x79\x20\x1d\x5e\x91\x80\xd4\x58\x9b\xf3\xbd\x7c\x7f\xbe\x09\x50\x5e\x9a\x0f\x5f\x7c\xdb\x91\xff\x97\x47\x22\x41\xe0\x77\xd2\x25\x77\x9c\xdf\x78\x04\x04\x5e\x9b\x4d\x96\x5d\xf7\xa1\x79\xf8\x15\x90\x96\x6b\x3d\xd9\x06\x00\xc9\xdd\x4a\xf9\x2e\x81\x44\xb7\xaa\x13\x2f\x4a\xfb\xe6\x6f\x42\xd6\xdd\x4a\xf9\x72\x81\xbc\x02\x77\xad\x7e\xd2\x68\x74\x9b\x4d\xd5\xec\xab\xd6\xba\xbf\xd9\x14\xde\x46\xd3\xf7\x7b\xdd\xad\x8d\x7e\xdd\xef\xb6\x5a\xa2\xb9\xb5\x5e\x17\xdf\xf5\x36\x1b\x00\xf2\x06\xdc\xaf\x13\xa5\xd4\xf0\x7b\x2c\xe8\xd1\x1b\xc6\xc7\x91\xc6\x41\x8a\x51\xac\x51\x45\xf1\xf1\xfe\x01\x80\xbc\x24\x6e\xad\x5d\x40\x7e\x89\xb5\x17\xc2\x9b\x00\x00\x00\xff\xff\xa1\x4d\x73\x18\x88\x05\x00\x00")

func PledgeWasmBytes() ([]byte, error) {
	return bindataRead(
		_PledgeWasm,
		"../pledge.wasm",
	)
}

func PledgeWasm() (*asset, error) {
	bytes, err := PledgeWasmBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../pledge.wasm", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _BchainBchainWasm = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\xd6\x7b\x70\x13\xc7\x1d\x07\xf0\xdd\xbd\x93\x2d\xf9\x2c\x73\x26\x80\x63\x82\xe8\x5a\x93\x86\x47\x1b\x08\x6e\x30\x30\x19\x0a\x27\x23\x0c\x2d\x21\xe6\x31\x40\xc0\x8c\x2d\xdb\x87\x91\x2c\x4b\x70\x92\x6d\x5c\x6c\x4b\x10\x20\x4e\x30\xc4\xd0\x52\x08\x8f\x40\x13\x02\x49\x3a\x04\x1a\x1a\x4a\x78\x45\x98\x0e\xd0\x09\x05\x4a\x69\x03\x64\x00\x13\x08\x2d\xb8\x2d\xb4\x81\x4e\x1e\x6d\xd3\xd9\xef\x4a\xd8\x4d\xa7\x7f\x74\xf8\xe3\x73\xbf\x7d\xfd\xf6\xb7\x7b\x9c\x4c\x7c\x91\x1a\x27\x21\x84\x1e\x88\xc7\xe3\x44\x2b\x63\xb1\x18\x29\xa3\x31\x1a\x2b\xa3\x31\x52\xa6\xc4\x62\xe2\x51\x8d\xc5\x44\x33\x21\x65\x84\x36\x8b\xa8\x39\x46\xca\x58\xac\x59\x0c\x10\x8f\xb4\x99\xb0\xfd\x4b\xc4\x7c\xc5\x0c\xd5\x39\x7d\x15\x51\x7f\x38\x54\x1a\x31\x43\x95\xa6\x45\x98\x68\x4b\xf3\x45\x22\xa6\x15\x25\x44\x04\x99\xe5\xc1\x70\x45\x75\x69\xa8\xb6\xa6\xdc\xb4\x48\x9a\x68\xca\x92\x4d\x0b\xac\x70\x65\x6d\xc5\x83\x49\x95\xe5\xa5\x55\x66\x94\x28\xa9\x20\x62\x46\x89\x2a\x82\x6c\xcb\x5c\x58\xeb\xb7\xcc\xa9\x66\x7d\xd8\xaa\x34\x6a\xa3\xf3\x89\x4d\xb4\x3b\x22\x66\x74\xaa\x19\xa9\x0d\x26\x33\x39\x22\x51\x2b\x3f\x18\xae\x4f\xad\x68\x8f\x44\xad\x40\xd8\x1f\x4a\xad\x19\x89\x5a\x41\x33\x44\xa8\xb2\x46\x54\xaf\xa7\x33\x3b\x61\x8c\x31\x87\x2d\x83\xd9\x6c\x36\x9b\xba\x4c\x34\xd3\x05\x84\xd8\x9e\xc3\x13\xa1\xe9\xfb\x96\xc6\xe3\x24\x3b\xad\xc6\xac\x09\x5b\x0d\x8c\xf4\x2d\x9d\x3d\xb9\xa0\xbc\x62\xbe\xcf\x1f\x1a\x19\xb5\x7c\xa1\xc8\x3c\xd3\xf2\x16\x57\x34\x4c\x7b\xa2\x94\x68\x0f\x75\x75\x16\x58\xd8\xaa\xb7\x8e\x64\xe6\x76\xb5\x0e\x1b\x96\x9a\x33\xde\x34\xbd\x0d\xc4\x99\xd3\xd5\x37\xaa\xdc\x17\x34\x43\x15\xe6\x33\xf3\xbc\xc5\x15\x24\xab\x4f\xb7\x9e\x2a\x33\x3a\xad\x76\xc1\x82\x60\x83\xb7\x8e\xf4\xf8\x8f\xe5\xaa\xcc\xe8\x38\xb3\xc2\x5f\xe3\x0b\x46\xbc\x75\x44\xff\xfa\xa4\x86\x9a\xf2\x70\xd0\x5b\x47\xb2\x7b\x75\x75\x8c\xa8\x32\xa3\x93\x7d\x35\xa6\xb7\x8e\xf4\xb4\xa7\xb6\x43\x1e\x4a\x93\x1b\x26\xbd\xb4\x6e\x5b\x24\xbd\x1d\x0f\x76\x45\xfa\x38\x1e\xec\x83\xe4\x68\xdd\x32\x93\x87\x1d\x0f\x92\x91\xdc\xf4\xe4\xfa\xa4\xaf\x63\x89\x38\x44\x92\xd1\xf1\xbc\x38\xec\x0f\xc5\x0b\x43\x69\x8c\x8d\x35\x88\x41\x06\x32\xd5\x38\xcb\xaa\xdd\x6a\x01\x53\xb9\xa2\x67\x18\x7c\x82\xa1\xeb\x94\xab\xc6\x36\x16\xd0\x09\x57\x3d\x64\x84\xb2\x01\x31\x0d\x18\xdb\x92\x1d\x32\xd4\x1d\x6e\x85\xab\xc6\x06\x1a\xd0\x55\xc3\x3e\xde\xe0\x62\xde\x20\x31\x9a\xcd\x36\x12\x44\x44\xa9\xb8\xa9\x6b\x11\x9e\x9c\x63\xd8\x75\x9b\x5c\xde\xce\xa9\x6e\xe7\xaa\xa1\x23\x01\xc5\x93\xee\x70\x8b\x07\x7b\x40\x57\xfb\xca\x65\xec\x9c\x35\x8a\xb1\xa2\x97\xcb\x3e\xb1\x84\x41\xb8\xa8\x20\x50\xc0\x54\x4d\x3b\x24\x6a\x63\x2c\xc6\x9a\xbb\xaa\x4b\xd0\x6a\x37\x2b\x60\xaa\x9e\xc6\x99\xd1\x46\x03\xba\xa2\xb3\x3c\x85\x33\x99\x99\x25\xb3\xa2\x0b\x11\x32\x33\x99\xd9\x13\x4f\x2c\xfb\x6a\x51\x9e\xca\xc6\x72\xc5\x43\xa7\x38\x09\x1e\xbe\xea\xe8\xcc\x99\xe9\x24\x1e\x92\xa7\x66\x52\x4d\xc0\x15\x4f\xdb\x4f\x7a\xc4\xdc\x8a\xe7\x04\x9d\xe1\x24\x5c\xd9\x69\x28\x51\x63\x0f\x0d\x0c\x52\x48\x9e\xaa\x71\xc6\x99\x28\x40\x6d\x4c\xa5\xe4\x32\x45\xb2\x00\x66\x24\xa8\x2c\xe0\x77\xb8\x1c\xd6\xed\x72\xe2\x4a\xb5\x5b\xc1\xe5\x88\x2a\x75\x85\x2b\xc9\x9b\x51\x70\x33\x8c\x2b\x62\xe7\xa2\x04\x74\xc8\x50\x77\xb8\x45\xc7\x06\x96\xbc\x99\x0e\x71\x19\xca\x20\x31\x9c\xce\x36\xe2\x54\x44\xa9\xb8\xa9\x6b\x15\x9e\x9c\x84\xab\x51\xe4\x01\x29\xc9\x03\x42\x7a\x44\xc9\xc5\xe5\xd5\x60\x19\x3b\xa7\x8d\xa9\xb1\x5c\xf6\x25\x2b\x53\x8c\xb8\x22\x2b\x3b\x15\xff\xfa\x6b\xd7\x96\xbc\x18\x79\x15\x2d\x14\x6f\x01\x93\xaf\x00\xd3\x1d\x9c\x19\x2d\x14\x39\xf0\x60\xd8\xf5\x74\x1c\x55\x5b\xf2\xa8\x3a\xc5\x82\x4c\x89\xd1\x6e\x77\xad\xcb\xb3\xf2\xc4\xdb\x3f\x1c\x8a\x5b\x11\x0f\x23\x14\xbb\x41\xf2\x98\x18\x35\x90\x6d\xa6\x6e\xea\x75\x12\x65\x2c\x57\x3d\x19\xcd\x79\x2a\x67\x06\x0d\xb8\x19\xa7\x13\x9d\x44\x13\xaf\xe6\x08\xc5\xae\xa5\x2a\x48\x47\x05\xba\xcc\xf7\x02\xfe\x13\x19\x9b\xa9\xa1\xea\xe9\x5a\xab\x8c\xb6\x51\x63\x1b\xd5\x33\xba\x1a\xee\x3a\x8d\xbb\xce\x6e\x0d\x5c\xfc\x13\xf5\x68\x5a\x4b\xf2\x04\x38\xd1\x33\xb5\x17\x53\x01\xe5\x44\x77\xca\x3e\x0c\xd6\xb3\xba\x0f\xec\xf1\x7f\xa7\xd5\xd6\xae\x8e\xc7\x49\x06\x31\x54\x4d\x25\x45\x84\x10\x43\xd7\x32\xc5\x57\x93\x9b\x8b\x2a\x4c\xb3\x92\x10\x83\x6b\xae\xd4\xb7\x85\x57\x99\x51\x2e\x7f\x30\xb8\xa7\x90\x9b\x96\x15\xb6\x08\x31\x12\x44\xcb\x7d\x30\xc4\x1f\x8a\xd4\xce\x9b\xe7\xaf\xf0\x9b\xa1\x28\xf7\x14\x12\x62\x74\x10\x8d\x77\xfb\x3a\xfd\x8f\x45\xe2\x54\xeb\xd7\x7d\xd4\x7f\xaf\xb3\x99\x6a\xaa\x9d\x10\x14\xa4\x29\x68\xda\x43\xb5\x44\x26\x69\xef\xd4\x44\x73\x7c\xe7\x27\x90\xb7\x5d\x85\x57\x16\x5f\x84\xef\x05\x7e\x0b\x77\x97\xfc\x1a\xbe\x3f\xe5\x04\xfc\xc3\xc4\xa3\xf0\x31\xef\x41\xd8\x50\xf8\x2e\xec\x18\xb7\x1b\xce\x9a\xf8\x26\xbc\x56\xfc\x1a\x5c\x34\x67\x0b\xec\xe7\x5f\x0f\x4f\x35\xac\x81\xcb\x57\xad\x84\x45\xdb\x97\x43\xfd\x48\x4c\xe6\xbf\xd4\x20\xf3\x7e\x11\x81\xcf\xb9\x42\x70\x66\xe1\x7c\x38\x38\x58\x0e\xbf\x5c\x5f\x02\x4f\x1e\x9f\x01\x5f\xfa\xac\x18\x3e\xf3\xf8\xf7\x61\xba\xe9\x85\xef\x6e\x1d\x23\xe7\x5f\x1e\x05\xff\xe4\x7a\x12\x5a\xbe\x21\xf0\x8f\x3b\x07\xc2\x49\x7f\x75\xc3\x2d\xa3\x5d\xf0\x52\x6b\x0e\x24\x57\xb3\x61\x8f\x27\x33\xa1\xad\x35\x4d\xd6\x7b\x93\xc0\xad\xe3\xbe\xcc\x10\x8e\x7d\xfd\x3e\x6c\xb7\xdf\x85\x8f\x86\x6f\xc3\xca\x0b\x37\x60\x7c\xdc\x55\xb8\x68\xef\x45\x38\xf1\x9b\xe7\xe1\xdf\x5e\x3e\x0d\xad\xec\x5f\xc1\x13\x2b\x8f\xc1\xbb\xf6\x23\xb0\xb3\x65\x3f\xfc\x85\xb6\x17\x4e\x5b\xb3\x0b\x1e\xce\x7d\x03\x7e\xba\xfd\x55\x78\x3f\x7f\x0b\x4c\x1c\x5b\x0f\xa7\xcf\x5c\x0b\xdf\xfe\x73\x2b\x3c\xbd\xe2\x79\xf8\x8e\x7b\x29\x7c\xf6\x68\x23\x3c\x5a\x5e\x0f\x6f\xa5\x59\xf0\xcc\xae\xa0\xdc\xcf\xcc\x2a\x78\x56\x29\x97\xfd\xef\x94\xc0\x03\x95\x33\x61\x61\xee\x54\xb8\xf4\xec\x24\x58\xfb\x62\x11\xec\x57\xe4\x91\x71\xda\x68\xd8\x7c\x72\x04\x2c\x58\x95\x0f\x37\x4f\x7f\x5c\xee\xab\xff\x40\x58\xdd\xe9\x86\xe7\x0e\xb8\xe0\xb5\xd6\x87\xe1\x7a\xdf\x43\xf2\x7c\x46\x66\x49\x7b\x3a\xe0\x8f\xee\x28\xf0\xc2\xe9\x7f\x39\x84\x87\xf6\x7c\x0e\xc7\xac\xbb\x07\x4b\x1b\xef\xc0\x47\xe6\xdd\x86\xd6\xd4\x4f\x60\xb9\xd1\x01\x3b\xf3\x3f\x82\x19\x03\x7e\x0f\xdf\x77\xfd\x06\xda\x73\x4e\xc1\x1b\x3d\x4f\xc0\x29\xd9\xed\x70\x72\xcf\xc3\xf0\x52\xef\xfd\xf0\xde\x23\x7b\xe1\xc6\x47\xdf\x86\x1f\x0c\x7d\x13\x2e\x1f\xbd\x1d\x9e\x7c\x7a\x2b\x5c\x5b\xb6\x11\x7e\x5c\xbb\x0e\xee\x5b\xd5\x06\x9d\x6f\xac\x84\x37\x7f\xb9\x02\x8e\xbd\xb1\x04\xe6\xa6\x35\xc1\x05\x83\x17\xc1\x71\x53\x22\xf0\x87\x3f\x08\xc1\x92\x1d\x7e\xf8\xca\xf9\x4a\x38\xcb\x56\x06\x5b\x46\xce\x81\x43\x42\x33\xe0\x84\x1d\x53\xe0\xf9\x8e\x49\xf0\x9c\x6b\x82\xcc\x33\xb7\x10\xba\xb6\x7e\x17\x56\x5d\x1f\x05\xfb\x0f\x1e\x2e\xfb\x23\x4f\xc0\x53\x47\xbe\x05\x13\x99\x03\x60\xbf\xb9\x6e\x78\x65\xb7\x0b\xda\x6c\xb9\xb0\xad\xa4\x17\x6c\xda\xd7\x43\xce\xeb\xa9\xc1\xc5\x0b\xd3\xe0\x0b\xe7\xa8\x3c\xaf\xe1\xff\xc4\xe7\xe7\xd0\xa6\xcf\xe0\x55\xe5\x1e\x2c\x09\xdd\x81\xf9\x1f\xdd\x86\xa5\xdf\xbb\x09\xaf\x1f\xbc\x26\xc7\x0f\xbd\x0c\x6f\xbd\x7e\x01\xfa\xfb\x9f\x87\x4f\xbd\x7c\x06\xce\xed\xf3\x01\x3c\xbd\xee\x38\x5c\x9d\xd3\x0e\x37\x6d\x3a\x0c\x3f\xcd\x7b\x0f\xfe\xf8\xad\x9f\xc3\xa5\xc3\xf7\xc0\x7d\xed\x3f\x95\x79\xa7\xee\x94\xe3\xae\xbf\x0a\x3f\xaf\x7f\x05\x16\x66\x6d\x84\x67\x76\xac\x83\x9b\x8d\x35\xf0\xad\x2b\xad\xf0\xef\x4d\x2d\x30\xfe\x8d\x65\xb0\xe8\x58\x0c\x4e\xf4\x2f\x86\x2b\xb2\xeb\xe1\x3f\x0e\x5a\x72\xde\xfc\x10\x7c\x29\x27\x20\xe3\x93\x26\xfc\x62\xb1\x0f\x36\x0d\x9b\x0b\xbf\x73\x6b\x16\xec\xfb\xda\x74\x38\xac\xb4\x18\xd6\xbb\x26\xc1\xce\xcb\x45\xb0\x75\x5b\x21\x7c\xd6\x1c\x03\x8b\xbf\xfd\x14\x5c\x78\xbf\x00\x1e\x48\xe4\xc3\xc1\xab\x87\xc0\xc4\xdc\xc1\xb0\x71\xe8\x63\x70\x8e\xe2\x96\xe7\x77\xd1\x05\x97\xfc\x2c\x17\x1e\x5f\xd9\x5b\xce\x9b\x9f\x0d\x77\x8d\x77\xc2\xa7\x07\x38\x60\x46\x86\x0d\x7e\xfc\x17\x02\xc5\xef\xa1\x96\xe5\xe1\x85\xe2\xef\x6c\x3e\x3d\x5c\x6d\x86\xc8\xbf\x03\x00\x00\xff\xff\xd9\x0e\xec\x99\x71\x0d\x00\x00")

func BchainBchainWasmBytes() ([]byte, error) {
	return bindataRead(
		_BchainBchainWasm,
		"../../bchain/bchain.wasm",
	)
}

func BchainBchainWasm() (*asset, error) {
	bytes, err := BchainBchainWasmBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../bchain/bchain.wasm", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _SystemSystemWasm = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x5c\x90\x31\x4f\xeb\x30\x14\x85\xcf\xbd\x71\xfa\x5e\x5f\xd4\x57\x0f\x6f\x78\x0b\x92\x41\x0c\x1d\x41\x42\x9e\x6b\xa1\xae\x08\x31\xb2\x34\x51\xf0\xd0\xd0\xa4\x28\x31\x95\x3a\xb9\x15\x0c\x8c\x5d\x58\x91\xf8\x0f\x6c\xfc\x11\x46\x7e\x0a\x72\x52\x84\x60\xc9\x39\xf7\xcb\xd1\x3d\xb6\x91\x35\xe5\x00\x00\x3d\xae\xd7\x6b\x88\x94\x3c\xd2\xc8\x7b\x8f\x94\x3c\xf9\x94\xbd\x07\xbf\x85\x5f\x71\x64\xab\xe5\x20\xcb\xdd\x6c\x51\x4d\x1b\x5b\x5d\xd9\x1a\x08\x6c\x98\x2f\x2a\x57\x67\xb9\x9b\xe6\xb5\xcd\x9c\x05\x05\xda\xb7\xe5\xcc\x4d\x96\xb6\x72\x88\xda\xb9\xb1\xee\xc2\x36\xb7\xf3\xdd\xdc\x6b\x5c\x3d\xb7\x15\x38\xba\x0b\xdb\x39\x82\xb8\x0f\x86\x6e\x80\xb8\x45\x04\xfa\xf5\x1a\x4c\xd4\x2b\x6d\xb9\xa8\x57\x8c\xbd\xe9\xe5\x99\x6e\x56\x8d\xb3\xe5\xf1\x49\xd7\x76\xba\x2b\x9f\x9c\xe7\x88\xff\x7e\x67\xe8\xf5\x37\x61\x01\xfe\xbc\x6c\x42\xc7\x73\xf8\x12\x79\x1e\x1b\x18\x8c\x58\x98\x2d\x5d\x1f\xb0\x66\xa1\xd8\x1c\x15\x12\xad\x28\x52\x6c\xde\x51\xc8\x4f\x15\xfb\xc1\x49\xcd\x2a\x88\x14\x9a\x0f\x15\x1b\x68\xfe\xaf\x58\x91\xe6\x7f\x8a\xbb\xa4\x66\xd9\x5a\x59\x68\x1e\x75\x91\xdf\x5d\x44\x7c\x45\x42\x89\x2a\x14\xcb\x1d\x51\x24\x23\x13\xe0\x96\x0a\xcd\x22\x49\x1e\xda\x23\x2b\x28\xc8\x38\x49\x9e\xda\xd7\x81\x11\x89\x50\x63\x00\x46\x26\xc3\x1f\xb7\xfc\x08\x00\x00\xff\xff\x64\xcf\xc7\xf7\xc3\x01\x00\x00")

func SystemSystemWasmBytes() ([]byte, error) {
	return bindataRead(
		_SystemSystemWasm,
		"../../system/system.wasm",
	)
}

func SystemSystemWasm() (*asset, error) {
	bytes, err := SystemSystemWasmBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "../../system/system.wasm", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

// Asset loads and returns the asset for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func Asset(name string) ([]byte, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("Asset %s can't read by error: %v", name, err)
		}
		return a.bytes, nil
	}
	return nil, fmt.Errorf("Asset %s not found", name)
}

// MustAsset is like Asset but panics when Asset would return an error.
// It simplifies safe initialization of global variables.
func MustAsset(name string) []byte {
	a, err := Asset(name)
	if err != nil {
		panic("asset: Asset(" + name + "): " + err.Error())
	}

	return a
}

// AssetInfo loads and returns the asset info for the given name.
// It returns an error if the asset could not be found or
// could not be loaded.
func AssetInfo(name string) (os.FileInfo, error) {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	if f, ok := _bindata[cannonicalName]; ok {
		a, err := f()
		if err != nil {
			return nil, fmt.Errorf("AssetInfo %s can't read by error: %v", name, err)
		}
		return a.info, nil
	}
	return nil, fmt.Errorf("AssetInfo %s not found", name)
}

// AssetNames returns the names of the assets.
func AssetNames() []string {
	names := make([]string, 0, len(_bindata))
	for name := range _bindata {
		names = append(names, name)
	}
	return names
}

// _bindata is a table, holding each asset generator, mapped to its name.
var _bindata = map[string]func() (*asset, error){
	"../pledge.wasm": PledgeWasm,
	"../../bchain/bchain.wasm": BchainBchainWasm,
	"../../system/system.wasm": SystemSystemWasm,
}

// AssetDir returns the file names below a certain
// directory embedded in the file by go-bindata.
// For example if you run go-bindata on data/... and data contains the
// following hierarchy:
//     data/
//       foo.txt
//       img/
//         a.png
//         b.png
// then AssetDir("data") would return []string{"foo.txt", "img"}
// AssetDir("data/img") would return []string{"a.png", "b.png"}
// AssetDir("foo.txt") and AssetDir("notexist") would return an error
// AssetDir("") will return []string{"data"}.
func AssetDir(name string) ([]string, error) {
	node := _bintree
	if len(name) != 0 {
		cannonicalName := strings.Replace(name, "\\", "/", -1)
		pathList := strings.Split(cannonicalName, "/")
		for _, p := range pathList {
			node = node.Children[p]
			if node == nil {
				return nil, fmt.Errorf("Asset %s not found", name)
			}
		}
	}
	if node.Func != nil {
		return nil, fmt.Errorf("Asset %s not found", name)
	}
	rv := make([]string, 0, len(node.Children))
	for childName := range node.Children {
		rv = append(rv, childName)
	}
	return rv, nil
}

type bintree struct {
	Func     func() (*asset, error)
	Children map[string]*bintree
}
var _bintree = &bintree{nil, map[string]*bintree{
	"..": &bintree{nil, map[string]*bintree{
		"..": &bintree{nil, map[string]*bintree{
			"bchain": &bintree{nil, map[string]*bintree{
				"bchain.wasm": &bintree{BchainBchainWasm, map[string]*bintree{}},
			}},
			"system": &bintree{nil, map[string]*bintree{
				"system.wasm": &bintree{SystemSystemWasm, map[string]*bintree{}},
			}},
		}},
		"pledge.wasm": &bintree{PledgeWasm, map[string]*bintree{}},
	}},
}}

// RestoreAsset restores an asset under the given directory
func RestoreAsset(dir, name string) error {
	data, err := Asset(name)
	if err != nil {
		return err
	}
	info, err := AssetInfo(name)
	if err != nil {
		return err
	}
	err = os.MkdirAll(_filePath(dir, filepath.Dir(name)), os.FileMode(0755))
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(_filePath(dir, name), data, info.Mode())
	if err != nil {
		return err
	}
	err = os.Chtimes(_filePath(dir, name), info.ModTime(), info.ModTime())
	if err != nil {
		return err
	}
	return nil
}

// RestoreAssets restores an asset under the given directory recursively
func RestoreAssets(dir, name string) error {
	children, err := AssetDir(name)
	// File
	if err != nil {
		return RestoreAsset(dir, name)
	}
	// Dir
	for _, child := range children {
		err = RestoreAssets(dir, filepath.Join(name, child))
		if err != nil {
			return err
		}
	}
	return nil
}

func _filePath(dir, name string) string {
	cannonicalName := strings.Replace(name, "\\", "/", -1)
	return filepath.Join(append([]string{dir}, strings.Split(cannonicalName, "/")...)...)
}
