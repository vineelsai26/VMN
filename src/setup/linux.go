//go:build (linux && ignore) || (darwin && ignore) || !windows
// +build linux,ignore darwin,ignore !windows

package setup

import (
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"vineelsai.com/vmn/src/utils"
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

func SetPath(path string) {
	exec.Command("export", "VMN_VERSION="+path).Run()

	shells := []string{".bashrc", ".zshrc"}

	home := utils.GetHome()

	for _, shell := range shells {
		if _, err := os.Stat(filepath.Join(home, shell)); err == nil {
			f, err := os.OpenFile(filepath.Join(home, shell), os.O_RDONLY, 0644)
			if err != nil {
				panic(err)
			}
			defer f.Close()

			if _, err := f.Stat(); err == nil {
				b := make([]byte, 1024*1024)
				n, err := f.Read(b)
				if err != nil {
					panic(err)
				}

				if !strings.Contains(string(b[:n]), "eval \"`vmn env`\"") {
					file, err := os.OpenFile(filepath.Join(home, shell), os.O_APPEND|os.O_WRONLY, 0644)
					if err != nil {
						panic(err)
					}
					defer file.Close()

					if _, err := file.WriteString("\n#VMN \neval \"`vmn env`\""); err != nil {
						panic(err)
					}
				}

			}
		}
	}
}
