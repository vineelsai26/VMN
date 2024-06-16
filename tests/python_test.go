//go:build (linux && ignore) || (darwin && ignore) || !windows
// +build linux,ignore darwin,ignore !windows

package tests

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"vineelsai.com/vmn/src/python"
	"vineelsai.com/vmn/src/utils"
)

func TestPython_3_12_Install(t *testing.T) {
	msg, err := python.Install("3.12.3", false, "")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_12_Use(t *testing.T) {
	version, err := python.Use("3.12.3")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"))
	if err != nil {
		t.Fatal(err)
	}

	path, err := utils.GetVersionPath("v"+version, "python")
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != path {
		t.Fatal("invalid version")
	}
}

func TestPython_3_12_Uninstall(t *testing.T) {
	msg, err := python.Uninstall("3.12.3")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_12_Compiled_Install(t *testing.T) {
	msg, err := python.Install("3.12.3", true, "--enable-optimizations --enable-loadable-sqlite-extensions --enable-shared --with-computed-gotos --with-lto --enable-ipv6")
	if err != nil {
		t.Fatal(err)
	}

	if strings.Contains(msg, "is already installed") {
		t.Fatal("Python version shouldn't be already installed")
	}

	t.Log(msg)
}

func TestPython_3_12_Compiled_Use(t *testing.T) {
	version, err := python.Use("3.12.3")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"))
	if err != nil {
		t.Fatal(err)
	}

	path, err := utils.GetVersionPath("v"+version, "python")
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != path {
		t.Fatal("invalid version")
	}
}

func TestPython_3_12_Compiled_Uninstall(t *testing.T) {
	msg, err := python.Uninstall("3.12.3")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}
