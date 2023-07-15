package node

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"vineelsai.com/vmn/utils"
)

func GetDownloadURL(version string) (string, error) {
	if runtime.GOOS == "windows" {
		if runtime.GOARCH == "amd64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-win-x64.zip", nil
		} else if runtime.GOARCH == "386" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-win-x86.zip", nil
		} else {
			return "", fmt.Errorf("unsupported os or architecture")
		}
	} else if runtime.GOOS == "linux" {
		if runtime.GOARCH == "amd64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-linux-x64.tar.gz", nil
		} else if runtime.GOARCH == "386" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-linux-x86.tar.gz", nil
		} else if runtime.GOARCH == "arm64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-linux-arm64.tar.gz", nil
		} else {
			return "", fmt.Errorf("unsupported os or architecture")
		}
	} else if runtime.GOOS == "darwin" {
		if os.Getenv("VMN_USE_ROSETTA") == "true" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-darwin-x64.tar.gz", nil
		} else if runtime.GOARCH == "amd64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-darwin-x64.tar.gz", nil
		} else if runtime.GOARCH == "arm64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-darwin-arm64.tar.gz", nil
		} else {
			return "", fmt.Errorf("unsupported os or architecture")
		}
	}
	return "", fmt.Errorf("unsupported os or architecture")
}

func Install(version string) {
	fullURLFile, err := GetDownloadURL(version)
	if err != nil {
		panic(err)
	}

	// Download file
	fmt.Println("Downloading Node.js from " + fullURLFile)

	fileName, err := utils.Download(fullURLFile)
	if err != nil {
		panic(err)
	}

	// Unzip file
	fmt.Println("Installing Node.js version " + version + "...")
	if strings.HasSuffix(fileName, ".zip") {
		if err := utils.Unzip(fileName, utils.GetDestination(version)); err != nil {
			panic(err)
		}
	} else if strings.HasSuffix(fileName, ".tar.gz") {
		if err := utils.Untar(fileName, utils.GetDestination(version)); err != nil {
			panic(err)
		}
	}

	// Delete file
	fmt.Println("Cleaning up...")
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
}

func InstallLatest() {
	// Install latest version
	Install(GetLatestVersion())
}

func InstallLatestLTS() {
	// Install latest LTS version
	Install(GetLatestLTSVersion())
}

func InstallSpecific(version string) {
	// Determine Specific SemVer Version from input
	if len(strings.Split(version, ".")) == 3 {
		if strings.Contains(version, "v") {
			Install(version)
		} else {
			Install("v" + version)
		}
	} else if len(strings.Split(version, ".")) == 2 {
		if strings.Contains(version, "v") {
			version = strings.Split(version, "v")[1]
			Install(GetLatestVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1]))
		} else {
			Install(GetLatestVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1]))
		}
	} else if len(strings.Split(version, ".")) == 1 {
		if strings.Contains(version, "v") {
			Install(GetLatestVersionOfVersion(strings.Split(version, "v")[1], ""))
		} else {
			Install(GetLatestVersionOfVersion(version, ""))
		}
	} else {
		panic("invalid version")
	}
}
