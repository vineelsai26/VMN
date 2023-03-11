package main

import (
	"os"
	"strconv"

	"vineelsai.com/vmn/node"
)

func main() {
	if len(os.Args) > 2 {
		if os.Args[1] == "install" {
			if os.Args[2] == "latest" {
				node.InstallLatest()
			} else if os.Args[2] == "lts" {
				node.InstallLatestLTS()
			} else if _, err := strconv.Atoi(os.Args[2]); err == nil {
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
			} else {
				panic("Invalid list type")
			}
		} else {
			panic("Invalid command")
		}
	}
}
