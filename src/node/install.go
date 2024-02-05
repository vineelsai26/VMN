package node

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func getDownloadURL(version string) (string, error) {
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

func installVersion(version string) (string, error) {
	fullURLFile, err := getDownloadURL(version)
	if err != nil {
		return "", err
	}
	downloadDir := filepath.Join(utils.GetHome(), ".cache", "vmn")
	downloadedFilePath := filepath.Join(downloadDir, strings.Split(fullURLFile, "/")[len(strings.Split(fullURLFile, "/"))-1])

	// Download file
	fmt.Println("Downloading Node.js from " + fullURLFile)

	fileName, err := utils.Download(downloadDir, fullURLFile)
	if err != nil {
		return "", err
	}

	// Unzip file
	fmt.Println("Installing Node.js version " + version + "...")
	if strings.HasSuffix(fileName, ".zip") {
		if err := utils.Unzip(downloadedFilePath, utils.GetDestination(version, "node")); err != nil {
			return "", err
		}
	} else if strings.HasSuffix(fileName, ".tar.gz") {
		if err := utils.UnGzip(downloadedFilePath, utils.GetDestination(version, "node")); err != nil {
			return "", err
		}
	}

	// Delete file
	fmt.Println("Cleaning up...")
	if err := os.Remove(downloadedFilePath); err != nil {
		return "", err
	}

	return "Node.js version " + version + " installed", nil
}

func Install(version string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	if version == "latest" {
		version = GetLatestVersion()
	} else if version == "lts" {
		version = GetLatestLTSVersion()
	} else if len(strings.Split(version, ".")) == 3 {
		version = "v" + version
	} else if len(strings.Split(version, ".")) == 2 {
		version = GetLatestVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1])
	} else if len(strings.Split(version, ".")) == 1 {
		version = GetLatestVersionOfVersion(version, "")
	} else {
		panic("invalid version")
	}

	if utils.IsInstalled(version, "node") {
		return "Node version " + version + " is already installed", nil
	} else {
		return installVersion(version)
	}
}
