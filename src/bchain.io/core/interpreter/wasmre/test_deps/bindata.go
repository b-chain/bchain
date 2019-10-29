// Code generated by go-bindata.
// sources:
// test.wasm
// DO NOT EDIT!

package test_deps

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

var _testWasm = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xac\x56\x4d\x6f\x1c\x45\x13\xae\xee\xf9\xde\xde\xb5\x3b\xc9\x9b\xbc\x79\xf7\xe5\xa3\xc7\x27\x8b\x10\x69\x77\x23\xaf\xed\x90\xc3\x56\x22\x47\x39\x05\x4b\x41\x1c\x72\xb1\x67\x77\xdb\xf1\x7e\x5b\x33\x83\xb1\x25\xc4\x6c\x20\x22\x10\x24\x30\xca\x1f\xc8\x15\x71\xc9\x81\x23\x87\xdc\xb8\x5a\xdc\xb8\xe5\x82\xc4\xd1\x88\x3f\x80\xaa\x67\xd6\x76\x0c\x07\x84\x90\x0f\xdb\xf5\x74\x55\x3d\x4f\x57\x55\x4f\x1b\xa2\x64\x54\x01\x00\xf6\xfd\x74\x3a\x85\xd2\xa6\x9d\x65\x59\x06\x9b\x2c\x83\x4d\x2b\xcb\x32\x96\x6d\xf2\xdc\x66\x99\x01\x60\x13\xd8\xc7\x84\xb1\x6c\x13\xc8\xc8\xf8\xcf\x0f\xa7\x53\x98\xb7\xf4\x78\xb7\x12\x75\xd2\xde\x64\xbc\x91\xe8\x71\x57\xc7\xc0\x08\x73\xa3\x24\xd1\x71\x0a\x16\x19\xe5\xf6\x70\xd2\x19\x6c\x8c\x3f\x18\xb5\x75\x0c\x2e\x41\xa5\x6e\x7b\x43\x8f\x76\x86\x51\x47\x03\x98\x80\x6e\x7b\xe3\x81\x4e\x81\xcf\x8c\x44\xa7\xf9\x4e\xa0\x47\xbd\x74\x6d\x57\x8f\x67\xd9\x7a\xc9\x1d\xbd\x87\xdd\x6e\xac\x93\x04\x6c\x82\xac\xe1\xe4\x41\xbe\x72\x47\x7a\xd4\x19\xed\x14\x69\xc8\xd8\xd9\x3f\x31\x92\x19\x81\x9d\x6c\x47\x75\x70\x0c\x9e\x6c\x47\x8d\xa5\xe6\x89\xb1\x54\x6f\x80\x63\x3d\xa5\xc2\x9c\xb7\xc1\xf1\x7c\xcb\xb6\x03\xdf\xf7\xbd\xc0\x77\x03\xdf\x7e\x44\x1b\x6c\x07\xc0\xf9\xd4\xac\x80\x79\xdf\x7c\x32\x9d\xc2\x05\x22\x98\xc4\xfb\x1c\xc4\xc6\xfd\x66\x92\xc6\x43\x3d\x5e\xef\xc0\xfc\xe5\x8d\xfb\x77\x97\xd2\xc9\x40\x8f\x9b\x9d\x58\x47\xa9\x5e\x5b\xef\xdc\xab\x6d\xf4\x7a\x20\x2f\x1e\x6f\x2d\xa7\x71\x34\x4e\x74\xbc\xb6\xde\xe9\xc1\xb9\x4b\xc7\xf8\x6a\x3b\x1a\xea\x71\x47\xbf\xbb\xb5\xb6\xde\x81\xf3\x6e\x9e\x01\x2e\x78\x85\x3f\xfc\x27\x38\xf6\x80\x8b\x4e\xaa\x93\xb4\x0e\x97\x04\xfd\xa2\x69\x40\x1d\xfe\x7b\xca\x6a\xc0\xe5\x0a\x59\xb7\xe2\xfd\x49\x3a\xc1\x9d\x1e\xfc\x2f\x20\xfb\x5e\x1a\x75\x06\x50\xf5\x76\xa3\x78\x3d\x8a\x23\xf8\xbf\x89\x79\xbf\xb0\x5e\x33\x3e\x79\x07\x5e\x9f\xa7\xf5\x4d\xea\xe7\xdd\xbc\x9d\x6f\x78\x23\x3d\x7a\x4f\x27\x29\xbc\x29\x86\x93\xc9\xce\xed\x49\xac\x77\x75\x0c\x2a\x78\x48\xf5\x81\xd2\xe3\x2f\xa9\x92\x3f\x98\x62\x59\x19\xcf\x78\x4b\xc1\x55\x80\xb5\x0a\x28\x40\xd6\x0f\x19\x42\x08\x56\x4b\x31\x05\xfd\x90\x1b\x6c\xc1\x0a\x41\xf1\xab\x00\x15\x10\xca\x9a\x17\x08\x42\x3c\xa5\x22\x33\x37\xe3\x2d\x04\x84\x45\x6e\xe3\x0b\x18\x2c\x04\x4d\x6e\x23\x84\x3e\x42\xe8\x9d\x4d\x4c\x1b\xae\xd5\x52\xb6\x72\xfb\xa1\xa3\x5c\x4a\xec\x85\xae\x72\xf2\xc4\x02\x25\x3a\x0a\x94\x27\x2d\xde\x52\xac\x88\x64\x14\xe9\x15\x91\xde\xa9\x48\xff\x74\xa4\x42\x5b\x31\xe5\x4b\x4b\x05\xca\x6a\xf2\x1b\xbc\xa5\x38\xb2\x3b\x15\x3a\x87\x85\xa5\x61\x68\x29\x8e\x59\x7f\x81\xd3\x01\x72\x17\xa1\x02\x09\xbc\xc5\x5b\x2a\x28\xa8\x02\x64\x71\xe8\x14\x54\x0e\x51\x59\x39\x15\x0f\x5d\x65\x19\xaa\x32\x17\x02\x21\xe4\x94\x84\xab\x00\x6f\xf4\xd1\x96\x16\x52\xec\x0b\xe8\x37\xb9\x2d\xc4\x2f\x7f\x59\x17\xbf\xc9\x6d\xe5\xa3\xec\x4b\xa0\xd2\x28\x1f\xa1\xc9\xcb\xf9\x0f\x15\xcb\xe5\x2d\xe5\x5f\x05\x49\x3a\xc8\x8d\xa4\x58\x08\xa1\x43\x07\x50\x4e\x3f\xb4\x95\x43\x52\xdc\xd0\x51\x76\x71\x6a\xe3\xa8\x5c\xe5\x63\xb9\x2f\x6d\xb4\x6f\x63\x4d\xb2\xb3\xa5\x73\x8b\x24\xee\xa9\x24\xde\xe9\x24\x4c\x91\x1a\xbf\x2f\xed\xaa\xf2\x95\xbf\xc8\xcb\x8a\x0f\x8c\x36\x32\x7c\xc5\xfb\xb9\x42\xef\x5f\x50\x89\xb6\x74\xfe\xb1\x40\x0a\x46\xe2\x9d\x55\xfa\xc8\x8c\xb0\x93\xf1\x6c\x56\x69\x39\x58\x70\xf3\x01\xa4\x21\xa1\x0a\x9f\x65\xe3\x08\xa1\x6d\xb5\x14\x57\x36\xb5\xd7\x26\x36\x27\xb4\x8b\xf6\x1a\x36\x0a\x2d\x9b\x72\xb8\x8b\xbc\x4c\x53\xab\x5c\x94\x44\xa9\x6c\x21\xbe\x36\x97\x08\x5f\x00\x1e\x02\x72\xfc\xd5\x93\x52\x3c\x31\x98\xa2\x3f\x26\xcf\x89\xcf\x8f\x4d\x79\x5e\x7c\x67\x44\x32\x12\xf9\x12\xa4\x5f\x45\x2b\x64\xbc\x85\x47\x20\xbd\xb5\x0a\xe0\x33\x46\x18\x0b\x99\x30\xf7\x4d\x88\xaf\x4c\x30\xcf\x90\xe1\x11\x93\x8c\x6e\xda\x17\x39\x25\x18\x40\x7c\xf6\xb0\x48\x78\x3c\x5f\x07\x7c\xb0\x00\x74\xec\x29\xc7\x40\x01\xbe\x64\x7d\x59\xc6\x03\x2e\xfd\xaa\x02\x2c\xbd\xc3\x9e\xf0\x02\xf5\xab\x33\x9f\x03\xd6\x97\x15\x7c\x76\xe2\x73\xc8\x0a\xf4\xd8\x47\xce\xe1\xf3\x93\xfd\x29\x53\x46\x3e\xdd\xe5\x03\x9e\x37\xe0\xb7\xe9\x59\x29\x53\xf6\xaa\x14\xc0\xb7\xa5\x58\x80\x57\xe5\xbc\x95\x67\x2a\x5c\x6a\x67\x84\xec\xe4\x98\x5f\x9d\x85\xcc\x58\xa7\x2c\x67\x7d\x94\x97\x83\x09\xf3\xbd\x07\x10\x3f\xfe\x49\xc6\x21\xe4\x32\x48\x41\x93\xdf\xa1\xcf\xcf\x75\x58\x54\x80\xfc\x7a\xde\x18\x5c\xec\x37\xb9\xa2\x6d\xe9\x9a\xe4\x87\xc5\x4c\x3d\x36\x39\x25\xbf\xc9\x3e\x12\x3f\x99\xbc\xfc\xd4\x7c\xd5\x06\x0b\x8c\xd2\x32\xfc\x1d\xf0\x8a\x14\x55\xc5\xf2\x03\x51\x17\x15\x33\x47\x92\xa5\xaa\x31\x71\xca\x15\xc3\x40\x06\x21\x20\x0d\x5f\xcd\x4c\x10\xcc\x28\xac\x56\x19\x84\x10\xdf\xd2\xd7\x62\x0e\xd0\x16\x76\x0b\x01\x00\xa5\x70\x93\xfd\xd1\x64\x08\x80\x4a\x38\xe3\x68\xa4\x01\xb0\x26\x3c\x1d\xc7\x93\xb8\x64\x26\x4f\xd8\x7b\x7b\x6d\x20\xcd\xb4\xda\x03\xc0\x97\x20\xe6\xb6\xf5\x70\x38\x51\x1f\x4e\xe2\x61\x37\x24\xbf\x23\x10\x57\x6a\x7b\x0d\xdd\x5c\x69\xd7\x96\x56\xae\xd5\x1a\xf5\xee\xf2\x4a\xa7\xde\x68\x6c\x2d\xd7\x57\xb7\x3a\x2b\x8d\xda\xb5\xe6\x52\x63\x35\x5a\xad\x2d\x2d\xd7\xbb\x40\xa3\x28\x6e\xfd\xfd\x00\xd5\x4b\x54\xa4\x76\xa3\x61\xaf\xab\xb6\xf5\x9e\x8a\xf2\xf7\xdf\x30\x33\x51\xa6\x57\x49\x15\xff\x73\x00\x4e\xb9\x28\x45\xed\x4e\x57\x6f\x3d\xd8\xee\x01\x0d\x90\x28\xd3\x8b\xaf\x62\x9d\x5e\x57\x14\xf3\x8c\x8b\xb9\xfc\xe1\x3f\xc1\x9e\xe7\xd8\x52\xbd\x31\xc3\xfe\x08\x00\x00\xff\xff\xa8\x55\x6e\xa7\x27\x09\x00\x00")

func testWasmBytes() ([]byte, error) {
	return bindataRead(
		_testWasm,
		"test.wasm",
	)
}

func testWasm() (*asset, error) {
	bytes, err := testWasmBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "test.wasm", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"test.wasm": testWasm,
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
	"test.wasm": {testWasm, map[string]*bintree{}},
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