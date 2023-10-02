package node

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"vineelsai.com/vmn/src/setup"
	"vineelsai.com/vmn/src/utils"
)

func useVersion(version string) {
	path, err := utils.GetVersionPath(version, "node")
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(filepath.Join(utils.GetHome(), ".vmn", "node-version"), os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err := f.Stat(); err == nil {
		f.Truncate(0)
		f.Seek(0, 0)
		f.WriteString(path)
	}

	if _, err := os.Stat(path); err == nil {
		fmt.Println("Setting VMN_VERSION to " + version + " ... ")
		setup.SetPath(path)
	} else {
		fmt.Println("Node.js version " + version + " is not installed")
	}
}

func Use(version string) {
	version = strings.TrimPrefix(version, "v")
	if version == "latest" {
		version = GetLatestVersion()
	} else if len(strings.Split(version, ".")) == 3 {
		version = "v" + version
	} else if len(strings.Split(version, ".")) == 2 {
		version = GetLatestInstalledVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1])
	} else if len(strings.Split(version, ".")) == 1 {
		version = GetLatestInstalledVersionOfVersion(version, "")
	} else {
		panic("invalid version")
	}

	if !utils.IsInstalled(version, "node") {
		fmt.Println("installing python version " + version + "...")
		installVersion(version)
	}
	useVersion(version)
}
