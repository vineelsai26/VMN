package tests

import (
	"os"
	"path/filepath"
	"testing"

	"vineelsai.com/vmn/src/node"
	"vineelsai.com/vmn/src/utils"
)

func TestNode_20_Install(t *testing.T) {
	msg, err := node.Install("20")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestNode_20_Use(t *testing.T) {
	version, err := node.Use("20")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "node-version"))
	if err != nil {
		t.Fatal(err)
	}

	path, err := utils.GetVersionPath(version, "node")
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != path {
		t.Fatal("invalid version")
	}
}

func TestNode_20_Uninstall(t *testing.T) {
	msg, err := node.Uninstall("20")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}
