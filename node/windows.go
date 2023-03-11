//go:build !linux && !darwin
// +build !linux,!darwin

package node

import (
	"fmt"
	"os/exec"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func SetPathWindows(path string) {
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
}

func SetPathLinux(path string) {}
