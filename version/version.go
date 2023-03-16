package version

import "fmt"

var Version = "dev"

func GetVersion() {
	fmt.Println(Version)
}
