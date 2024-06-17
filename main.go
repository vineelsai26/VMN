package main

import (
	"flag"
	"fmt"

	"vineelsai.com/vmn/src"
	"vineelsai.com/vmn/src/node"
	"vineelsai.com/vmn/src/setup"
	"vineelsai.com/vmn/src/shell"
)

func help() {
	fmt.Println(
		`Usage: vmn <command> [version]
Commands:
	install [version]        				Install a specific version of node or latest or lts version of node (default: lts)
	use [version]            				Use a specific version of node
	list [type]              				List all versions, installed versions or lts versions
	uninstall [version]      				Uninstall a specific version of node
	help                     				Print this help
Examples:
	vmn install latest       				Install the latest version of node
	vmn use latest           				Use the latest version of node
	vmn install lts          				Install the latest lts version of node
	vmn use lts              				Use the latest lts version of node
	vmn install 18.15.0      				Install a specific version of node
	vmn use 18.15.0          				Use a specific version of node
	vmn install 20           				Install a specific version of node
	vmn use 20               				Use a specific version of node
	vmn list all             				List all versions of node
	vmn list installed       				List installed versions of node
	vmn list lts             				List lts versions of node
	vmn uninstall all        				Uninstall all versions of node
	vmn uninstall lts        				Uninstall all lts versions of node
	vmn env                  				Print environment variables
	vmn help                 				Print this help
	vmn version              				Print version
	VMN_USE_ROSETTA=true vmn install 14 	Install x86 version of node on Apple Silicon`)
}

func handleNodeVersionManagement(args []string) {
	if args[0] == "help" {
		help()
	} else if args[0] == "install" {
		msg, err := node.Install(args[1])
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	} else if args[0] == "use" {
		msg, err := node.Use(args[1])
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	} else if args[0] == "list" {
		node.List(args[1])
	} else if args[0] == "uninstall" {
		msg, err := node.Uninstall(args[1])
		if err != nil {
			panic(err)
		}
		fmt.Println(msg)
	} else if args[0] == "env" {
		shell.RunShellSpecificCommands(args)
	} else {
		panic("Invalid command")
	}
}

func main() {
	// Parse flags
	flag.Parse()

	// Get clean arguments after flags are parsed
	args := flag.Args()

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
		} else {
			handleNodeVersionManagement(args)
		}
	} else if len(args) >= 3 {
		handleNodeVersionManagement(args)
	} else {
		panic("Too many arguments")
	}
}
