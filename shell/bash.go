package shell

import "fmt"

func SetEnvForBash() {
	fmt.Println(`
cd() {
	builtin cd "$@"
	if [ -f .vmnrc ]; then
		export PATH="$PATH:$HOME/.vmn/node/$(cat .vmnrc)"
	fi
}

if [ -f .vmnrc ]; then
	export PATH="$PATH:$HOME/.vmn/node/$(cat .vmnrc)"
fi
	`)
}
