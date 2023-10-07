package python

import (
	"fmt"
	"os"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func uninstallVersion(version string) (string, error) {
	if !utils.IsInstalled(version, "python") {
		return "", fmt.Errorf("Python version " + version + " is not installed")
	}
	fmt.Printf("Uninstalling Python %s\n", version)
	path, err := utils.GetVersionPath(version, "python")
	if err != nil {
		return "", err
	}
	os.RemoveAll(path)

	return "Python version " + version + " uninstalled successfully", nil
}

func Uninstall(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	if version == "all" {
		for _, version := range GetAllVersions() {
			if _, err := uninstallVersion(version); err != nil {
				return "", err
			}
		}
		return "All Python versions uninstalled successfully", nil
	} else if version == "latest" {
		version = "v" + GetLatestVersion()
	} else if version != "" {
		version = "v" + version
	} else {
		return "", fmt.Errorf("invalid version")
	}
	return uninstallVersion(version)
}