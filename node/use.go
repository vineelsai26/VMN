package node

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"golang.org/x/sys/windows/registry"
)

func GetVersionPath(version string) string {
	if runtime.GOOS == "windows" {
		return GetDestination(version)
	} else if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination(version), "bin")
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination(version), "bin")
	}
	return ""
}

func Use(version string) {
	path := GetVersionPath(version)
	if path == "" {
		fmt.Println("Unsupported OS or architecture")
		return
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Println("Setting VMN_VERSION to " + version + " ... ")

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
	} else {
		fmt.Println("Node.js version " + version + " is not installed")
	}
}

func UseLatest() {
	Use(GetLatestVersion())
}

func UseLatestLTS() {
	Use(GetLatestLTSVersion())
}

func UseSpecific(version string) {
	Use(GetLatestVersionOfVersion(version))
}
