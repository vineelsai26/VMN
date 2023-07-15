package shell

import "fmt"

func SetEnvForPosix() {
	fmt.Println(`
export PATH="$(cat $HOME/.vmn/current):$PATH"

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

cd() {
	builtin cd "$@"
	setNodeVersion
}

setNodeVersion
	`)
}
