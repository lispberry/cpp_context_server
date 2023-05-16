package dot

import (
	"fmt"
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

	ref := graph.AddNode(&Node{
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

func TestCreateGraph(t *testing.T) {
	graph := NewMemoryGraph()

	ref := graph.AddNode(&Node{
		Address: "0xF",
		Data:    "10",
		Next:    "0x0",
	})
	graph.AddPointer(&Pointer{
		Name:    "head",
		Address: ref,
	})

	graph.Add()

	fmt.Println(graph.String())
}
