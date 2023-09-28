package python

import (
	"fmt"
	"os"

	"vineelsai.com/vmn/src/utils"
)

func uninstallVersion(version string) {
	if !utils.IsInstalled(version, "python") {
		panic("Python version " + version + " is not installed")
	}
	fmt.Printf("Uninstalling Python %s\n", version)
	path, err := utils.GetVersionPath(version, "python")
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
	} else if version == "latest" {
		uninstallVersion(GetLatestVersion())
	} else if version != "" {
		uninstallVersion(version)
	} else {
		panic("Invalid version")
	}
}
