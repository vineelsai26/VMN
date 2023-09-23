package node

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"vineelsai.com/vmn/src/setup"
	"vineelsai.com/vmn/src/utils"
)

func SetPath(path string) {
	if runtime.GOOS == "windows" || runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		setup.SetPath(path)
	} else {
		fmt.Println("Not implemented for this OS")
	}
}

func Use(version string) {
	path, err := utils.GetVersionPath(version)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(filepath.Join(utils.GetHome(), ".vmn", "current"), os.O_RDWR|os.O_CREATE, 0755)
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
		SetPath(path)
	} else {
		fmt.Println("Node.js version " + version + " is not installed")
	}
}

func UseLatest() {
	Use(GetLatestVersion())
}

func UseLatestLTS() {
	Use(GetLatestLTSVersion())
}

func UseSpecific(version string) {
	if len(strings.Split(version, ".")) == 3 {
		if strings.Contains(version, "v") {
			Use(version)
		} else {
			Use("v" + version)
		}
	} else if len(strings.Split(version, ".")) == 2 {
		if strings.Contains(version, "v") {
			version = strings.Split(version, "v")[1]
			Use(GetLatestInstalledVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1]))
		} else {
			Use(GetLatestInstalledVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1]))
		}
	} else if len(strings.Split(version, ".")) == 1 {
		if strings.Contains(version, "v") {
			Use(GetLatestInstalledVersionOfVersion(strings.Split(version, "v")[1], ""))
		} else {
			Use(GetLatestInstalledVersionOfVersion(version, ""))
		}
	} else {
		panic("invalid version")
	}
}
