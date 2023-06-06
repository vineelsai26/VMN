package main

import (
	"fmt"
	"os"
	"strconv"

	"vineelsai.com/vmn/node"
	"vineelsai.com/vmn/setup"
	"vineelsai.com/vmn/shell"
	"vineelsai.com/vmn/version"
)

func help() {
	fmt.Println("Usage: vmn <command> [version]")
	fmt.Println("Commands:")
	fmt.Println("  install [version]        			Install a specific version of node or latest or lts version of node (default: lts)")
	fmt.Println("  use [version]            			Use a specific version of node")
	fmt.Println("  list [type]              			List all versions, installed versions or lts versions")
	fmt.Println("  uninstall [version]      			Uninstall a specific version of node")
	fmt.Println("  env                      			Print environment variables")
	fmt.Println("  help                     			Print this help")
	fmt.Println("  version                  			Print version")
	fmt.Println("Examples:")
	fmt.Println("  vmn install latest       			Install the latest version of node")
	fmt.Println("  vmn use latest           			Use the latest version of node")
	fmt.Println("  vmn install lts          			Install the latest lts version of node")
	fmt.Println("  vmn use lts              			Use the latest lts version of node")
	fmt.Println("  vmn install 18.15.0      			Install a specific version of node")
	fmt.Println("  vmn use 18.15.0          			Use a specific version of node")
	fmt.Println("  vmn install 20           			Install a specific version of node")
	fmt.Println("  vmn use 20               			Use a specific version of node")
	fmt.Println("  vmn list all             			List all versions of node")
	fmt.Println("  vmn list installed       			List installed versions of node")
	fmt.Println("  vmn list lts             			List lts versions of node")
	fmt.Println("  vmn uninstall all        			Uninstall all versions of node")
	fmt.Println("  vmn uninstall lts        			Uninstall all lts versions of node")
	fmt.Println("  vmn env                  			Print environment variables")
	fmt.Println("  vmn help                 			Print this help")
	fmt.Println("  vmn version              			Print version")
	fmt.Println("  VMN_USE_ROSETTA=true vmn install 14  Install x86 version of node on Apple Silicon")
}

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
		} else if os.Args[1] == "help" {
			help()
		} else if os.Args[1] == "version" {
			version.GetVersion()
		} else if os.Args[1] == "setup" {
			setup.Install()
		} else {
			panic("Invalid command")
		}
	} else {
		help()
	}
}
