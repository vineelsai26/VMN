package tests

import (
	"testing"

	"vineelsai.com/vmn/src/node"
)

func TestNode_20_Install(t *testing.T) {
	msg, err := node.Install("20")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestNode_18_Install(t *testing.T) {
	msg, err := node.Install("18")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestNode_16_Install(t *testing.T) {
	msg, err := node.Install("16")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}

func TestNode_14_Install(t *testing.T) {
	msg, err := node.Install("14")
	if err != nil {
		t.Fatal(err)
	}

	t.Log(msg)
}
