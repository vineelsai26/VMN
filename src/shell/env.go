package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"vineelsai.com/vmn/src/node"
	"vineelsai.com/vmn/src/python"
	"vineelsai.com/vmn/src/utils"
)

func setEnvForPosixShell() {
	fmt.Println(`
export PATH="$(cat $HOME/.vmn/node-version):$PATH"
export PATH="$(cat $HOME/.vmn/python-version):$PATH"

function vmn {
	$(whereis vmn | cut -d" " -f2) $@
	if [[ "$1" == "node" ]]
	then
		export PATH="$(cat $HOME/.vmn/node-version):$PATH"
	elif [[ "$1" == "python" ]]
	then
		export PATH="$(cat $HOME/.vmn/python-version):$PATH"
	fi
}

function setNodeVersion {
	if [ -f .vmnrc ]
	then
		echo "Found .vmnrc file"
		if [ -f $HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .vmnrc)" | tail -1 | cut -f2 -d"v")/bin/node ]
		then
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .vmnrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn node install $(cat .vmnrc)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .vmnrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	elif [ -f .nvmrc ]
	then
		echo "Found .nvmrc file"
		if [ -f $HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .nvmrc)" | tail -1 | cut -f2 -d"v")/bin/node ]
		then
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .nvmrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn node install $(cat .nvmrc)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .nvmrc)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	elif [ -f .node-version ]
	then
		echo "Found .node-version file"
		if [ -f $HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin/node ]
		then
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn node install $(cat .node-version)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	fi
}

function setPythonVersion {
	if [ -f .python-version ]
	then
		echo "Found .python-version file"
		if [ -d $HOME/.vmn/python/v$(ls "$HOME/.vmn/python" 2> /dev/null | grep "$(cat .python-version)" | tail -1 | cut -f2 -d"v")/bin ]
		then
			export PATH="$HOME/.vmn/python/v$(ls "$HOME/.vmn/python" 2> /dev/null | grep "$(cat .python-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		else
			vmn --compile python install $(cat .python-version)
			export PATH="$HOME/.vmn/python/v$(ls "$HOME/.vmn/python" 2> /dev/null | grep "$(cat .python-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi

		if [[ $(ls "$HOME/.vmn/python" 2> /dev/null | grep "$(cat .python-version)" | tail -1) != "" ]]; then
			echo "Using python version v$(ls "$HOME/.vmn/python" 2> /dev/null | grep "$(cat .python-version)" | tail -1 | cut -f2 -d"v")"
		fi
	fi
}

function cd {
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
	}
}

func RunShellSpecificCommands(args []string) {
	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn")); os.IsNotExist(err) {
		if err := os.MkdirAll(filepath.Join(utils.GetHome(), ".vmn"), 0755); err != nil {
			panic(err)
		}
	}

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn", "node-version")); os.IsNotExist(err) {
		f, err := os.OpenFile(filepath.Join(utils.GetHome(), ".vmn", "node-version"), os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		installedNodeVersions := node.GetInstalledVersions()
		if len(installedNodeVersions) > 0 {
			if _, err := f.Stat(); err == nil {
				f.Truncate(0)
				f.Seek(0, 0)
				f.WriteString(utils.GetDestination(installedNodeVersions[0], "node"))
			}
		}

	}

	if _, err := os.Stat(filepath.Join(utils.GetHome(), ".vmn", "python-version")); os.IsNotExist(err) {
		f, err := os.OpenFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"), os.O_RDWR|os.O_CREATE, 0755)
		if err != nil {
			panic(err)
		}

		defer f.Close()
		installedPythonVersions := python.GetInstalledVersions()
		if len(installedPythonVersions) > 0 {
			if _, err := f.Stat(); err == nil {
				f.Truncate(0)
				f.Seek(0, 0)
				f.WriteString(utils.GetDestination(python.GetInstalledVersions()[0], "python"))
			}
		}

	}

	if args[1] == "powershell" {
		setEnvForPowershell()
	} else {
		setEnvForPosixShell()
	}
}
