package node

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"vineelsai.com/vmn/src/utils"
)

func getVersions() []map[string]interface{} {
	res, err := http.Get("https://nodejs.org/dist/index.json")
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
	for _, version := range versions {
		allVersions = append(allVersions, version["version"].(string))
	}
	return allVersions
}

func GetAllLTSVersions() []string {
	versions := getVersions()
	var ltsVersions []string
	for _, version := range versions {
		if version["lts"] != nil && version["lts"] != false && version["lts"] != "" {
			ltsVersions = append(ltsVersions, version["version"].(string))
		}
	}
	return ltsVersions
}

func GetInstalledVersions() []string {
	var installedVersions []string
	for _, version := range GetAllVersions() {
		if utils.IsInstalled(version, "node") {
			installedVersions = append(installedVersions, version)
		}
	}
	return installedVersions
}

func GetLatestVersion() string {
	versions := getVersions()
	return versions[0]["version"].(string)
}

func GetLatestLTSVersion() string {
	versions := getVersions()
	for _, version := range versions {
		if version["lts"] != nil && version["lts"] != false && version["lts"] != "" {
			return version["version"].(string)
		}
	}
	return ""
}

func GetLatestVersionOfVersion(major string, minor string) string {
	versions := getVersions()
	if minor != "" {
		for _, version := range versions {
			if strings.Split(version["version"].(string), ".")[0] == "v"+major && strings.Split(version["version"].(string), ".")[1] == minor {
				return version["version"].(string)
			}
		}
	} else {
		for _, version := range versions {
			if strings.Split(version["version"].(string), ".")[0] == "v"+major {
				return version["version"].(string)
			}
		}
	}
	panic("version not found")
}

func GetLatestInstalledVersionOfVersion(major string, minor string) string {
	versions := getVersions()
	if minor != "" {
		for _, version := range versions {
			if strings.Split(version["version"].(string), ".")[0] == "v"+major && strings.Split(version["version"].(string), ".")[1] == minor {
				if utils.IsInstalled(version["version"].(string), "node") {
					return version["version"].(string)
				}
			}
		}
	} else {
		for _, version := range versions {
			if strings.Split(version["version"].(string), ".")[0] == "v"+major {
				if utils.IsInstalled(version["version"].(string), "node") {
					return version["version"].(string)
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
	} else if status == "lts" {
		for _, version := range GetAllLTSVersions() {
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
