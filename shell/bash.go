package shell

import "fmt"

func SetEnvForBash() {
	fmt.Println(`
cd() {
	builtin cd "$@"
	if [ -f .vmnrc ]; then
		echo "Found .vmnrc file"
		export PATH="$HOME/.vmn/node/$(cat .vmnrc)/bin:$PATH"
		echo "Using node version $(node --version)"
	fi

	if [ -f .nvmrc ]; then
		echo "Found .vmnrc file"
		export PATH="$HOME/.vmn/node/$(cat .nvmrc)/bin:$PATH"
		echo "Using node version $(node --version)"
	fi

	if [ -f .node-version ]; then
		echo "Found .vmnrc file"
		export PATH="$HOME/.vmn/node/$(cat .node-version)/bin:$PATH"
		echo "Using node version $(node --version)"
	fi
}

if [ -f .vmnrc ]; then
	echo "Found .vmnrc file"
	export PATH="$HOME/.vmn/node/$(cat .vmnrc)/bin:$PATH"
	echo "Using node version $(node --version)"
fi

if [ -f .nvmrc ]; then
	echo "Found .vmnrc file"
	export PATH="$HOME/.vmn/node/$(cat .nvmrc)/bin:$PATH"
	echo "Using node version $(node --version)"
fi

if [ -f .node-version ]; then
	echo "Found .vmnrc file"
	export PATH="$HOME/.vmn/node/$(cat .node-version)/bin:$PATH"
	echo "Using node version $(node --version)"
fi
		`)
}
