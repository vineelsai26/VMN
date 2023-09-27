package utils

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

func GetCPUCount() int {
	return runtime.NumCPU()
}

func GetDestination(version string, pl string) string {
	return filepath.Join(GetHome(), ".vmn", pl, version)
}

func GetVersionPath(version string, pl string) (string, error) {
	if runtime.GOOS == "windows" {
		return GetDestination(version, pl), nil
	} else if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination(version, pl), "bin"), nil
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination(version, pl), "bin"), nil
	}
	return "", fmt.Errorf("unsupported os")
}

func GetBinaryPath(version string, pl string) (string, error) {
	if runtime.GOOS == "windows" {
		return filepath.Join(GetDestination(version, pl), pl+".exe"), nil
	} else if runtime.GOOS == "linux" {
		return filepath.Join(GetDestination(version, pl), "bin", pl), nil
	} else if runtime.GOOS == "darwin" {
		return filepath.Join(GetDestination(version, pl), "bin", pl), nil
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
