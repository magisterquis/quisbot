// Code generated by go-bindata.
// sources:
// templates/footer.tmpl
// templates/header.tmpl
// templates/index.tmpl
// DO NOT EDIT!

package main

import (
	"github.com/elazarl/go-bindata-assetfs"
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

var _templatesFooterTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb2\xd1\x77\xf2\x77\x89\xb4\xe3\xb2\xd1\xf7\x08\xf1\xf5\xb1\xe3\x02\x04\x00\x00\xff\xff\xee\xe8\x85\xfb\x10\x00\x00\x00")

func templatesFooterTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesFooterTmpl,
		"templates/footer.tmpl",
	)
}

func templatesFooterTmpl() (*asset, error) {
	bytes, err := templatesFooterTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/footer.tmpl", size: 16, mode: os.FileMode(420), modTime: time.Unix(1472009823, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesHeaderTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xb2\x51\x74\xf1\x77\x0e\x89\x0c\x70\x55\xf0\x08\xf1\xf5\xb1\xe3\xb2\x81\x51\xae\x8e\x2e\x76\x5c\x0a\x50\x60\x13\xe2\x19\xe2\xe3\x6a\x57\x5d\xad\xa0\x17\x92\x59\x92\x93\xaa\x50\x5b\x6b\xa3\x0f\x11\xe3\xb2\xd1\x87\xa8\xb5\x71\xf2\x77\x89\xb4\xe3\x02\x04\x00\x00\xff\xff\x99\xc1\x68\xec\x51\x00\x00\x00")

func templatesHeaderTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesHeaderTmpl,
		"templates/header.tmpl",
	)
}

func templatesHeaderTmpl() (*asset, error) {
	bytes, err := templatesHeaderTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/header.tmpl", size: 81, mode: os.FileMode(420), modTime: time.Unix(1472009834, 0)}
	a := &asset{bytes: bytes, info: info}
	return a, nil
}

var _templatesIndexTmpl = []byte("\x1f\x8b\x08\x00\x00\x09\x6e\x88\x00\xff\xaa\xae\x56\x28\x49\xcd\x2d\xc8\x49\x2c\x49\x55\x50\x82\xb1\x8a\xf5\x33\x52\x13\x53\x52\x8b\xf4\x4a\x72\x0b\x72\x94\x14\x6a\x6b\xb9\x6c\x3c\x0c\xed\xc2\x53\x73\x92\xf3\x73\x53\x6d\xf4\x3d\x0c\xed\xb8\x6c\x02\xec\x42\x32\x32\x8b\x15\x32\x8b\x15\x12\x15\xf2\x32\x93\x53\x15\x0a\x12\xd3\x53\x6d\xf4\x03\xec\xb8\xb8\x70\x98\x99\x96\x9f\x5f\x82\x6c\x26\x20\x00\x00\xff\xff\x91\x85\xa4\x86\x7b\x00\x00\x00")

func templatesIndexTmplBytes() ([]byte, error) {
	return bindataRead(
		_templatesIndexTmpl,
		"templates/index.tmpl",
	)
}

func templatesIndexTmpl() (*asset, error) {
	bytes, err := templatesIndexTmplBytes()
	if err != nil {
		return nil, err
	}

	info := bindataFileInfo{name: "templates/index.tmpl", size: 123, mode: os.FileMode(420), modTime: time.Unix(1472016191, 0)}
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
	"templates/footer.tmpl": templatesFooterTmpl,
	"templates/header.tmpl": templatesHeaderTmpl,
	"templates/index.tmpl": templatesIndexTmpl,
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
	"templates": &bintree{nil, map[string]*bintree{
		"footer.tmpl": &bintree{templatesFooterTmpl, map[string]*bintree{}},
		"header.tmpl": &bintree{templatesHeaderTmpl, map[string]*bintree{}},
		"index.tmpl": &bintree{templatesIndexTmpl, map[string]*bintree{}},
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


func assetFS() *assetfs.AssetFS {
	assetInfo := func(path string) (os.FileInfo, error) {
		return os.Stat(path)
	}
	for k := range _bintree.Children {
		return &assetfs.AssetFS{Asset: Asset, AssetDir: AssetDir, AssetInfo: assetInfo, Prefix: k}
	}
	panic("unreachable")
}