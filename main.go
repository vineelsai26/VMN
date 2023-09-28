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
}

func nodeHelp() {
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
	if args[0] == "help" {
		nodeHelp()
	} else if args[0] == "install" {
		node.Install(args[1])
	} else if args[0] == "use" {
		node.Use(args[1])
	} else if args[0] == "list" {
		node.List(args[1])
	} else if args[0] == "uninstall" {
		node.Uninstall(args[1])
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

	if args[0] == "help" {
		help()
	} else if args[0] == "install" {
		python.Install(args[1])
	} else if args[0] == "use" {
		python.Use(args[1])
	} else if args[0] == "list" {
		python.List(args[1])
	} else if args[0] == "uninstall" {
		python.Uninstall(args[1])
	} else if args[0] == "env" {
		shell.RunShellSpecificCommands(args)
	} else {
		panic("Invalid command")
	}
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		help()
	} else if len(args) == 1 {
		if args[0] == "env" {
			shell.PrintEnv()
		} else if args[0] == "list" {
			for _, version := range node.GetInstalledVersions() {
				println(version)
			}
		} else if args[0] == "help" {
			help()
		} else if args[0] == "version" {
			src.GetVersion()
		} else if args[0] == "setup" {
			setup.Install()
		}
	} else if len(args) == 2 {
		// python or node
		if args[0] == "python" {
			handlePythonVersionManagement(args)
		} else if args[0] == "node" {
			handleNodeVersionManagement(args)
		} else {
			handleNodeVersionManagement(args)
		}
	} else if len(args) == 3 {
		// python or node
		if args[0] == "python" {
			handlePythonVersionManagement(args[1:])
		} else if args[0] == "node" {
			handleNodeVersionManagement(args[1:])
		} else {
			panic("Invalid command")
		}
	} else {
		panic("Too many arguments")
	}
}
