//go:build linux || darwin
// +build linux darwin

package location

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"vineelsai.com/vmn/utils"
)

func Set(path string) {
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
