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

func SetPath(path string) {
	if runtime.GOOS == "windows" {
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
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		exec.Command("export", "VMN_VERSION="+path).Run()

		home, err := os.UserHomeDir()
		if err != nil {
			panic(err)
		}

		shells := []string{".bashrc", ".zshrc"}

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
	} else {
		fmt.Println("Not implemented for this OS")
	}
}

func Use(version string) {
	path := GetVersionPath(version)
	if path == "" {
		fmt.Println("Unsupported OS or architecture")
		return
	}

	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	f, err := os.OpenFile(filepath.Join(home, ".vmn", "current"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Stat(); err == nil {
		f.Truncate(0)
		f.Seek(0, 0)
		f.WriteString(path)
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Println("Setting VMN_VERSION to " + version + " ... ")
		SetPath(path)
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
