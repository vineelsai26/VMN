package node

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
)

func SetPath(path string) {
	if runtime.GOOS == "windows" {
		SetPathWindows(path)
	} else if runtime.GOOS == "linux" || runtime.GOOS == "darwin" {
		SetPathLinux(path)
	} else {
		fmt.Println("Not implemented for this OS")
	}
}

func Use(version string) {
	path, err := GetVersionPath(version)
	if err != nil {
		panic(err)
	}

	if _, err := os.Stat(filepath.Join(GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	f, err := os.OpenFile(filepath.Join(GetHome(), ".vmn", "current"), os.O_RDWR|os.O_CREATE, 0755)
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
	Use(GetLatestVersionOfVersion(version))
}
