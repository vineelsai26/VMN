package python

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func getVersions() []map[string]interface{} {
	res, err := http.Get("https://api.github.com/repos/python/cpython/git/matching-refs/tags/v")
	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		print(err)
	}

	var versions []map[string]interface{}
	json.Unmarshal(body, &versions)
	return versions
}

func GetAllVersions() []string {
	versions := getVersions()
	var allVersions []string
	for i := range versions {
		version := versions[len(versions)-i-1]
		allVersions = append(allVersions, strings.Replace(version["ref"].(string), "refs/tags/v", "", 1))
	}

	return allVersions
}

func GetInstalledVersions() []string {
	var installedVersions []string
	for _, version := range GetAllVersions() {
		if utils.IsInstalled(version, "python") {
			installedVersions = append(installedVersions, version)
		}
	}
	return installedVersions
}

func GetLatestVersion() string {
	versions := GetAllVersions()

	for _, version := range versions {
		if !strings.Contains(version, "a") || strings.Contains(version, "b") || strings.Contains(version, "c") {
			return version
		}
	}

	return versions[0]
}

func GetLatestVersionOfVersion(major string, minor string) string {
	versions := GetAllVersions()
	if minor != "" {
		for _, version := range versions {
			if strings.Split(version, ".")[0] == major && strings.Split(version, ".")[1] == minor {
				return version
			}
		}
	} else {
		for _, version := range versions {
			if strings.Split(version, ".")[0] == major {
				return version
			}
		}
	}
	panic("version not found")
}

func GetLatestInstalledVersionOfVersion(major string, minor string) string {
	versions := GetAllVersions()
	if minor != "" {
		for _, version := range versions {
			if strings.Split(version, ".")[0] == major && strings.Split(version, ".")[1] == minor {
				if utils.IsInstalled("v"+version, "python") {
					return version
				}
			}
		}
	} else {
		for _, version := range versions {
			if strings.Split(version, ".")[0] == "v"+major {
				if utils.IsInstalled(version, "python") {
					return version
				}
			}
		}
	}
	panic("version not installed")
}

func List(status string) {
	if status == "all" {
		for _, version := range GetAllVersions() {
			println(version)
		}
	} else if status == "installed" {
		for _, version := range GetInstalledVersions() {
			println(version)
		}
	} else {
		panic("Invalid list type")
	}
}
