package test

import (
	"path"
	"runtime"
)

// GetTestPath test path
func GetTestPath() string {
	_, filename, _, ok := runtime.Caller(1)
	if !ok {
		panic("GetCurrentFilePath failed")
	}
	defaultPath := path.Dir(filename)
	return defaultPath
}
