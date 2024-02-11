package python

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vineelsai.com/vmn/src/setup"
	"vineelsai.com/vmn/src/utils"
)

func useVersion(version string) (string, error) {
	path, err := utils.GetVersionPath("v"+version, "python")
	if err != nil {
		return "", err
	}

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			return "", err
		}
	}

	f, err := os.OpenFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return "", err
	}
	defer f.Close()

	if _, err := f.Stat(); err == nil {
		f.Truncate(0)
		f.Seek(0, 0)
		f.WriteString(path)
	}

	if _, err := os.Stat(path); err == nil {
		setup.SetPath(path)
	} else {
		return "", fmt.Errorf("Python version " + version + " is not installed")
	}

	return version, nil
}

func Use(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	if version == "latest" {
		version = GetLatestVersion()
	} else if len(strings.Split(version, ".")) == 3 {
		version = strings.TrimPrefix(version, "v")
	} else if len(strings.Split(version, ".")) == 2 {
		version = GetLatestVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1])
	} else if len(strings.Split(version, ".")) == 1 {
		version = GetLatestVersionOfVersion(version, "")
	} else {
		panic("invalid version")
	}

	if !utils.IsInstalled(version, "python") {
		fmt.Println("installing python version " + version + "...")
		installVersion(version, false, "")
	}
	return useVersion(version)
}
