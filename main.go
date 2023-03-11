package main

import (
	"os"
	"strconv"

	"vineelsai.com/vmn/node"
	"vineelsai.com/vmn/shell"
)

func main() {
	if len(os.Args) > 2 {
		IsVersionExist := false
		for _, version := range node.GetAllVersions() {
			if os.Args[2] == version {
				IsVersionExist = true
			}
		}
		if os.Args[1] == "install" {
			if os.Args[2] == "latest" {
				node.InstallLatest()
			} else if os.Args[2] == "lts" {
				node.InstallLatestLTS()
			} else if _, err := strconv.Atoi(os.Args[2]); err == nil {
				node.InstallSpecific(os.Args[2])
			} else if IsVersionExist {
				node.InstallSpecific(os.Args[2])
			} else {
				panic("Invalid version")
			}
		} else if os.Args[1] == "use" {
			if os.Args[2] == "latest" {
				node.UseLatest()
			} else if os.Args[2] == "lts" {
				node.UseLatestLTS()
			} else if _, err := strconv.Atoi(os.Args[2]); err == nil {
				node.UseSpecific(os.Args[2])
			} else if IsVersionExist {
				node.UseSpecific(os.Args[2])
			} else {
				panic("Invalid version")
			}
		} else if os.Args[1] == "list" {
			if os.Args[2] == "all" {
				for _, version := range node.GetAllVersions() {
					println(version)
				}
			} else if os.Args[2] == "lts" {
				for _, version := range node.GetAllLTSVersions() {
					println(version)
				}
			} else if os.Args[2] == "installed" {
				for _, version := range node.GetInstalledVersions() {
					println(version)
				}
			} else {
				panic("Invalid list type")
			}
		} else if os.Args[1] == "uninstall" {
			if os.Args[2] == "all" {
				node.UninstallAll()
			} else if os.Args[2] == "lts" {
				node.UninstallAllLTS()
			} else if os.Args[2] == "latest" {
				node.UninstallLatest()
			} else if os.Args[2] != "" {
				node.UninstallSpecific(os.Args[2])
			} else {
				panic("Invalid version")
			}
		} else if os.Args[1] == "env" {
			shell.RunShellSpecificCommands(os.Args)
		} else {
			panic("Invalid command")
		}
	} else if len(os.Args) == 2 {
		if os.Args[1] == "env" {
			shell.PrintEnv()
		} else if os.Args[1] == "list" {
			for _, version := range node.GetInstalledVersions() {
				println(version)
			}
		} else {
			panic("Invalid command")
		}
	} else {
		panic("Invalid command")
	}
}
