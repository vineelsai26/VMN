package utils

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
)

func Download(fullURLFile string) (string, error) {
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return "", err
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

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

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s, Unable to download the file make sure the specified version or architecture is available", resp.Status)
	}

	defer resp.Body.Close()

	// Create blank file
	file, err := os.Create(fileName)
	if err != nil {
		return "", err
	}

	io.Copy(file, resp.Body)

	defer file.Close()

	return fileName, nil
}
