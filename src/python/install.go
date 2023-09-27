package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func Install(version string) {
	fullURLFile := "https://www.python.org/ftp/python/" + version + "/Python-" + version + ".tgz"
	downloadDir := filepath.Join(utils.GetHome(), ".cache", "vmn")
	buildDir := filepath.Join(downloadDir, "build")
	downloadedFilePath := filepath.Join(downloadDir, strings.Split(fullURLFile, "/")[len(strings.Split(fullURLFile, "/"))-1])

	// Download file
	fmt.Println("Downloading Python from " + fullURLFile)
	fileName, err := utils.Download(downloadDir, fullURLFile)
	if err != nil {
		panic(err)
	}

	// Unzip file
	fmt.Println("Installing Python version " + version + "...")
	if strings.HasSuffix(fileName, ".tgz") {
		if err := utils.UnGzip(downloadedFilePath, buildDir); err != nil {
			panic(err)
		}

		err := exec.Command(
			"/bin/bash",
			"-c",
			"cd "+buildDir+" && ./configure --prefix="+utils.GetDestination(version, "python")+" --enable-optimizations && make -j"+strconv.Itoa(utils.GetCPUCount())+" && sudo make altinstall",
		).Run()
		if err != nil {
			panic(err)
		}
	}

	// Delete file
	fmt.Println("Cleaning up...")
	if err := os.Remove(downloadedFilePath); err != nil {
		panic(err)
	}
	if err := os.RemoveAll(buildDir); err != nil {
		panic(err)
	}
}

func InstallLatest() {
	// Install latest version
	Install(GetLatestVersion())
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
