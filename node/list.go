package node

import (
	"vineelsai.com/vmn/utils"
)

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
		if utils.IsInstalled(version) {
			installedVersions = append(installedVersions, version)
		}
	}
	return installedVersions
}
