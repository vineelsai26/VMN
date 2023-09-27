package main

import (
	"fmt"
	"os"
	"runtime"

	"vineelsai.com/vmn/src"
	"vineelsai.com/vmn/src/node"
	"vineelsai.com/vmn/src/python"
	"vineelsai.com/vmn/src/setup"
	"vineelsai.com/vmn/src/shell"
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

func handleNodeVersionManagement(args []string) {
	if args[0] == "install" {
		if args[1] == "latest" {
			node.InstallLatest()
		} else if args[1] == "lts" {
			node.InstallLatestLTS()
		} else {
			node.InstallSpecific(args[1])
		}
	} else if args[0] == "use" {
		if args[1] == "latest" {
			node.UseLatest()
		} else if args[1] == "lts" {
			node.UseLatestLTS()
		} else {
			node.UseSpecific(args[1])
		}
	} else if args[0] == "list" {
		if args[1] == "all" {
			for _, version := range node.GetAllVersions() {
				println(version)
			}
		} else if args[1] == "lts" {
			for _, version := range node.GetAllLTSVersions() {
				println(version)
			}
		} else if args[1] == "installed" {
			for _, version := range node.GetInstalledVersions() {
				println(version)
			}
		} else {
			panic("Invalid list type")
		}
	} else if args[0] == "uninstall" {
		if args[1] == "all" {
			node.UninstallAll()
		} else if args[1] == "lts" {
			node.UninstallAllLTS()
		} else if args[1] == "latest" {
			node.UninstallLatest()
		} else if args[1] != "" {
			node.UninstallSpecific(args[1])
		} else {
			panic("Invalid version")
		}
	} else if args[0] == "env" {
		shell.RunShellSpecificCommands(args)
	} else {
		panic("Invalid command")
	}
}

func handlePythonVersionManagement(args []string) {
	if runtime.GOOS == "windows" {
		panic("Python version management is not supported on Windows")
	}

	if args[0] == "install" {
		if args[1] == "latest" {
			python.InstallLatest()
		} else {
			python.InstallSpecific(args[1])
		}
	} else if args[0] == "list" {
		if args[1] == "all" {
			for _, version := range python.GetAllVersions() {
				println(version)
			}
		} else if args[1] == "installed" {
			for _, version := range python.GetInstalledVersions() {
				println(version)
			}
		} else {
			panic("Invalid list type")
		}
	} else if args[0] == "uninstall" {
		if args[1] == "all" {
			python.UninstallAll()
		} else if args[1] == "latest" {
			python.UninstallLatest()
		} else if args[1] != "" {
			python.UninstallSpecific(args[1])
		} else {
			panic("Invalid version")
		}
	} else if args[0] == "env" {
		shell.RunShellSpecificCommands(args)
	} else {
		panic("Invalid command")
	}
}

func main() {
	if len(os.Args) == 1 {
		help()
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
			src.GetVersion()
		} else if os.Args[1] == "setup" {
			setup.Install()
		}
	} else if len(os.Args) == 3 {
		handleNodeVersionManagement(os.Args[1:])
	} else if len(os.Args) == 4 {
		if os.Args[1] == "python" {
			handlePythonVersionManagement(os.Args[2:])
		} else if os.Args[1] == "node" {
			handleNodeVersionManagement(os.Args[2:])
		} else {
			panic("Invalid command")
		}
	} else {
		panic("Too many arguments")
	}
}
