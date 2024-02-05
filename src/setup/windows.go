//go:build !linux && !darwin
// +build !linux,!darwin

package setup

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
	"vineelsai.com/vmn/src/utils"
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

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	destFile, err := os.Create(filepath.Join(utils.GetHome(), ".vmn", "vmn.exe"))
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
		if env == filepath.Join(utils.GetHome(), ".vmn") {
			isPathVariableExists = true
			return
		}
	}

	if !isPathVariableExists {
		exec.Command("setx", "PATH", userPath+filepath.Join(utils.GetHome(), ".vmn")).Run()
	}
}

func SetPath(path string) {
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE)
	if err != nil {
		panic(err)
	}
	defer k.Close()

	userPath, _, err := k.GetStringValue("Path")
	if err != nil {
		panic(err)
	}

	out, err := exec.Command("setx", "VMN_VERSION", path).Output()
	if err != nil {
		panic(err)
	}

	fmt.Println(string(out))

	var isPathVariableExists bool = false

	for _, env := range strings.Split(userPath, ";") {
		if env == "%VMN_VERSION%" {
			isPathVariableExists = true
			return
		}
	}

	if !isPathVariableExists {
		out, err = exec.Command("setx", "PATH", "%VMN_VERSION%;"+userPath).Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(string(out))
	}

	shells := []string{
		filepath.Join(utils.GetHome(), "Documents\\WindowsPowerShell\\Microsoft.PowerShell_profile.ps1"),
		filepath.Join(utils.GetHome(), "Documents\\PowerShell\\Microsoft.PowerShell_profile.ps1"),
	}

	for _, shell := range shells {
		f, err := os.OpenFile(shell, os.O_RDONLY, 0755)
		if err != nil {
			continue
		}
		defer f.Close()

		if _, err := f.Stat(); err == nil {
			b := make([]byte, 1024)
			n, err := f.Read(b)
			if err != nil {
				panic(err)
			}

			if !strings.Contains(string(b[:n]), "vmn env | Out-String | Invoke-Expression") {
				file, err := os.OpenFile(shell, os.O_APPEND|os.O_WRONLY, 0755)
				if err != nil {
					panic(err)
				}
				defer file.Close()

				if _, err := file.WriteString("vmn env | Out-String | Invoke-Expression\n"); err != nil {
					panic(err)
				}
			}
		}
	}

}
