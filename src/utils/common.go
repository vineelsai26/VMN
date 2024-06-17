package utils

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

func GetHome() string {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}
	return home
}

func GetCPUCount() int {
	return runtime.NumCPU()
}

func GetDestination(version string, pl string) string {
	version = strings.TrimPrefix(version, "v")
	return filepath.Join(GetHome(), ".vmn", pl, "v"+version)
}

func GetVersionPath(version string, pl string) (string, error) {
	if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination(version, pl), "bin"), nil
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination(version, pl), "bin"), nil
	}
	return "", fmt.Errorf("unsupported os")
}

func GetBinaryPath(version string, pl string) (string, error) {
	binaryName := pl
	version = strings.TrimPrefix(version, "v")

	if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination("v"+version, pl), "bin", binaryName), nil
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination("v"+version, pl), "bin", binaryName), nil
	}
	return "", fmt.Errorf("unsupported os")
}

func IsInstalled(version string, pl string) bool {
	binaryPath, err := GetBinaryPath(version, pl)
	if err != nil {
		return false
	}

	if _, err := os.Stat(binaryPath); err == nil {
		return true
	}
	return false
}
