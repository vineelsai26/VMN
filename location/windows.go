//go:build !linux && !darwin
// +build !linux,!darwin

package location

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func Set(path string) {
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

	f, err := os.OpenFile("C:\\Users\\Vineel\\Documents\\WindowsPowerShell\\Microsoft.PowerShell_profile.ps1", os.O_RDONLY, 0755)
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

		shells := []string{
			filepath.Join(GetHome(), "Documents\\WindowsPowerShell\\Microsoft.PowerShell_profile.ps1"),
			filepath.Join(GetHome(), "Documents\\PowerShell\\Microsoft.PowerShell_profile.ps1"),
		}

		for _, shell := range shells {
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
