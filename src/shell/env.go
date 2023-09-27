package shell

import (
	"fmt"
	"runtime"
)

func setEnvForPosixShell() {
	fmt.Println(`
export PATH="$(cat $HOME/.vmn/node-version):$PATH"

setNodeVersion() {
	if [ -f .vmnrc ]; then
		echo "Found .vmnrc file"
		if [ -f $HOME/.vmn/node/$(cat .vmnrc)/bin/node ]; then
			export PATH="$HOME/.vmn/node/$(cat .vmnrc)/bin:$PATH"
		else
			vmn install $(cat .vmnrc)
			export PATH="$HOME/.vmn/node/$(cat .vmnrc)/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	elif [ -f .nvmrc ]; then
		echo "Found .nvmrc file"
		if [ -f $HOME/.vmn/node/$(cat .nvmrc)/bin/node ]; then
			export PATH="$HOME/.vmn/node/$(cat .nvmrc)/bin:$PATH"
		else
			vmn install $(cat .nvmrc)
			export PATH="$HOME/.vmn/node/$(cat .nvmrc)/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	elif [ -f .node-version ]; then
		echo "Found .node-version file"
		if [ -f $HOME/.vmn/node/$(cat .node-version)/bin/node ]; then
			export PATH="$HOME/.vmn/node/$(cat .node-version)/bin:$PATH"
		else
			vmn install $(cat .node-version)
			export PATH="$HOME/.vmn/node/$(cat .node-version)/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	fi
}

setPythonVersion() {
	if [ -f .python-version ]; then
		echo "Found .python-version file"
		if [ -f $HOME/.vmn/python/$(cat .python-version)/bin/python ]; then
			export PATH="$HOME/.vmn/python/$(cat .python-version)/bin:$PATH"
		else
			vmn python install $(cat .python-version)
			export PATH="$HOME/.vmn/python/$(cat .python-version)/bin:$PATH"
		fi
		echo "Using python version $(python --version)"
	fi
}

cd() {
	builtin cd "$@"
	setNodeVersion
	setPythonVersion
}

setNodeVersion
setPythonVersion
	`)
}

func setEnvForPowershell() {
	fmt.Println(`
function setNodeVersion {
	if (Test-Path .vmnrc) {
        Write-Output "Found .vmnrc file"
        Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .vmnrc);$Env:PATH"
        Write-Output "Using node version $(node --version)"
    }

	if (Test-Path .nvmrc) {
		Write-Output "Found .nvmrc file"
		Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .nvmrc);$Env:PATH"
		Write-Output "Using node version $(node --version)"
	}

	if (Test-Path .node-version) {
		Write-Output "Found .node-version file"
		Set-Item -Path Env:PATH -Value "C:\Users\Vineel\.vmn\node\$(Get-Content .node-version);$Env:PATH"
		Write-Output "Using node version $(node --version)"
	}
}

function changedir($argList) {
    Set-Location $argList
	setNodeVersion
}

setNodeVersion

Set-Alias -Name cd -Option AllScope -Value changedir
	`)
}

func PrintEnv() {
	if runtime.GOOS == "linux" {
		fmt.Println("eval \"SHELL=`ps -p $$ -o comm=`; `vmn env $SHELL`\"")
	} else if runtime.GOOS == "darwin" {
		fmt.Println("eval \"`vmn env zsh`\"")
	} else if runtime.GOOS == "windows" {
		fmt.Println("vmn env powershell | Out-String | Invoke-Expression")
	} else {
		fmt.Println("Not implemented for this OS")
	}
}

func RunShellSpecificCommands(args []string) {
	if args[2] == "zsh" || args[2] == "bash" {
		setEnvForPosixShell()
	} else if args[2] == "powershell" {
		setEnvForPowershell()
	} else {
		fmt.Println("Not implemented for this shell")
	}
}
