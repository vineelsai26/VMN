package node

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
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

func GetLatestVersionOfVersion(major string) string {
	versions := getVersions()
	for _, version := range versions {
		if strings.Split(version["version"].(string), ".")[0] == "v"+major {
			return version["version"].(string)
		}
	}
	return ""
}
