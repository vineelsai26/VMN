package node

import (
	"fmt"
	"os"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func uninstallVersion(version string) (string, error) {
	if !utils.IsInstalled(version, "node") {
		panic("Node.js version " + version + " is not installed")
	}
	fmt.Printf("Uninstalling Node.js %s\n", version)
	path, err := utils.GetVersionPath(version, "node")
	if err != nil {
		return "", err
	}
	os.RemoveAll(path)

	return "Node.js version " + version + " uninstalled successfully", nil
}

func Uninstall(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	if version == "all" {
		for _, version := range GetAllVersions() {
			if _, err := uninstallVersion(version); err != nil {
				return "", err
			}
		}
	} else if version == "lts" {
		for _, version := range GetAllLTSVersions() {
			if _, err := uninstallVersion(version); err != nil {
				return "", err
			}
		}
	} else if version == "latest" {
		version = "v" + GetLatestVersion()
	} else if version != "" {
		version = "v" + version
	} else {
		return "", fmt.Errorf("invalid version")
	}

	return uninstallVersion(version)
}
