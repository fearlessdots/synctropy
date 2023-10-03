package main

//
//// IMPORTS
//

import (
	// Modules in GOROOT
	"os"
	"path/filepath"
	"strings"
	// External modules
)

//
//// FILES/DIRECTORIES
//

func fileIsHidden(name string) bool {
	_, file := filepath.Split(name)
	return strings.HasPrefix(file, ".")
}

func filterHiddenFilesAndDirectories(unfilteredFiles []os.FileInfo) []os.FileInfo {
	var filteredFiles []os.FileInfo

	for _, element := range unfilteredFiles {
		if !fileIsHidden(element.Name()) {
			filteredFiles = append(filteredFiles, element)
		}
	}

	return filteredFiles
}
