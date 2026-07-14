package node

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/ProtonMail/go-crypto/openpgp"
	"vineelsai.com/vmn/src/utils"
)

var nodeReleaseKeyFingerprints = []string{
	"5BE8A3F6C8A5C01D106C0AD820B1A390B168D356",
	"DD792F5973C6DE52C432CBDAC77ABFA00DDBF2B7",
	"CC68F5A3106FF448322E48ED27F5E38D5B0A215F",
	"8FCCA13FEF1D0C2E91008E09770F7A9A5AE15600",
	"890C08DB8579162FEE0DF9DB8BEAB4DFCF555EF4",
	"C82FA3AE1CBEDC6BE46B9360C43CEC45C17AB93C",
	"108F52B48DB57BB0CC439B2997B01419BD92F80A",
	"A363A499291CBBC940DD62E41F10027AF002F8B0",
}

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
	if err := verifyNodeRelease(fullURLFile, downloadedFilePath); err != nil {
		return "", err
	}

	// Unzip file
	fmt.Println("Installing Node.js version " + version + "...")
	if strings.HasSuffix(fileName, ".zip") {
		if err := utils.Unzip(downloadedFilePath, utils.GetDestination(version, "node")); err != nil {
			return "", err
		}
	} else if strings.HasSuffix(fileName, ".tar.gz") {
		if err := utils.UnGzip(downloadedFilePath, utils.GetDestination(version, "node"), false); err != nil {
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

func verifyNodeRelease(archiveURL, archivePath string) error {
	baseURL := archiveURL[:strings.LastIndex(archiveURL, "/")+1]
	manifest, err := utils.FetchBytes(baseURL+"SHASUMS256.txt", 1<<20)
	if err != nil {
		return fmt.Errorf("download Node.js checksums: %w", err)
	}
	signature, err := utils.FetchBytes(baseURL+"SHASUMS256.txt.sig", 64<<10)
	if err != nil {
		return fmt.Errorf("download Node.js checksum signature: %w", err)
	}
	keyring, err := nodeReleaseKeyring()
	if err != nil {
		return err
	}
	if _, err := openpgp.CheckDetachedSignature(keyring, bytes.NewReader(manifest), bytes.NewReader(signature), nil); err != nil {
		return fmt.Errorf("verify Node.js signed checksums: %w", err)
	}

	archiveName := filepath.Base(archivePath)
	expected := ""
	for _, line := range strings.Split(string(manifest), "\n") {
		fields := strings.Fields(line)
		if len(fields) == 2 && fields[1] == archiveName {
			expected = fields[0]
			break
		}
	}
	if len(expected) != sha256.Size*2 {
		return fmt.Errorf("Node.js signed checksums do not contain %s", archiveName)
	}
	file, err := os.Open(archivePath)
	if err != nil {
		return err
	}
	defer file.Close()
	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return err
	}
	actual := hex.EncodeToString(hash.Sum(nil))
	if !strings.EqualFold(actual, expected) {
		return fmt.Errorf("Node.js archive checksum mismatch: got %s, want %s", actual, expected)
	}
	return nil
}

func nodeReleaseKeyring() (openpgp.EntityList, error) {
	var keyring openpgp.EntityList
	for _, fingerprint := range nodeReleaseKeyFingerprints {
		armored, err := utils.FetchBytes("https://keys.openpgp.org/vks/v1/by-fingerprint/"+fingerprint, 1<<20)
		if err != nil {
			continue
		}
		entities, err := openpgp.ReadArmoredKeyRing(bytes.NewReader(armored))
		if err != nil {
			continue
		}
		for _, entity := range entities {
			actual := strings.ToUpper(hex.EncodeToString(entity.PrimaryKey.Fingerprint[:]))
			if actual == fingerprint {
				keyring = append(keyring, entity)
			}
		}
	}
	if len(keyring) == 0 {
		return nil, fmt.Errorf("could not load any pinned Node.js release signing keys")
	}
	return keyring, nil
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
