//go:build (linux && ignore) || (darwin && ignore) || !windows
// +build linux,ignore darwin,ignore !windows

package setup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func Install() {
	fmt.Println("Installing VMN...")
	dir, err := os.Getwd()
	if err != nil {
		panic(err)
	}

	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	username := user.Username
	var installPath string

	if username == "root" {
		fmt.Println("Installing VMN as root, VMN will be installed in /usr/local/bin")
		installPath = "/usr/local/bin"
	} else {
		fmt.Println("Installing VMN as " + username + ", VMN will be installed in ~/.local/bin (make sure this is in your PATH)")
		installPath = filepath.Join(utils.GetHome(), ".local", "bin")
	}

	srcFile, err := os.Open(filepath.Join(dir, "vmn"))
	if err != nil {
		panic(err)
	}
	defer srcFile.Close()

	if _, err := os.Stat(installPath); os.IsNotExist(err) {
		if err := os.MkdirAll(installPath, 0755); err != nil {
			panic(err)
		}
	}

	destFile, err := os.Create(filepath.Join(installPath, "vmn"))
	if err != nil {
		panic(err)
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, srcFile); err != nil {
		panic(err)
	}
	fmt.Println("VMN installed successfully!")
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
