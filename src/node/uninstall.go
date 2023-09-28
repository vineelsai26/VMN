package node

import (
	"fmt"
	"os"

	"vineelsai.com/vmn/src/utils"
)

func uninstallVersion(version string) {
	if !utils.IsInstalled(version, "node") {
		panic("Node.js version " + version + " is not installed")
	}
	fmt.Printf("Uninstalling Node.js %s\n", version)
	path, err := utils.GetVersionPath(version, "node")
	if err != nil {
		panic(err)
	}
	os.RemoveAll(path)
}

func Uninstall(version string) {
	if version == "all" {
		for _, version := range GetAllVersions() {
			uninstallVersion(version)
		}
	} else if version == "lts" {
		for _, version := range GetAllLTSVersions() {
			uninstallVersion(version)
		}
	} else if version == "latest" {
		uninstallVersion(GetLatestVersion())
	} else if version != "" {
		uninstallVersion(version)
	} else {
		panic("Invalid version")
	}
}
