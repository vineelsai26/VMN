//go:build linux || darwin
// +build linux darwin

package node

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func SetPathWindows(path string) {}

func SetPathLinux(path string) {
	exec.Command("export", "VMN_VERSION="+path).Run()

	shells := []string{".bashrc", ".zshrc"}

	home := GetHome()

	for _, shell := range shells {
		if _, err := os.Stat(filepath.Join(home, shell)); err == nil {
			f, err := os.OpenFile(filepath.Join(home, shell), os.O_RDONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := f.Stat(); err == nil {
				b := make([]byte, 1024)
				n, err := f.Read(b)
				if err != nil {
					panic(err)
				}
				if !strings.Contains(string(b[:n]), "export PATH=$VMN_VERSION:$PATH") {
					exec.Command("echo", "export PATH=$VMN_VERSION:$PATH", ">>", filepath.Join(home, shell)).Run()
				}

				if !strings.Contains(string(b[:n]), "export VMN_VERSION=${cat "+filepath.Join(home, ".vmn", "current")+"}") {
					exec.Command("echo", "export VMN_VERSION=${cat "+filepath.Join(home, ".vmn", "current")+"}", ">>", filepath.Join(home, shell)).Run()
				}
			}
		}
	}
}
