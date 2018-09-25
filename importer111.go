// +build go1.11

package main

import (
	"fmt"
	"go/importer"
	"go/types"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

const (
	vendorDir = "/vendor/"
)

func moduleLookup(path string) (io.ReadCloser, error) {
	pack, err := os.Open(path)
	if err == nil && pack != nil {
		return pack, nil
	}
	if err != nil {
		fmt.Println("first open error: ", err.Error())
	}

	// try looking up package in vendor
	packageRoot := getPackageDirectoryPath()
	pathParts := []string{packageRoot, vendorDir, path}
	vendorPath := strings.Join(pathParts, "")
	fmt.Println("vendorPath: ", vendorPath)
	mod, err := os.Open(vendorPath)
	if err != nil {
		return nil, err
	}
	return mod, nil
}

type moduleImporterFrom struct{}

func (mif *moduleImporterFrom) Import(path string) (*types.Package, error) {
	pathParts := strings.Split(path, "/")
	var dir string
	pathDepth := len(pathParts)
	if pathDepth <= 0 {
		return nil, fmt.Errorf("Length of import path cannot be 0.  Invalid importation path.")
	}
	if pathDepth == 1 || pathDepth == 2{
		dir = pathParts[0]
	} else {
		dir = pathParts[len(path) - 3]
	}
	return mif.ImportFrom(path, dir, 0)
}

func (mif *moduleImporterFrom) ImportFrom(path string, dir string, mode types.ImportMode) (*types.Package, error) {
	// yet to be implemented
	return nil, nil
}

func getPackageDirectoryPath() (packageRoot string) {
	_, b, _, _ := runtime.Caller(0)
	return filepath.Dir(b)
}

func defaultImporter() types.Importer {
	return importer.For(runtime.Compiler, moduleLookup)
}
