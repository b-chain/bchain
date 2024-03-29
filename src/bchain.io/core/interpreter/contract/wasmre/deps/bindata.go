// Code generated by go-bindata.
// sources:
// pledge.json
// system.json
// bchain.json
// DO NOT EDIT!

package deps

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

var _pledgeJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x94\x56\x5d\x93\xa2\x48\x16\xfd\x2f\xf5\x3a\x1b\xdb\x16\x6a\x6c\xbb\x11\xf3\x70\x6e\x9a\x60\xaa\x50\x7d\xd1\x02\xe1\x65\x43\xa4\x4c\x15\x14\xbb\x2d\x8b\xc2\x8d\xfd\xef\x1b\x99\x48\x57\xcf\x6c\x77\xc4\xce\x83\x41\x72\xbd\x9f\xe7\x9c\xcc\xe4\xdf\x0f\xfb\xd3\xeb\xcb\xb7\x7f\x9d\xd6\xc7\x97\x87\x7f\x3e\xd4\xeb\xcb\xf1\xdb\xcb\xdf\xe3\xf5\xe5\x18\xca\x87\xbf\x3d\x6c\xaa\xdc\x98\xe1\xb9\xb7\x8c\x7b\x00\x40\x97\x44\x40\x03\xd3\x04\xee\x27\x78\xa0\x6d\xed\x7e\x4a\x10\x7c\xda\xf6\x3f\xe3\xfe\xa4\xad\x03\xd1\xbe\xbb\xbf\x61\x35\xd4\x34\x19\xdd\xff\x9f\x7e\xda\x6a\xd2\x58\x0d\x80\xd7\x98\xb4\x02\xc6\x08\xca\xec\x94\x2c\x93\x38\xe8\xad\xe3\xd1\x75\xe5\x04\xbb\xcc\x7b\x5f\xad\x57\x61\xf5\xec\xb9\x4d\xc2\x10\x70\xa2\x6b\xae\x1f\x77\x49\x3f\x3c\x67\xce\x70\xbb\x71\xa2\x6b\xea\x45\x0d\x80\x71\x1a\x0f\x1d\x3a\xba\xb7\x8d\x13\x35\x39\x20\xad\xef\x8e\x0e\x99\x33\xec\x6d\x8e\xee\x21\x77\x47\xbb\xd4\x0b\x9b\x74\x15\xdc\x3a\x7f\x37\x0e\xde\xb2\x53\xd8\x98\x9a\xbf\xae\x97\x16\x49\x39\x3a\xa5\x2b\x06\x6a\xdb\xa3\x97\x7a\x53\x53\xbb\x07\x24\x36\x8f\x88\xf3\x32\x77\xf3\x72\x1d\xe7\x55\x0e\xb8\x6d\x5c\x79\x4b\x57\xe1\x22\x5d\x05\x8f\xd9\x84\x41\xed\x7c\xd3\x4d\x3f\x6c\xfc\xe3\xfb\x5b\xee\x7c\xf4\x2d\x26\x41\x6f\x73\x3c\xbf\xad\xe3\x61\x0f\xf0\xdb\x59\x4e\xc6\xf6\x5e\x66\x1a\x84\xc1\x11\x06\x23\x81\xcd\x84\xc0\x4a\x18\x02\xd4\xc2\xda\xb0\x32\x74\xc4\x63\x6d\xad\x12\xe0\x97\x6f\x3a\x31\xeb\xc2\xcb\xe2\xe8\x35\xeb\x4f\x87\xd0\x48\x56\x8f\xe7\x27\x7f\x49\x75\xe6\x45\x45\xea\x44\x9c\x39\xa3\x4b\x70\xa2\x4b\x1a\x87\xa7\xf4\x39\x1a\x02\xf5\x2f\x7c\xa6\x65\xea\x45\x65\x66\x7d\x7a\xd9\xcf\x7c\x9e\x26\x5d\x9e\xd1\x31\x8c\xe8\x00\x0c\xb6\x3f\xf3\xf3\x97\x6e\xf7\xfe\x25\x2d\xa2\x41\x2e\x23\x4e\x6a\x7c\xa1\xae\x0f\x06\x53\x57\x8f\x11\x8a\x1f\xf2\x82\xd4\x7c\xe3\xbd\x97\xa9\x97\x97\x4b\x27\x75\x5f\x26\x0c\x59\x1f\xa9\x9d\x1b\xb3\xaa\xda\x1b\x3d\xaa\x7f\x0c\x66\xe6\x29\x70\xfa\x4c\x5b\x3d\x05\x33\x11\x84\x16\x24\x3d\xc6\xf1\xb2\xc7\x2d\x11\x24\x30\x56\x90\x7d\x0c\x34\x29\x04\xd4\x03\x7d\x95\x90\x1a\x3d\xa9\xd7\x3b\xf3\x5f\x20\x30\xde\x8c\x21\x30\xe6\x14\xb4\x2e\x24\x77\x7e\xe6\x59\xc0\xfd\x2a\x51\xed\xb1\xc7\x58\xe2\x79\xad\xe0\x9f\x51\x83\xa9\x91\x52\xc1\xd7\x68\x0a\x13\x4b\xdb\x1f\x73\xec\x21\x4c\xac\x20\x78\x5d\x7e\x29\x3d\x43\x13\xb0\xc3\x5c\x69\x65\xea\xf3\x7a\x07\xd7\xdb\x9b\x1e\xfc\x33\xde\xa1\xb1\xaa\xfb\x78\x27\xd2\x50\xda\xd4\xe3\x90\xbe\x32\x6b\xa6\x3d\xc6\x33\xf6\x21\xb1\xd9\x61\x81\x31\x6b\x58\xbf\x02\x53\x2a\xa0\x58\x34\xca\xd6\x93\x5e\xc5\x14\x56\x1a\x0d\x48\x21\x1c\xb5\x79\xe9\x73\x50\xfb\xcc\xa9\xe9\xe9\x5e\xb7\x90\x4a\x22\x31\x6b\x19\xe8\x19\xdb\x9c\x8b\x8b\x83\xb3\xc1\xc8\x6f\xb1\x39\x24\xa2\x6c\x71\x52\x81\x56\x49\x6b\xf7\x14\xe2\xca\xc1\x8e\x8d\x9d\x02\xad\x2c\x56\x72\xdc\xe2\xcc\xcc\x0e\x34\x91\x24\xcc\x15\xdb\x5a\x08\xb4\xf2\x6d\x6c\x21\xc9\xe4\x14\x96\x0f\xdf\xc1\xae\x36\x7d\xda\x19\x03\x3d\xdb\xe0\x79\x02\x48\xa9\x3f\x30\x37\x5c\x5a\x3c\x85\x89\x23\x5c\xe6\xfb\xea\xe9\x17\x7c\x4f\x34\xe2\x3f\xf1\xfd\xa8\x3b\x8e\xe9\x0f\x1c\xbf\x00\x1d\xc7\xca\xab\x34\x7a\x13\x0d\xaf\x62\xb1\x37\x18\xfa\x77\xfc\x5a\x4e\xe3\x62\x52\x83\x0c\x6e\xbf\xe0\x7c\xf9\x43\xce\xbf\xc8\xf9\x8c\x7d\x56\x70\x63\x5e\x83\xa4\xe1\xa5\xf9\x6e\x1b\x59\xbe\x92\x9f\xf0\x65\xea\x12\x26\x8a\x0d\x2f\x0a\x1d\xaf\xe6\x27\x71\xd9\x77\x9a\xf9\xae\xd3\x77\x68\x32\x7d\x2e\x14\x45\xc5\x84\x11\xda\x59\x6c\x1e\x79\xef\x5f\x9a\xbd\x72\xef\x5d\x7a\xe7\xae\x7f\xa2\x71\x22\x12\x71\xd7\xc5\x51\xb7\x5c\x2d\x0c\x57\xd3\xc2\xce\x39\x4e\x84\x67\xb1\x0e\xc8\x81\xb1\xab\x67\xb3\xc6\x32\x11\xf2\xbb\x16\x60\xfa\xaf\x03\xad\x24\x2f\xc1\xa2\x91\xb6\x0e\xc6\x89\x18\xdb\xd8\x29\x59\x3d\xe9\x36\x7f\xa0\xd5\xe6\xae\x3f\x36\x9c\x67\xf2\x43\x13\x03\xbb\x6f\x43\xba\xe3\x4d\x9d\x26\x9a\x81\xbd\x9b\x08\xab\xcf\x82\xa5\x04\x33\x2a\xe8\xd0\x68\xe8\x9b\xd2\xca\x81\x66\x0d\xd3\x71\xfd\xa4\xb0\x30\x58\x17\x46\x8b\xac\x8c\xb6\x77\x98\x99\x1e\xf6\x56\x6f\x06\xb3\x76\x6d\xf6\x9c\x90\x12\xca\xee\x2f\xd7\xcc\x26\x45\x7d\x0d\xdb\xb3\x56\xd2\xb6\xfe\xd3\x59\x63\xb4\x77\xb0\xda\x13\x76\x6f\xda\xb3\x46\xb2\x58\x40\x70\x08\x8d\x85\xe1\xa7\xd2\x28\xec\xfe\x69\x35\x66\xd7\x16\x6f\xa5\x30\xb5\x7a\xd9\xc3\x25\x8d\xa9\xd9\xff\x2c\x1a\x08\x4e\x34\xdd\xf5\xd2\xae\xad\x0e\xf4\xc7\x3c\xcf\x1e\xa3\xf0\xd8\xea\xeb\xff\x9f\x63\x76\xbf\x3f\x68\x6b\xcf\x09\xf0\x18\x3f\xb1\xf1\x55\xdd\xcf\x5b\x0d\x61\x6e\x9b\xc1\x7c\xaf\x2c\xd6\x50\xc6\x46\xf8\x22\xea\xab\xd0\x95\xf1\xa9\xc0\xcc\x73\x12\x44\x26\x80\x43\xcc\xa9\xff\xfd\x9c\x07\x2f\x30\x17\x2b\x73\xf7\x0e\x6f\xe9\xd1\xdc\x85\x44\x3e\x2e\x93\xcd\x31\x2a\xd2\x38\x7a\x35\xef\x35\x30\x9f\xdd\x68\x90\x1c\x36\x83\xf4\x20\x1b\x76\xd4\xed\x69\xa9\x7a\xc1\x01\x7d\x5e\xee\x04\xdf\xd2\x5d\x70\x48\xfa\x4f\xcb\xb2\xf4\x97\xcf\xc3\x70\xe9\xf7\x82\x65\x51\x3f\x2d\xf3\xf2\xc9\xd4\xd8\x02\xa2\x09\x9a\x34\x8e\xca\x6c\x59\xe9\xee\xfe\x50\x9e\xfb\x9a\xf5\xa3\x6b\x2e\xe8\xbc\x69\xe8\x9a\xf5\x59\xa7\xf1\xf0\x2d\x8f\xf3\xca\xd4\xad\x20\xe7\x6a\x15\x3c\x26\x65\xf8\x96\x7b\xee\xe5\xf9\x7f\xe2\xf0\xdb\x97\x05\xf5\x3a\xbb\x89\xe9\x41\xce\xa7\x3f\xc4\x44\x4e\x74\x4e\x9d\x5d\x4f\x4d\xee\xdf\x05\xf7\x98\xac\x1f\xee\xb2\x8f\x6f\x05\x98\x3d\x2d\xea\xc7\x6d\xd7\x5b\xe4\x8d\x7a\x49\x5c\x03\x46\x43\xfa\x12\xac\xfa\x7f\xcc\x07\xfc\xfe\xfb\xc3\x7f\xfe\x1b\x00\x00\xff\xff\x06\x7c\x7e\x2e\xb0\x09\x00\x00")

