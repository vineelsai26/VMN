//go:build (linux && ignore) || (darwin && ignore) || !windows
// +build linux,ignore darwin,ignore !windows

package setup

import (
	"io"
	"os"
	"path/filepath"

	"vineelsai.com/vmn/utils"
)

func Install() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	srcFile, err := os.Open(filepath.Join(dir, "vmn"))
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	destFile, err := os.Create(filepath.Join(utils.GetHome(), ".vmn", "vmn"))
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		panic(err)
	}
}
