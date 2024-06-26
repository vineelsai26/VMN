package shell

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"vineelsai.com/vmn/src/node"
	"vineelsai.com/vmn/src/utils"
)

func setEnvForPosixShell() {
	fmt.Println(`
export PATH="$(cat $HOME/.vmn/node-version):$PATH"

function vmn {
	$(whereis vmn | cut -d" " -f2) $@
	if [[ "$1" == "node" ]]
	then
		export PATH="$(cat $HOME/.vmn/node-version):$PATH"
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
			vmn install $(cat .vmnrc)
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
			vmn install $(cat .nvmrc)
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
			vmn install $(cat .node-version)
			export PATH="$HOME/.vmn/node/v$(ls "$HOME/.vmn/node" 2> /dev/null | grep "$(cat .node-version)" | tail -1 | cut -f2 -d"v")/bin:$PATH"
		fi
		echo "Using node version $(node --version)"
	fi
}

function cd {
	builtin cd "$@"
	setNodeVersion
}

setNodeVersion
	`)
}

func PrintEnv() {
	if runtime.GOOS == "linux" {
		fmt.Println("eval \"SHELL=`ps -p $$ -o comm=`; `vmn env $SHELL`\"")
	} else if runtime.GOOS == "darwin" {
		fmt.Println("eval \"`vmn env zsh`\"")
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

	setEnvForPosixShell()
}
