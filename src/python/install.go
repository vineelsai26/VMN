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

func installVersion(version string) {
	fullURLFile := "https://www.python.org/ftp/python/" + version + "/Python-" + version + ".tgz"
	downloadDir := filepath.Join(utils.GetHome(), ".cache", "vmn")
	buildDir := filepath.Join(downloadDir, "build", version)
	downloadedFilePath := filepath.Join(downloadDir, strings.Split(fullURLFile, "/")[len(strings.Split(fullURLFile, "/"))-1])

	// Download file
	fmt.Println("Downloading Python from " + fullURLFile)
	fileName, err := utils.Download(downloadDir, fullURLFile)
	if err != nil {
		panic(err)
	}

	// Unzip file
	fmt.Println("Building Python version " + version + " from source...")
	if strings.HasSuffix(fileName, ".tgz") {
		if err := utils.UnGzip(downloadedFilePath, buildDir); err != nil {
			panic(err)
		}

		cmd := exec.Command(
			"/bin/bash",
			"-c",
			"cd "+buildDir+" && ./configure --prefix="+utils.GetDestination("v"+version, "python")+" --enable-optimizations && make -j"+strconv.Itoa(utils.GetCPUCount())+" && sudo make altinstall",
		)
		out, err := cmd.StdoutPipe()
		if err != nil {
			panic(err)
		}

		if err = cmd.Start(); err != nil {
			panic(err)
		}
		for {
			tmp := make([]byte, 1024)
			_, err := out.Read(tmp)
			fmt.Print(string(tmp))
			if err != nil {
				break
			}
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

func Install(version string) {
	version = strings.Replace(version, "v", "", 1)
	if version == "latest" {
		version = GetLatestVersion()
	} else if len(strings.Split(version, ".")) == 3 {
		version = "v" + version
	} else if len(strings.Split(version, ".")) == 2 {
		version = GetLatestVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1])
	} else if len(strings.Split(version, ".")) == 1 {
		version = GetLatestVersionOfVersion(version, "")
	} else {
		panic("invalid version")
	}

	if utils.IsInstalled(version, "python") {
		fmt.Println("Python version " + version + " is already installed")
	} else {
		installVersion(version)
	}
}
