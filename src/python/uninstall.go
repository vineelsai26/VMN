package python

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

func UninstallSpecific(version string) {
	if utils.IsInstalled(version, "python") {
		Uninstall(version)
	}
}

func Uninstall(version string) {
	fmt.Printf("Uninstalling Python %s\n", version)
	path, err := utils.GetVersionPath(version, "python")
	if err != nil {
		panic(err)
	}
	os.RemoveAll(path)
}
