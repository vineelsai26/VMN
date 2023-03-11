package node

import (
	"fmt"
	"os"
)

func UninstallAll() {
	for _, version := range GetAllVersions() {
		UninstallSpecific(version)
	}
}

func UninstallLatest() {
	UninstallSpecific(GetLatestVersion())
}

func UninstallAllLTS() {
	for _, version := range GetAllLTSVersions() {
		UninstallSpecific(version)
	}
}

func UninstallSpecific(version string) {
	if IsInstalled(version) {
		Uninstall(version)
	}
}

func Uninstall(version string) {
	fmt.Printf("Uninstalling Node.js %s\n", version)
	path, err := GetVersionPath(version)
	if err != nil {
		panic(err)
	}
	os.RemoveAll(path)
}
