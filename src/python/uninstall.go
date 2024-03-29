package python

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func uninstallVersion(version string) (string, error) {
	if !utils.IsInstalled(version, "python") {
		return "", fmt.Errorf("Python version " + version + " is not installed")
	}
	fmt.Printf("Uninstalling Python %s\n", version)
	path := utils.GetDestination(version, "python")
	if err := os.RemoveAll(strings.TrimSpace(path)); err != nil {
		exec.Command("/bin/bash", "-c", "rm -rf "+path).Run()
	}

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
	} else if len(strings.Split(version, ".")) == 3 {
		version = "v" + version
	} else if version != "" {
		for _, ver := range GetInstalledVersions() {
			if strings.HasPrefix(ver, version) {
				version = "v" + ver
				if _, err := uninstallVersion(version); err != nil {
					return "", err
				}
			}
		}
		return "Python version " + version + " uninstalled successfully", nil
	} else {
		return "", fmt.Errorf("invalid version")
	}
	return uninstallVersion(version)
}
