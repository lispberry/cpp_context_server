package dot

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"strings"
	"testing"
)

func validate(dot string) error {
	graphAst, _ := gographviz.ParseString(dot)
	graph := gographviz.NewGraph()
	return gographviz.Analyse(graphAst, graph)
}

func TestEmptyGraph(t *testing.T) {
	graph := NewMemoryGraph()

	str := graph.String()
	if !strings.Contains(str, "nodesep") {
		t.Fail()
	}
}

func TestDeleteGraph(t *testing.T) {
	graph := NewMemoryGraph()
	pointer := graph.Pointer("hello", Ref("0xFFFF"))
	pointer.SetAddress("0xFFF")

	for _, change := range graph.Changes() {
		if err := validate(change); err != nil {
			t.Errorf("Invalid dot `%s` with error `%s`", change, err.Error())
		}
	}
}

func TestPointerGraph(t *testing.T) {
	graph := NewMemoryGraph()
	graph.Pointer("h1", Ref("0xFFFF"))
	graph.Pointer("h2", Ref("0xFFFF1"))
	graph.Pointer("h3", Ref("0x0"))
	graph.ListNode(Ref("0xFFFF"), "10", Ref("0xFFFF1"))
	graph.ListNode(Ref("0xFFFF1"), "10", Ref("0x0"))

	fmt.Println(graph.String())
}

func TestPointer1Graph(t *testing.T) {
	graph := NewMemoryGraph()
	ptr := graph.Pointer("h1", "0xFFFF1")
	graph.Pointer("h2", "0xFFFF")
	graph.Pointer("h3", "0x0")
	graph.ListNode("0xFFFF", "10", Ref("0xFFFF1"))
	graph.ListNode(Ref("0xFFFF1"), "10", Ref("0x0"))
	ptr.SetAddress("0xFFFF")

	fmt.Println(graph.Changes())
}

func TestNodeSetDataGraph(t *testing.T) {
	graph := NewMemoryGraph()
	graph.Pointer("h1", "0xFFFF1")
	graph.Pointer("h2", "0xFFFF")
	graph.Pointer("h3", "0x0")
	node := graph.ListNode("0xFFFF", "10", Ref("0xFFFF1"))
	graph.ListNode(Ref("0xFFFF1"), "10", Ref("0x0"))
	node.SetData("20")

	fmt.Println(graph.Changes())
}

func TestNodeSetNextGraph(t *testing.T) {
	graph := NewMemoryGraph()
	graph.Pointer("h1", "0xFFFF1")
	graph.Pointer("h2", "0xFFFF")
	graph.Pointer("h3", "0x0")
	graph.ListNode("0xFFFF", "10", Ref("0xFFFF1"))
	node := graph.ListNode(Ref("0xFFFF1"), "10", Ref("0x0"))
	graph.ListNode(Ref("0xFFFF2"), "30", Ref("0x0"))
	node.SetNext("0xFFFF2")

	fmt.Println(graph.Changes())
}