func pledgeJsonBytes() ([]byte, error) {
	return bindataRead(
		_pledgeJson,
		"pledge.json",
	)
}

func pledgeJson() (*asset, error) {
	bytes, err := pledgeJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "pledge.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _systemJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x44\x51\xd1\x72\x9b\x30\x10\xfc\x17\x9e\x3b\x0d\x91\xc3\x74\xd2\x99\x3c\xec\x09\x41\x44\x6d\x32\x22\x31\x58\xbc\x74\x2c\xc0\x32\x0e\xe0\x38\x8e\x4b\x9c\x4e\xff\xbd\x23\x33\x69\x9f\x6e\x67\xef\x76\x75\xba\xfd\xed\xb5\xc3\x5b\xf3\xfa\x73\x58\xf7\x8d\xf7\xdd\x1b\xd7\xc7\xfe\xb5\xf9\x5a\xac\x8f\x7d\x26\xbc\x2f\x5e\xb5\xaf\x1d\x8d\x38\xfa\x30\xca\x07\x00\xea\x34\x87\x05\x84\x46\x74\x85\x18\xe1\x66\x76\xeb\x2a\x6d\xc6\xe8\x4a\x23\xb9\xda\x8c\xe0\x4c\x5e\x66\x22\xb0\xfc\x54\xdb\xeb\xad\x9e\x65\x2f\x86\x05\x9b\x8a\xe5\xa7\x32\xce\xcf\x00\xc2\xb2\x08\x58\xc8\xd2\x5f\x66\xc8\xce\xba\x48\xfd\x15\x4b\xcf\x65\x11\xf9\xa5\x02\x4d\xba\xae\x33\x45\xe7\x67\xab\xb2\x33\x83\x02\xc6\xb4\x33\x83\x4e\x2a\x96\xfb\xcb\x3e\xff\xa8\x8b\x77\x1f\x58\x5c\x7c\x68\x48\xfd\xaa\x7f\xef\x8c\x05\xc7\xcd\x03\xac\x04\x60\x17\x20\xf9\x38\xe1\x95\x5b\xbd\x08\xad\xe4\x00\x04\xa0\x9a\x6f\x13\x5e\xc4\xa6\xc8\xdf\xcc\x2c\x09\x60\x51\xaf\xae\x5f\x1e\xd2\x21\x0d\xaa\x59\xd6\x99\x27\xe1\xeb\x59\xd2\xe9\x55\xd6\x29\x76\x7b\xaa\xef\x13\xf7\x8f\x68\x19\x2f\x40\x2a\xd8\x55\x7d\xbe\xad\xe3\x3c\x34\x2c\xf0\xab\x3e\xda\xd5\x40\xcc\x75\x3c\xbd\x87\x83\x6f\xb5\xf3\xff\xb1\x9d\xaa\xa0\xcd\x98\x40\x29\x22\x70\xcb\x49\xc4\x0a\xc5\xb1\xc5\x4e\x73\xe2\xe0\xea\x89\x0e\x02\xb0\x78\x16\xe3\xba\x05\x49\x24\x74\x03\xc7\x09\x8b\xe7\x7b\x8b\x78\xaf\x88\x8b\x0b\x2f\x42\xcd\xa5\xd3\x64\x50\x14\x6a\x9e\x38\xac\xc0\xb0\xb5\x16\x2d\x28\xb5\x72\x29\x21\x3f\x75\x0c\x5b\x38\x9e\xab\x8c\x0e\xa9\x95\x7b\xe7\x81\x50\x73\xce\xc1\x25\x04\x83\x55\x53\xbf\x01\xd6\x3b\xcd\xe1\xfc\x1e\xe9\x20\x21\xd5\x3f\xbe\x05\x09\xa4\x74\xe9\x95\x20\x37\x47\x38\xce\xdb\x29\x67\x48\xb7\x3b\x21\xe2\xe3\xa9\x99\xee\x2a\xa1\x94\x9a\x13\x27\x77\x0c\xa8\x0c\xf3\xf0\x33\xdf\xe5\xff\xcc\x81\xbb\x3b\xef\xcf\xdf\x00\x00\x00\xff\xff\xee\xe1\x3a\x58\x84\x02\x00\x00")

