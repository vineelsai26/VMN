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

func GetLatestVersionOfVersion(major string) string {
	versions := getVersions()
	for _, version := range versions {
		if strings.Split(version["version"].(string), ".")[0] == "v"+major {
			return version["version"].(string)
		}
	}
	return ""
}
