package dot

import (
	"strings"
	"testing"
)

func TestEmptyGraph(t *testing.T) {
	graph := NewMemoryGraph()

	str := graph.String()
	if !strings.Contains(str, "nodesep") {
		t.Fail()
	}
}

func TestDeleteGraph(t *testing.T) {
	graph := NewMemoryGraph()
	originStr := graph.String()

	const ref = "0xFFF"
	graph.AddNode(ref, &Node{
		Address: "0xFFFF",
		Data:    "10",
		Next:    "0x0",
	})
	graph.DeleteNode(ref)
	newStr := graph.String()

	if originStr != newStr {
		t.Fail()
	}
}