func systemJsonBytes() ([]byte, error) {
	return bindataRead(
		_systemJson,
		"system.json",
	)
}

func systemJson() (*asset, error) {
	bytes, err := systemJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "system.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _bchainJson = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x7c\x97\x5b\x73\xea\x30\x92\xc7\xbf\xcb\x79\xcd\xd6\x02\x06\x12\xd8\xaa\x7d\x68\x5d\x6c\x64\x6e\x91\x03\xf8\xf2\xb2\xe5\x0b\x08\x63\x1b\x43\x80\x18\xd8\x9a\xef\x3e\x25\xab\x73\xce\x49\x4d\xcd\x3c\x90\xb8\xfc\x53\x77\xff\xfb\xaf\x56\x44\xfe\xff\x57\x7e\xbc\x6e\x3f\xff\xef\x18\x57\xdb\x5f\xff\xf3\xab\x89\x2f\xd5\xe7\xf6\xbf\xfd\xf8\x52\x79\xfc\xd7\x7f\xfd\x4a\xeb\x4c\xbf\x06\xc7\x7e\x26\xb2\x0b\x00\x40\xae\x35\x05\x05\x30\x0d\xc1\xed\xec\x1a\xa2\x20\x18\x91\x9d\x05\xa4\x7d\xee\x8f\x3b\xbb\xc6\xee\x84\xe0\x75\x76\xfd\x71\x07\x1c\x00\x70\x80\xef\xfa\xe3\x97\x96\x1f\xc7\x2f\xe0\x00\xdb\xf5\x87\x9a\x91\x9d\x02\x7a\x0a\x9d\x36\x1f\x58\x9b\x5b\xa6\x7a\xfb\xb0\xef\x9d\x12\x6b\xb8\x4b\xad\xcd\x2d\x72\x36\x0f\x00\xc1\x22\x7f\x68\x91\xca\x7e\xa6\xd6\xe6\x91\x69\x15\xed\xda\x61\x9e\x38\xe3\x43\xdc\x1b\x37\x69\x35\x2e\x32\x7f\x51\xa6\x0a\xa8\x61\x51\x11\x96\xe3\x63\x14\x48\x80\x66\x51\x26\xc7\xd0\x89\x1c\x57\xe7\xec\x02\xc8\x36\x1f\x0f\xdc\x32\x0d\x36\xa7\xb4\xda\x7c\x44\x41\xf6\x95\x56\x1e\xc9\x02\xaf\x06\x58\xb7\x9c\x06\x8b\x32\xb3\xdd\x32\xed\x6f\x2e\x7f\x6a\x96\xcf\x6c\x22\x1e\x89\x33\xee\x47\x81\x00\x50\x6d\x6e\x91\xf6\xbd\x47\x5c\x8d\x4f\xc9\xf1\x4f\x3d\xfd\x2e\x71\x36\x37\x00\xce\x8a\xb0\xf5\x4c\x12\x25\x26\x00\x82\x82\x12\x82\xc8\x82\x12\xb9\xb6\x89\xfc\xe0\x4a\x50\x00\xb0\x1b\xfd\x43\x0d\xda\xb5\x04\x80\x4f\xae\xf5\x54\x3f\x7b\xa4\xea\x95\x89\x3f\x7e\x6c\xa5\x00\xa7\x1c\xc7\xab\x43\x94\x87\xd6\x7e\x1f\xfb\x83\x41\x36\x71\xf7\xc9\x71\x51\x45\x81\x6b\xaf\x9d\xc5\x70\xfd\x24\x3b\x80\x7a\x15\xf4\x4e\xcb\x45\xe5\x1e\x62\xc7\x3e\x25\x87\xe8\x81\x3d\xda\x99\x82\x99\xb3\xf9\x91\xe3\x3e\x0f\xbc\x47\xe8\x0f\x9f\x51\xb5\x79\x78\xd5\xa6\xf4\x82\x02\x18\xc9\x76\x7e\x31\xb0\xc2\x6a\x51\x87\x7e\x79\x5b\xfa\xee\x5e\xf7\x13\x5a\x9b\xf7\xa8\xd8\xc8\xb0\x81\x85\xfd\x53\xcb\x30\xb2\x36\xdd\x75\x7f\xd3\xa4\xce\x7d\xe8\x05\x21\xb0\x7d\xf9\x23\xc7\x7c\x65\x1f\xa3\xc0\xe3\x91\xbf\x38\x25\xbe\x7d\x49\xbb\x1b\x0b\x60\xe4\xff\xd4\x5a\xea\x35\xab\xad\xdf\xcb\x13\xeb\xae\xf5\x4a\xdb\xfe\x51\xa7\xaf\xeb\xac\x2a\xfb\x1a\xad\x37\x16\x10\x2e\xfe\x78\x20\x80\xab\xef\x5e\x25\xf0\xe6\xda\x4d\x2b\xfb\x96\x5a\x51\x99\x16\x51\x19\x49\x58\xd3\xbf\xfb\x50\xb0\xa1\x7e\x56\x66\xf6\xa2\x97\x4e\xc8\x65\x2b\xc1\xa7\x56\x56\x66\xdc\x2b\x43\xab\xbc\x86\xfe\xfd\x09\x24\x75\x4d\x5f\xe5\x35\xac\xc6\x17\x20\x6a\xf2\xbb\xbe\x84\x88\xea\xf9\x15\x7a\x40\xce\x5e\xae\xf7\x91\xb0\x40\x85\x7a\x3f\x39\xd9\x35\x2e\x48\x49\x08\x50\x45\x09\x9f\x48\xa8\x2e\x39\x61\x21\x25\x14\x18\x87\x92\x08\xbe\x27\x9c\x00\x11\xe0\x91\x1a\xdc\x33\x07\x50\x84\x0b\x58\x34\xcb\x10\x3e\x80\xcb\x18\x48\x5c\x4c\x00\x9f\x69\x9c\xe3\xbb\x3d\x08\xa1\xe6\x8a\x70\x27\x04\xbf\x96\xd0\xe5\xc2\x2b\xb8\xe2\xc0\x15\xa1\x05\xab\x80\x2b\x28\x4f\xa4\x01\x90\x3a\x56\x80\x3c\xc1\x50\x11\x01\xee\xf8\x47\xee\x1c\x98\xae\x5d\x81\x7d\x96\x52\x49\x42\x81\x4b\x05\x7d\x68\x94\x02\x0f\x26\x9a\x71\x47\xe7\xe1\xaa\x8d\xf1\xc8\x99\x83\xca\xdb\x67\xb9\x3f\x73\x98\xc7\x02\x64\x5b\x93\x52\xa0\x3b\x96\xea\xdf\xed\x3a\xbd\x13\x84\x73\x11\x17\x5c\x70\xf0\x08\xe8\xf7\x19\xd0\xf8\x10\x52\x02\x97\xd9\x60\xd0\x9e\x7b\x0a\xc7\x11\xdd\xa9\x1f\x3e\x01\xf8\x97\x1c\xf4\x3a\x02\xb6\x00\x97\x14\x60\x9f\x39\x08\x05\x85\xf6\x66\x2e\x65\x06\x44\x7f\x38\x14\x5a\x0b\x35\xba\xb4\x16\xa1\x20\x2f\x18\xe7\x82\x68\x2d\x5c\x7b\xfe\xa1\xfd\x75\xdb\x3e\x28\xd1\x5e\x88\x13\xdc\x61\x0f\x0f\xa0\xbf\x7b\x2d\xb8\x44\x9f\xa9\x8c\x5a\x5f\xe8\x5f\xbd\xd2\xef\x5e\x65\xcd\xe0\x2d\x9c\xc5\xa3\x41\x97\x09\x29\xa9\xa4\xc0\xa4\xb2\x3d\xa6\x9d\xa3\x5c\xef\x1c\xed\xbc\xb0\x4f\xbb\x17\x2e\x40\xb3\x13\xdd\xb3\x01\x3d\x4a\xb9\xe8\xca\xad\xce\x5b\x30\xbd\x8e\x3e\x80\xb6\x3a\x1a\xa5\xc8\xa4\xf9\xa3\x21\x37\x5a\x29\x6a\x95\x12\x14\x68\x5d\x7e\x6d\x81\x92\x33\x3a\x7e\x21\xed\xac\x49\xb7\x03\x05\x21\xc0\x39\x4c\x41\x70\x19\x02\x8b\x1f\x82\x2d\x94\xe0\x02\x6c\x0a\xb6\x4b\x06\x66\xef\x59\xeb\xf9\x5e\xd7\x5b\xfc\x9e\x31\xe8\xe2\x8c\xe5\xc0\xda\x7e\x4d\xef\x0c\x67\xcc\xbc\x6b\x67\x4c\x28\xe8\x3a\x21\x54\xdf\x33\xe6\x00\x78\x5a\xe3\x5c\xcf\x13\x15\x60\xc7\x7a\x8e\xf4\xdc\xc1\x03\xd8\x54\xb6\x39\xc9\x6e\x95\xb2\x0a\xda\x58\x09\x7e\xad\x4c\x1d\x45\xcd\x2c\x48\xac\x3f\x17\x5a\xd3\xef\xf9\xea\xea\xf3\x51\x2b\xe8\x72\xf9\xa7\x36\x17\xf1\x1e\x98\xd3\xce\xe9\xbc\xf5\x0b\x02\xf4\xab\xdb\xfa\xd5\xf6\xf5\xc3\xaf\xae\x03\x60\xa1\x5f\xa3\x87\x39\x9b\xd2\xfe\xe9\x97\xee\xef\x21\xa8\xf1\xcb\xa5\xc0\x52\x96\x6b\x0d\xed\xcc\xbb\xdf\xf3\x4e\x39\xe8\x7d\x71\x44\x7b\xc6\xee\x35\x3e\xb7\xb5\x1c\xb3\x37\x6d\x7f\xa6\xd6\xcb\xcd\xd4\x52\x8b\x0e\x04\x03\x2a\x39\x07\x29\xa1\x06\xa5\xcf\xd0\x25\x87\x67\x48\x09\x9f\x42\x73\x4d\xbf\x84\xd4\x1e\x98\x67\xed\x83\x94\xb0\x87\xf6\x0c\x40\x0d\x9f\x0d\x11\xca\xb6\x99\x04\x26\xdb\xf3\x78\x7a\x69\xd7\x17\x9c\xc4\xb9\xd0\x9e\xbb\x4c\xc2\x4c\xc0\x5c\x11\x7d\xde\xe0\x62\x7c\x6a\xbd\x8d\x48\x3b\x77\x1e\x39\xeb\xde\x68\x73\x73\xcd\xdd\x42\xc8\x17\xd8\x84\x10\x70\xe8\x00\x3d\xd1\xe7\x6c\x3d\xd1\xfb\xe9\x72\x08\x67\x07\xd1\xde\x3b\x20\x63\x58\xea\x0f\x87\x42\x12\x75\x9b\x9b\x78\xd0\x37\x2b\xe8\x33\x0d\x7b\x98\xd2\x41\xfe\xdb\x57\x01\x20\x69\x73\x9b\x9a\x75\x9c\xec\x1e\x7a\x3e\x40\x32\xb8\x89\xbf\x63\x39\x74\x67\xb9\xc0\xbb\x2d\x18\x29\x20\xb0\xa4\x83\x0a\xb5\x38\x23\x58\x73\xfe\x43\x47\x0a\x44\x7f\x7e\xea\x20\xa4\x86\x21\xa9\x61\x20\xa9\x07\x0e\x6d\x6e\xa3\x43\xad\xdf\x37\x20\xa5\x9c\x11\x46\x26\xba\xa0\xf4\x60\xc6\x1c\x7d\x67\x8e\x54\x14\xec\x0f\x91\xbf\x29\x74\xac\x80\x4b\xf6\xd7\x5d\xa1\xf4\xdf\x71\x31\x59\x94\x49\xe5\x95\x69\x4e\xa8\x7c\x90\x32\x3d\xba\x5f\xa9\x6a\xeb\x03\xbd\x97\x7f\xee\x8f\x9c\x9c\x92\xe3\xa2\x17\x55\xd1\x29\xb4\xca\x32\x39\x4a\x25\x8b\xb9\xf6\x51\x01\x5c\x36\x7f\xdd\xc3\x4e\xe4\xaf\x55\x54\x6d\x4a\xe1\x94\x4f\xc1\xda\xfe\x00\xe4\x45\xfd\xcb\x9a\xff\x54\x3f\x06\x42\xef\xf7\x9f\xf7\xd7\xc7\xbf\xd3\xe0\x8c\x40\x5e\x78\xeb\x36\x10\x3d\xbf\x33\xe8\xba\x4c\x3f\x77\x81\xcf\xe8\x66\xfc\x58\x1d\x37\xd7\xb0\xd2\xdf\xa7\x08\x19\x00\x9f\x35\xd0\x40\x73\xbb\xcc\xda\x10\x41\x8f\x43\x79\x31\x8f\x2e\xdb\x51\xfd\x4d\x04\xd8\x76\x3b\x2e\x0c\x9f\x95\xe7\x6e\x63\xf8\x67\xef\xb9\x30\x9c\x76\xd6\x69\x6a\xf8\x32\x73\x1b\x8c\x77\x0b\xff\xcd\x70\x32\x94\x6f\x6b\xe4\x84\x7f\x62\xbc\xcf\xab\x13\xe6\xdf\xaf\xa7\x73\xc3\x27\xfb\xe4\x88\xf1\x13\xa7\x0a\x90\xcf\xb7\x11\x47\x7d\xfe\x22\xc7\x78\xef\xf5\xc3\x36\x1c\xe4\xd7\x71\x64\xf8\xe2\x2d\xde\x62\xfc\x67\xe7\xf4\x44\x7d\x6c\xe2\x0f\x0c\xb7\x33\x16\x83\xe1\xd3\xea\x9e\x63\xfe\xb7\xb2\xd7\x18\x3e\x3f\x38\x1b\x65\xf8\xa1\xb3\x97\xa8\xff\x7d\xb6\x46\x7f\x20\x2b\x3d\xcc\x7f\x8d\x46\x1d\xe4\x41\x7f\x85\xfe\x2c\x9b\x6c\x81\xf9\xd3\x8a\xdf\x30\x7f\x7d\x7a\x28\xc3\xf9\xf3\x2e\x30\xff\x71\xd5\xcd\x90\xc7\x07\x82\xfe\x01\xdb\x71\xcc\xcf\x9e\x72\x8e\xfd\xd9\x07\x08\x31\x3f\xab\x00\xf3\x1f\xbb\xab\x37\xaa\xda\xfa\x70\xda\xad\xa7\x46\xbf\x10\x23\xa8\x8d\xff\x4e\xe7\x8c\xbc\xb0\x96\xd2\x70\x41\x78\xbf\x31\x7c\x3b\x7b\x8f\x90\xbb\x6e\x86\xfc\x7d\x1a\x3d\x95\xe1\xe9\x7e\xe2\x1a\xce\x26\xf9\x68\x8e\x1c\xc4\x17\xf2\xd7\x41\x3e\x30\x9c\xf6\xe9\x1b\x72\x7b\xe9\x7d\x4a\xc3\xbf\xbc\xaa\xc6\xf8\xe7\x69\x2a\x30\x7e\xfd\x38\x22\x6f\xe6\xcb\x00\xeb\xaf\x37\x43\xe4\xd3\xb7\xeb\x01\xf5\x3f\xeb\x78\x82\xf9\x1f\xa2\x46\x6e\x1f\xc9\x0e\x79\xe3\xbb\x7d\xcc\x7f\x22\x47\xe4\xf3\x97\x6b\x82\x3c\x2d\xb3\x23\x72\x09\x16\xf2\x65\xf5\x08\x90\xdf\x2a\xef\xbb\x3e\x73\x36\xc8\xc5\x4c\xae\x90\x67\xe1\x09\xeb\x43\xea\x75\x91\x4f\xd6\xce\x3b\xea\xdf\x5d\x14\xd6\x07\xeb\xb0\x42\xee\xde\xd7\x33\xe4\x9f\xfb\x01\xf6\x4f\x3e\x5f\x1f\xc8\x17\xd3\xd7\x09\xf2\x41\xd8\x84\xd8\x5f\x10\x7a\xdf\xfe\x8a\x01\x43\x7f\x47\x7b\x81\xfe\xd3\xf5\x5b\x83\x7c\x71\x58\x76\x54\x61\xfc\x9d\xb1\x21\x95\xed\x63\xa7\xdc\x49\xd7\xf8\x33\x1e\xbd\x35\x86\x3b\x56\x72\x36\x9c\x3c\xd6\xb7\xb5\xe1\x8e\xfe\x3e\x60\xf8\xeb\x73\x96\x60\xfc\xd4\x1d\x87\x86\xcf\x46\x59\x57\x1a\x4e\xc9\x6e\x6e\x38\x2b\xf8\x25\xc5\xfc\xdc\x6b\xb0\xfe\x6a\x3f\x1d\x23\x8f\xed\x44\x19\xfe\x2e\x93\x0b\xd6\xaf\xf2\x8f\x1b\xf2\xd9\xec\xb5\x30\x5c\xac\xc7\x27\xac\xff\xe0\xd6\xce\x70\xfa\x11\xba\xb5\xe1\x4b\xb7\x57\x62\xfd\x6b\x6d\x49\xd4\xb7\x38\x8b\x0b\xf2\x84\xef\xb1\xbe\x5c\x7e\x10\xe4\x11\x39\x76\x0d\x9f\xc0\x79\x8b\xf9\x3d\xdb\x7e\x62\xfe\xd5\xb6\x1a\xa0\xfe\x47\x15\x63\xfe\xe3\xd5\x2b\xd0\x9f\x90\xf8\x80\xfa\xd4\xc0\x47\xfd\xb3\x6a\xe8\x23\x17\xa7\x92\x1b\xbe\x58\xa9\x15\xe6\xef\xee\x7b\x02\x39\x04\xeb\xb9\xe1\xc4\xda\xbd\xa3\xbe\xe8\x65\xfe\xbd\x3f\xdb\xe9\x01\xfd\xe7\xe4\x39\xc7\xf8\xe7\x6c\xf0\x89\xfe\x00\xa3\xe8\x2f\xb9\x67\x02\xe3\xfb\xe7\x51\x86\xf1\x36\xf1\xd0\x3f\x97\xa4\x36\xc6\xef\x92\xf0\x1d\xe3\xe7\x5c\xa1\x3f\x93\xeb\x9d\x60\x7c\x9e\x09\xf4\x87\x75\x17\xbb\xae\x30\x4d\x3d\x92\x17\x30\x43\xdb\x6f\xde\xbb\x78\x4f\xa4\xc9\x6d\x64\x38\x9b\x25\xaf\xca\xf0\xa0\xfb\x59\x19\xce\x0e\xb7\x17\x8e\xf1\xcd\x57\x5f\x1a\xfe\x1a\x07\x21\xc6\x9f\x26\x57\x89\xf1\xb7\xe8\xd9\x60\xbc\x97\xcc\x30\x7e\x59\xa6\xa1\xe1\x62\x10\x34\x98\xbf\x1a\xde\xc6\xc8\x9f\x22\x29\x0c\x77\xe3\xf3\x15\xf5\xa9\x83\xdf\x18\x4e\xad\xe6\x7c\xc1\xfc\xee\xc7\x09\xf3\x8f\x5e\x96\x39\xc6\x0f\xb3\xe1\xc0\x70\x3e\x5f\x56\x98\xff\xf9\x1a\x6f\x30\x3e\x95\x11\xea\x9f\x7e\x6e\x0f\x98\x7f\x7c\xb4\x04\x72\x98\x0c\x50\x3f\xdf\xb1\x1d\xe6\xff\xf0\xcb\x37\xe4\x0e\x0f\x52\xc3\xdf\x1f\x6f\x29\xe6\x3f\x57\xfb\x1b\xf2\x77\xdb\xaa\xbf\xf3\x8b\x08\xfd\x79\x39\x5e\xf7\xe8\xcf\x78\xbe\x41\xff\xd9\xfb\xc9\xc7\xfa\xce\x79\xb1\x36\x1c\xee\x41\x0f\x0c\x9f\xa4\xd9\x0a\xf9\xeb\xb8\x37\xc1\xfc\xf3\x9a\xcf\x0d\xb7\xcb\xf2\x1d\xf5\xad\x69\x83\xfa\xc8\x7d\xbc\x42\x7f\x67\x9f\x30\x47\x7d\xd3\x7a\x80\xfa\x18\x04\x1f\xa8\x6f\xd2\x7c\xb9\xc8\xfd\x66\x8e\xfe\x91\xcc\xf2\x50\x9f\x70\xee\x0e\xea\xbf\xd7\x17\xf4\x0f\xce\x31\xf9\xf6\xaf\xe7\x30\xe4\xab\x8f\xd4\x45\x3e\xa5\x6a\x6d\xf8\x72\x76\x03\xd4\x0f\xdc\x51\x4c\x5d\x96\x32\x27\xac\xfd\x7f\x3e\x27\xeb\xc4\xba\x96\x89\x82\xff\xfd\xf5\x8f\x7f\x06\x00\x00\xff\xff\x18\x8b\x9a\xc0\x70\x12\x00\x00")

func bchainJsonBytes() ([]byte, error) {
	return bindataRead(
		_bchainJson,
		"bchain.json",
	)
}

func bchainJson() (*asset, error) {
	bytes, err := bchainJsonBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "bchain.json", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"pledge.json": pledgeJson,
	"system.json": systemJson,
	"bchain.json": bchainJson,
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
	"bchain.json": &bintree{bchainJson, map[string]*bintree{}},
	"pledge.json": &bintree{pledgeJson, map[string]*bintree{}},
	"system.json": &bintree{systemJson, map[string]*bintree{}},
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

