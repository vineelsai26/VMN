//go:build (linux && ignore) || (darwin && ignore) || !windows
// +build linux,ignore darwin,ignore !windows

package tests

import (
	"os"
	"path/filepath"
	"testing"

	"vineelsai.com/vmn/src/python"
	"vineelsai.com/vmn/src/utils"
)

func TestPython_3_11_Install(t *testing.T) {
	msg, err := python.Install("3.11")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_11_Use(t *testing.T) {
	version, err := python.Use("3.11")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != filepath.Join(utils.GetHome(), ".vmn", "python", "v"+version, "bin") {
		t.Fatal("invalid version")
	}
}

func TestPython_3_11_Uninstall(t *testing.T) {
	msg, err := python.Uninstall("3.11")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_10_Install(t *testing.T) {
	msg, err := python.Install("3.10")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_10_Use(t *testing.T) {
	version, err := python.Use("3.10")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != filepath.Join(utils.GetHome(), ".vmn", "python", "v"+version, "bin") {
		t.Fatal("invalid version")
	}
}

func TestPython_3_10_Uninstall(t *testing.T) {
	msg, err := python.Uninstall("3.10")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_9_Install(t *testing.T) {
	msg, err := python.Install("3.9")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_9_Use(t *testing.T) {
	version, err := python.Use("3.9")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != filepath.Join(utils.GetHome(), ".vmn", "python", "v"+version, "bin") {
		t.Fatal("invalid version")
	}
}

func TestPython_3_9_Uninstall(t *testing.T) {
	msg, err := python.Uninstall("3.9")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_8_Install(t *testing.T) {
	msg, err := python.Install("3.8")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestPython_3_8_Use(t *testing.T) {
	version, err := python.Use("3.8")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(version)

	content, err := os.ReadFile(filepath.Join(utils.GetHome(), ".vmn", "python-version"))
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != filepath.Join(utils.GetHome(), ".vmn", "python", "v"+version, "bin") {
		t.Fatal("invalid version")
	}
}

func TestPython_3_8_Uninstall(t *testing.T) {
	msg, err := python.Uninstall("3.8")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}
