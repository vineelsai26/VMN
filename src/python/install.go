package python

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func installPythonFromSource(version string, compile_flags_override string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	fullURLFile := "https://www.python.org/ftp/python/" + version + "/Python-" + version + ".tgz"
	downloadDir := filepath.Join(utils.GetHome(), ".cache", "vmn")
	buildDir := filepath.Join(downloadDir, "build", version)
	downloadedFilePath := filepath.Join(downloadDir, strings.Split(fullURLFile, "/")[len(strings.Split(fullURLFile, "/"))-1])

	// Check if make command is available
	make_cmd := exec.Command(
		"make",
	)
	_, err := make_cmd.StdoutPipe()
	if err != nil {
		return "", err
	}

	if err = make_cmd.Start(); err != nil {
		return "", fmt.Errorf("make command not found. Please install 'make' and try again")
	}

	// Download file
	fmt.Println("Downloading Python from " + fullURLFile)
	fileName, err := utils.Download(downloadDir, fullURLFile)
	if err != nil {
		return "", err
	}

	// Unzip file
	fmt.Println("Building Python version " + version + " from source...")
	if strings.HasSuffix(fileName, ".tgz") {
		if err := utils.UnGzip(downloadedFilePath, buildDir, false); err != nil {
			return "", err
		}

		build_flags := "--enable-optimizations"

		if compile_flags_override != "" {
			build_flags = compile_flags_override
		}
		build_install_command := "./configure --prefix=" + utils.GetDestination("v"+version, "python") + " " + build_flags + " && make -j" + strconv.Itoa(utils.GetCPUCount()) + " && make altinstall"

		fmt.Println(build_install_command)
		cmd := exec.Command(
			"bash",
			"-c",
			"cd "+buildDir+" && "+build_install_command,
		)
		out, err := cmd.StdoutPipe()
		if err != nil {
			return "", err
		}

		if err = cmd.Start(); err != nil {
			return "", err
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

	// symlink python3.x to python and python3
	fmt.Println("Symlinking python3.x to python and python3...")
	pythonBinaryPath, err := utils.GetBinaryPath(version, "python")
	if err != nil {
		return "", err
	}
	pythonSymlinkPath, err := utils.GetVersionPath("v"+version, "python")
	if err != nil {
		return "", err
	}
	pythonSymlinks := []string{"python", "python3"}

	for _, symlink := range pythonSymlinks {
		if err := os.Symlink(pythonBinaryPath, filepath.Join(pythonSymlinkPath, symlink)); err != nil {
			return "", err
		}
	}

	// symlink pip3.x to pip and pip3
	pipBinaryPath, err := utils.GetBinaryPath(version, "pip")
	if err != nil {
		return "", err
	}
	pipSymlinks := []string{"pip", "pip3"}

	for _, symlink := range pipSymlinks {
		if err := os.Symlink(pipBinaryPath, filepath.Join(pythonSymlinkPath, symlink)); err != nil {
			return "", err
		}
	}

	// Delete file
	fmt.Println("Cleaning up...")
	if err := os.Remove(downloadedFilePath); err != nil {
		return "", err
	}
	if err := os.RemoveAll(buildDir); err != nil {
		exec.Command("bash", "-c", "rm -rf "+buildDir).Run()
	}

	return "Python version " + version + " installed successfully", nil
}

func installPython(version string) (string, error) {
	fullURLFile := "https://repo.vineelsai.com/" + runtime.GOOS + "/generic/" + runtime.GOARCH + "/python-" + strings.TrimPrefix(version, "v") + ".tar.gz"
	downloadDir := filepath.Join(utils.GetHome(), ".cache", "vmn")
	downloadedFilePath := filepath.Join(downloadDir, strings.Split(fullURLFile, "/")[len(strings.Split(fullURLFile, "/"))-1])

	// Download file
	fmt.Println("Downloading Python from " + fullURLFile)
	fileName, err := utils.Download(downloadDir, fullURLFile)
	if err != nil {
		return "", fmt.Errorf("Python version " + version + " not found in precompiled package repo, please run the following command to compile from source \n ```\n vmn --compile python install " + version + "\n ``` \n NOTE: compiling from source might take a while depending on your system resources.")
	}

	// Unzip file
	fmt.Println("Installing Python version " + version + "...")
	if strings.HasSuffix(fileName, ".zip") {
		if err := utils.Unzip(downloadedFilePath, utils.GetDestination(version, "python")); err != nil {
			return "", err
		}
	} else if strings.HasSuffix(fileName, ".tar.gz") {
		if err := utils.UnGzip(downloadedFilePath, utils.GetDestination(version, "python"), true); err != nil {
			return "", err
		}
	}

	// Delete file
	fmt.Println("Cleaning up...")
	if err := os.Remove(downloadedFilePath); err != nil {
		return "", err
	}

	return "Python version " + version + " installed", nil
}

func installVersion(version string, compile bool, compile_flags_override string) (string, error) {
	if compile {
		return installPythonFromSource(version, compile_flags_override)
	} else {
		return installPython(version)
	}
}

func Install(version string, compile bool, compile_flags_override string) (string, error) {
	version = strings.TrimPrefix(version, "v")
	if version == "latest" {
		version = GetLatestVersion()
	} else if len(strings.Split(version, ".")) == 3 {
		version = "v" + version
	} else if len(strings.Split(version, ".")) == 2 {
		version = GetLatestVersionOfVersion(strings.Split(version, ".")[0], strings.Split(version, ".")[1])
	} else if len(strings.Split(version, ".")) == 1 {
		version = GetLatestVersionOfVersion(version, "")
	} else {
		return "", fmt.Errorf("invalid version")
	}

	if utils.IsInstalled(version, "python") {
		return "Python version " + version + " is already installed", nil
	} else {
		return installVersion(version, compile, compile_flags_override)
	}
}
