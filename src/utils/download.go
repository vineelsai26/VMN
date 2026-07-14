package utils

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
)

func Download(downloadDir string, fullURLFile string) (string, error) {
	fileURL, err := url.Parse(fullURLFile)
	if err != nil {
		return "", err
	}
	if fileURL.Scheme != "https" || fileURL.Hostname() == "" {
		return "", fmt.Errorf("download URL must use HTTPS: %s", fullURLFile)
	}

	if _, err := os.Stat(downloadDir); os.IsNotExist(err) {
		if err := os.MkdirAll(downloadDir, 0755); err != nil {
			panic(err)
		}
	}

	path := fileURL.Path
	segments := strings.Split(path, "/")
	fileName := segments[len(segments)-1]

	client := secureHTTPClient(fileURL.Hostname())

	// Put content on file
	resp, err := client.Get(fullURLFile)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("bad status: %s, Unable to download the file make sure the specified version or architecture is available", resp.Status)
	}

	// Create blank file
	destination := filepath.Join(downloadDir, fileName)
	file, err := os.CreateTemp(downloadDir, fileName+"-*.part")
	if err != nil {
		return "", err
	}
	temporary := file.Name()
	defer os.Remove(temporary)
	if _, err := io.Copy(file, resp.Body); err != nil {
		_ = file.Close()
		return "", err
	}
	if err := file.Sync(); err != nil {
		_ = file.Close()
		return "", err
	}
	if err := file.Close(); err != nil {
		return "", err
	}
	if err := os.Rename(temporary, destination); err != nil {
		return "", err
	}

	return fileName, nil
}

func FetchBytes(fullURL string, maxBytes int64) ([]byte, error) {
	parsed, err := url.Parse(fullURL)
	if err != nil {
		return nil, err
	}
	if parsed.Scheme != "https" || parsed.Hostname() == "" {
		return nil, fmt.Errorf("download URL must use HTTPS: %s", fullURL)
	}
	client := secureHTTPClient(parsed.Hostname())
	response, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("download %s: bad status %s", fullURL, response.Status)
	}
	var output bytes.Buffer
	if _, err := io.Copy(&output, io.LimitReader(response.Body, maxBytes+1)); err != nil {
		return nil, err
	}
	if int64(output.Len()) > maxBytes {
		return nil, fmt.Errorf("download %s exceeded %d bytes", fullURL, maxBytes)
	}
	return output.Bytes(), nil
}

func secureHTTPClient(expectedHost string) http.Client {
	return http.Client{CheckRedirect: func(request *http.Request, via []*http.Request) error {
		if len(via) >= 5 {
			return fmt.Errorf("too many redirects")
		}
		if request.URL.Scheme != "https" || request.URL.Hostname() != expectedHost {
			return fmt.Errorf("refusing cross-origin or non-HTTPS redirect to %s", request.URL)
		}
		return nil
	}}
}
