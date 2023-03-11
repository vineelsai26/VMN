package node

import (
	"archive/zip"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
		if runtime.GOARCH == "amd64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-darwin-x64.tar.gz", nil
		} else if runtime.GOARCH == "arm64" {
			return "https://nodejs.org/dist/" + version + "/node-" + version + "-darwin-arm64.tar.gz", nil
		} else {
			return "", fmt.Errorf("unsupported os or architecture")
		}
	}
	return "", fmt.Errorf("unsupported os or architecture")
}

func Download(version string) (string, error) {
	fullURLFile, err := GetDownloadURL(version)
	if err != nil {
		return "", err
	}

	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return "", err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	client := http.Client{
		CheckRedirect: func(r *http.Request, via []*http.Request) error {
			r.URL.Opaque = r.URL.Path
			return nil
		},
	}

	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	io.Copy(file, resp.Body)

	defer file.Close()

	return fileName, nil
}

func Unzip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer func() {
		if err := r.Close(); err != nil {
			panic(err)
		}
	}()

	os.MkdirAll(dest, 0755)

	// Closure to address file descriptors issue with all the deferred .Close() methods
	extractAndWriteFile := func(f *zip.File) error {
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer func() {
			if err := rc.Close(); err != nil {
				panic(err)
			}
		}()

		path := filepath.Join(dest, strings.Join(strings.Split(f.Name, "/")[1:], "/"))

		// Check for ZipSlip (Directory traversal)
		if !strings.HasPrefix(path, filepath.Clean(dest)+string(os.PathSeparator)) {
			return fmt.Errorf("illegal file path: %s", path)
		}

		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.Mode())
		} else {
			os.MkdirAll(filepath.Dir(path), f.Mode())
			f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				return err
			}
			defer func() {
				if err := f.Close(); err != nil {
					panic(err)
				}
			}()

			_, err = io.Copy(f, rc)
			if err != nil {
				return err
			}
		}
		return nil
	}

	for _, f := range r.File[1:] {
		err := extractAndWriteFile(f)
		if err != nil {
			return err
		}
	}

	return nil
}

func Install(version string) {
	// Download file
	fmt.Println("Downloading Node.js version " + version + "...")
	fileName, err := Download(version)
	if err != nil {
		panic(err)
	}

	// Unzip file
	fmt.Println("Installing Node.js version " + version + "...")
	if err := Unzip(fileName, GetDestination(version)); err != nil {
		panic(err)
	}

	// Delete file
	fmt.Println("Cleaning up...")
	if err := os.Remove(fileName); err != nil {
		panic(err)
	}
}

func InstallLatest() {
	Install(GetLatestVersion())
}

func InstallLatestLTS() {
	Install(GetLatestLTSVersion())
}

func InstallSpecific(version string) {
	Install(GetLatestVersionOfVersion(version))
}
