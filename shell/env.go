package shell

import (
	"fmt"
	"runtime"
)

func PrintEnv() {
	if runtime.GOOS == "linux" {
		fmt.Println("eval \"SHELL=`ps -p $$ -o comm=`; `vmn env $SHELL`\"")
	} else if runtime.GOOS == "darwin" {
		fmt.Println("eval \"SHELL=`ps -p $$ -o comm=`; `vmn env $SHELL`\"")
	} else if runtime.GOOS == "windows" {
		fmt.Println("vmn env powershell | Out-String | Invoke-Expression")
	} else {
		fmt.Println("Not implemented for this OS")
	}
}

func RunShellSpecificCommands(args []string) {
	if args[2] == "zsh" || args[2] == "-zsh" {
		SetEnvForZSH()
	} else if args[2] == "bash" {
		SetEnvForBash()
	} else if args[2] == "powershell" {
		SetEnvForPowershell()
	} else {
		fmt.Println("Not implemented for this shell")
	}
}
