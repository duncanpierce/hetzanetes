// Code generated for package tmpl by go-bindata DO NOT EDIT. (@generated)
// sources:
// cloudinit.tmpl.yaml
// get-private-interface.tmpl.sh
package tmpl

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

// Name return file name
func (fi bindataFileInfo) Name() string {
	return fi.name
}

// Size return file size
func (fi bindataFileInfo) Size() int64 {
	return fi.size
}

// Mode return file mode
func (fi bindataFileInfo) Mode() os.FileMode {
	return fi.mode
}

// Mode return file modify time
func (fi bindataFileInfo) ModTime() time.Time {
	return fi.modTime
}

// IsDir return file whether a directory
func (fi bindataFileInfo) IsDir() bool {
	return fi.mode&os.ModeDir != 0
}

// Sys return file is sys mode
func (fi bindataFileInfo) Sys() interface{} {
	return nil
}

var _cloudinitTmplYaml = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\xdc\x56\x4d\x6f\xdb\x46\x13\xbe\xeb\x57\x0c\xf8\x1e\x0c\x04\x59\x32\x89\x8d\xb7\x85\x00\xb5\x08\x2a\x17\x49\x6d\xd8\x81\xed\xb4\x87\xa2\x10\x46\xe4\x90\xda\x6a\xb9\xbb\xd9\x1d\x4a\x55\x18\xfd\xf7\x82\x5f\x96\x44\x4b\x96\x8a\xb6\x68\x51\x5e\x04\xcd\xce\xcc\x3e\xf3\xb9\xcf\xff\x62\x65\x8a\x44\xc4\x46\xa7\x32\x1b\x0c\x2c\xc6\x73\xcc\x68\x52\xd8\x04\x99\x86\xc0\xae\xa0\x2d\x61\xe6\x30\xe9\x49\xfd\x70\x00\x20\xe0\xe7\xa0\x48\x97\xc1\xcb\x20\x2e\x9c\x0a\x5e\x06\xbf\x7e\x0a\x7e\x79\xb4\x73\x34\x35\x86\x27\x32\x9d\x38\xfa\x54\x48\x47\x49\xeb\x62\xe9\x24\xd3\x24\x95\xaa\xf3\x62\x91\x67\x43\x08\xca\x32\x7c\xaf\x3d\xa3\x52\x63\xe9\x28\x66\xe3\x56\xeb\x75\x34\x2f\x3c\x9b\x5c\x7e\x46\x96\x46\x87\x2b\xcc\x55\x30\x00\x00\x88\x8d\x66\xd2\x3c\x84\x2f\xf5\x5f\x00\xb4\xf2\x47\x72\x5e\x1a\x3d\x84\xce\x88\xc2\x26\xc4\x70\xfe\xb5\x0f\xa5\x89\x16\xaf\xa7\xc4\xf8\xba\xb5\x98\x4b\x9d\x0c\xe1\x6a\xfb\x82\xf6\xc4\x91\x37\x85\x8b\x1b\x80\xcd\x27\x60\x46\xfc\x59\x93\x13\x9e\x62\x47\x1c\x72\x6e\xd5\x9e\xd3\x38\xce\x6b\x94\xfb\x8e\xbc\xdc\x3e\xb2\xc8\xf1\x6c\xf7\x8a\x5a\xd4\xe9\x1c\x49\x4c\x0f\xce\x89\x99\x59\xec\x06\x7f\x5f\x1b\xb7\xa2\x9c\x18\x13\x64\xdc\x20\xd2\x98\xd3\x10\x66\x75\xb7\xec\x08\xbd\xc5\x98\xaa\x3c\x4f\x49\xf8\x95\x67\xca\xdb\x63\xcf\x4e\xea\x6c\xbc\xe3\x85\xcd\x9c\x74\x13\xc7\xbb\x06\xf4\x5b\x2b\x1f\x2a\xe1\x7a\x1d\x1c\x0f\x74\x93\x94\x13\xe2\x43\x6b\x7d\xd4\x0b\x72\x4c\x56\x99\x55\x4e\xfa\xb4\x40\xc5\xe3\x70\xb0\x33\x4a\x91\x13\x39\x6a\xcc\xc8\x9d\x98\x01\x4b\xf1\x56\xec\x94\x5b\x55\x0d\xd5\xa3\xa4\xaf\xd1\x85\x84\x52\x93\xf3\xbb\xf2\x2a\x37\x7f\x0c\xda\xc6\x63\x9e\xa3\x4e\xfa\xee\x1a\x97\x41\x34\x95\x3a\x3a\xe2\x32\xd8\x6f\x2a\x5a\x03\xeb\xcc\x42\x26\xe4\x46\x8d\x9b\x83\xda\x8a\x30\x21\x27\x48\x51\xcc\xa3\x14\x95\xa7\x83\xaa\xa8\x94\x59\x8a\x42\x33\x66\x19\xb5\xc0\x9e\x55\x8e\x91\x49\x68\x93\x90\x88\x65\xe2\xfc\xa8\x5a\x2f\xcf\xc0\x2e\x3c\x57\x53\x28\x13\x37\x2a\xcb\xf0\x83\x49\xde\xdb\x3b\xd4\x19\x35\x5d\xb8\xfb\x91\x5e\xec\xcf\x5d\x53\x8e\x77\xdf\x5d\xdf\x7e\x1c\x4f\x6e\x2e\x1f\x7e\xba\xbd\xbb\xda\xa3\x08\xb0\x40\x55\x50\xd3\xd4\x1f\x9c\x5c\x20\xd3\x0d\xf1\xd2\xb8\xf9\x0d\xe6\x7b\xaf\x7c\x34\xfa\xde\x99\xbc\xbb\x5c\x08\xb1\x3b\xb0\x8c\x4c\x69\xa1\xee\x1f\xa7\xf6\x70\xf7\x1f\x6b\x73\x2f\xb7\xca\xfe\x4f\x36\x77\x85\xa4\xee\xa7\x2a\x8a\x3d\xcd\x7c\x5a\x31\x1e\x6e\xaf\x2e\x6f\x4e\xcb\x6a\xff\x6b\x76\xe8\x15\xad\xee\x28\x3d\xa4\x73\x60\x1b\xf6\xbf\x39\xad\x86\xcd\xc6\x3b\x50\xc1\x31\x52\x6e\xf4\x5f\x54\xbf\xaa\xf9\xff\x0d\x6b\xc9\x4b\x91\x38\xb9\xf8\xef\xd4\xce\x15\x3a\xce\x93\x86\x9f\x14\xe9\x12\xea\xdd\x04\xd6\x19\x36\xc0\xb1\x85\xd4\x99\x1c\x50\xaf\x80\x4d\xfd\x63\x8d\x63\x78\xf3\xa6\x2c\x65\x0a\xe1\x5b\x2b\xef\xc9\x2d\xc8\xad\xd7\x2f\xff\x7f\x71\x71\x5e\x96\xa4\x93\xf5\xba\xe7\xac\x76\xd1\x5b\x44\x07\x54\x2a\x6f\x32\xa6\x63\x6a\xed\x9e\x79\xaa\x26\x52\x48\x28\xc5\x42\x31\x24\xa4\x57\x20\x75\x6c\x72\xa9\xb3\x7d\x0a\x8d\x4b\x53\x70\x66\x7a\x1a\xa4\x71\xaa\x68\xf0\x24\xc6\x01\x40\x59\x46\x2f\xe0\xe1\x76\x7c\x0b\x3c\x23\x47\x20\x3d\xa0\xf2\x06\x84\x48\x15\x6a\x4d\xaa\xe6\x9a\x30\x84\x65\x77\x2c\x19\x50\x27\xb0\x9c\x21\x9f\x79\x90\x1a\x24\x7f\x0b\x2f\xa2\x16\x74\xcd\x28\x41\xf8\xf4\x1a\xce\x66\xcc\xd6\x0f\xa3\x28\x23\x0e\xe7\xe7\x15\x8f\x3b\x83\x2f\xe0\x67\x20\x3c\x08\x10\x22\x91\xbe\xc2\x05\xbe\x49\x92\x9a\x6e\xc9\xaa\x37\x42\x09\xcf\xc6\x61\x46\x1b\xf9\x93\x47\x0f\x84\xa8\x46\x46\x11\x0b\x74\x19\xf4\x9e\x38\xfa\x8d\xc9\x69\x54\x5b\xe1\xc8\x14\x63\x1a\x95\x65\x37\x48\x10\x64\xc4\xc2\x36\x05\x10\x52\x33\xb9\x4a\x23\x80\xb0\xe3\x38\xfb\x42\x72\xb8\x0c\x33\xc9\xb3\x62\x5a\x78\x72\x2d\xb1\x09\x63\x93\x77\xec\xae\x06\x72\xec\xa5\x8e\x72\xac\x5e\xb6\x28\xa9\x79\x4e\xb4\x78\x1d\x7e\x15\xbe\x12\xba\x79\x6d\x7c\xcd\x9c\xce\xe0\x1b\x78\x9e\x46\x76\xbc\xf5\x4f\x83\xdd\xec\x82\x0a\xc9\x45\xf8\xaa\xc3\x55\x25\xd8\x69\x62\xf2\xd1\x66\x6d\x84\xab\x93\xc0\xb5\xcc\xb9\x05\x87\xd6\xaa\x15\x88\xf9\x7e\xab\x60\x50\x96\xa4\x3c\x9d\xdc\x4a\x57\xe7\xf7\x93\x8f\x77\xd7\xa3\xee\xb4\x2c\xab\xee\xbe\xd4\x89\x35\x52\xf3\x7a\x3d\xac\x66\xb8\xd6\xaa\xd7\x54\x45\x1e\x7e\x30\x52\xb7\x0c\x76\xab\x13\xff\xce\x16\xea\x76\xc8\xef\x01\x00\x00\xff\xff\x40\x66\xa4\x44\xbb\x0d\x00\x00")

