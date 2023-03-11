package node

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func GetHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func GetDestination(version string) string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return filepath.Join(home, ".vmn", "node", version)
}

func GetVersionPath(version string) (string, error) {
	if runtime.GOOS == "windows" {
		return GetDestination(version), nil
	} else if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination(version), "bin"), nil
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination(version), "bin"), nil
	}
	return "", fmt.Errorf("unsupported os")
}

func GetNodePath(version string) (string, error) {
	if runtime.GOOS == "windows" {
		return filepath.Join(GetDestination(version), "node.exe"), nil
	} else if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination(version), "bin", "node"), nil
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination(version), "bin", "node"), nil
	}
	return "", fmt.Errorf("unsupported os")
}

func IsInstalled(version string) bool {
	nodePath, err := GetNodePath(version)
	if err != nil {
		return false
	}

	if _, err := os.Stat(nodePath); err == nil {
		return true
	}
	return false
}
