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
	fmt.Println(
		`Usage: vmn <runtime> <command> [version]
use vmn <runtime> help for more information
runtimes currently supported: node, python`)
}

func nodeHelp() {
	fmt.Println(
		`Usage: vmn node <command> [version]
Commands:
	node install [version]        				Install a specific version of node or latest or lts version of node (default: lts)
	node use [version]            				Use a specific version of node
	node list [type]              				List all versions, installed versions or lts versions
	node uninstall [version]      				Uninstall a specific version of node
	node help                     				Print this help
Examples:
	vmn node install latest       				Install the latest version of node
	vmn node use latest           				Use the latest version of node
	vmn node install lts          				Install the latest lts version of node
	vmn node use lts              				Use the latest lts version of node
	vmn node install 18.15.0      				Install a specific version of node
	vmn node use 18.15.0          				Use a specific version of node
	vmn node install 20           				Install a specific version of node
	vmn node use 20               				Use a specific version of node
	vmn node list all             				List all versions of node
	vmn node list installed       				List installed versions of node
	vmn node list lts             				List lts versions of node
	vmn node uninstall all        				Uninstall all versions of node
	vmn node uninstall lts        				Uninstall all lts versions of node
	vmn node env                  				Print environment variables
	vmn node help                 				Print this help
	vmn node version              				Print version
	VMN_USE_ROSETTA=true vmn node install 14 	Install x86 version of node on Apple Silicon`)
}

func pythonHelp() {
	fmt.Println(
		`Usage: vmn python <command> [version]
Commands:
	python install [version]        			Install a specific version of python or latest or lts version of python (default: lts)
	python use [version]            			Use a specific version of python
	python list [type]              			List all versions, installed versions or lts versions
	python uninstall [version]      			Uninstall a specific version of python
	python help                     			Print this help section
Examples:
	vmn python install latest       			Install the latest version of python
	vmn python use latest           			Use the latest version of python
	vmn python install 3.11         			Install a specific version of python
	vmn python use 3.11		          			Use a specific version of python
	vmn python list all             			List all versions of python
	vmn python list installed       			List installed versions of python
	vmn python uninstall all        			Uninstall all versions of python
	vmn python help                 			Print this help`)
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
		pythonHelp()
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
		if args[0] == "env" {
			shell.RunShellSpecificCommands(args)
		} else if args[0] == "python" {
			handlePythonVersionManagement(args)
		} else if args[0] == "node" {
			handleNodeVersionManagement(args[1:])
		} else {
			handleNodeVersionManagement(args[1:])
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