func cloudinitTmplYamlBytes() ([]byte, error) {
	return bindataRead(
		_cloudinitTmplYaml,
		"cloudinit.tmpl.yaml",
	)
}

func cloudinitTmplYaml() (*asset, error) {
	bytes, err := cloudinitTmplYamlBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "cloudinit.tmpl.yaml", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _getPrivateInterfaceTmplSh = []byte("\x1f\x8b\x08\x00\x00\x00\x00\x00\x00\xff\x52\xd1\xc8\x2c\x50\xd0\xcd\x52\x28\xca\x2f\x2d\x49\x55\xc8\xc9\x2c\x2e\x51\xa8\xae\xd6\x0b\x28\xca\x2c\x4b\x2c\x49\xf5\x2c\x08\x4a\xcc\x4b\x4f\xad\xad\x55\xa8\x51\xc8\x2a\x54\xd0\x2d\x52\xd0\x8b\x36\x88\xd5\x4b\x49\x2d\xd3\x04\x04\x00\x00\xff\xff\x1f\x13\x6f\x42\x38\x00\x00\x00")

func getPrivateInterfaceTmplShBytes() ([]byte, error) {
	return bindataRead(
		_getPrivateInterfaceTmplSh,
		"get-private-interface.tmpl.sh",
	)
}

func getPrivateInterfaceTmplSh() (*asset, error) {
	bytes, err := getPrivateInterfaceTmplShBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "get-private-interface.tmpl.sh", size: 0, mode: os.FileMode(0), modTime: time.Unix(0, 0)}
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
	"cloudinit.tmpl.yaml":           cloudinitTmplYaml,
	"get-private-interface.tmpl.sh": getPrivateInterfaceTmplSh,
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
	"cloudinit.tmpl.yaml":           &bintree{cloudinitTmplYaml, map[string]*bintree{}},
	"get-private-interface.tmpl.sh": &bintree{getPrivateInterfaceTmplSh, map[string]*bintree{}},
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
