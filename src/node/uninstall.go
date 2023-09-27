package node

import (
	"fmt"
	"os"

	"vineelsai.com/vmn/src/utils"
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
	if utils.IsInstalled(version, "node") {
		Uninstall(version)
	}
}

func Uninstall(version string) {
	fmt.Printf("Uninstalling Node.js %s\n", version)
	path, err := utils.GetVersionPath(version, "node")
	if err != nil {
		panic(err)
	}
	os.RemoveAll(path)
}
