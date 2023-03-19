//go:build !linux && !darwin
// +build !linux,!darwin

package setup

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
	"vineelsai.com/vmn/node"
)

func Install() {
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	srcFile, err := os.Open(filepath.Join(dir, "vmn.exe"))
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	if _, err := os.Stat(filepath.Join(node.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(node.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	destFile, err := os.Create(filepath.Join(node.GetHome(), ".vmn", "vmn.exe"))
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		panic(err)
	}

	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	userPath, _, err := k.GetStringValue("Path")
	if err != nil {
		panic(err)
	}

	var isPathVariableExists bool = false

	for _, env := range strings.Split(userPath, ";") {
		if env == filepath.Join(node.GetHome(), ".vmn") {
			isPathVariableExists = true
			return
		}
	}

	if !isPathVariableExists {
		exec.Command("setx", "PATH", userPath+filepath.Join(node.GetHome(), ".vmn")).Run()
	}
}
