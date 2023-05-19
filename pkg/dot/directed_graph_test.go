package dot

import (
	"reflect"
	"testing"
)

func TestDirectedGraph_SubGraphsSingle(t *testing.T) {
	// Create a new DirectedGraph
	graph := NewDirectedGraph[string]()

	// Add vertices and edges
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")

	// Get the subgraphs
	subGraphs := graph.SubGraphs()

	if len(subGraphs) != 1 {
		t.Errorf("Expected 1 subgraph")
	}

	if !reflect.DeepEqual(subGraphs[0], []string{"A", "B", "C"}) {
		t.Errorf("Got %v", subGraphs[0])
	}
}

func TestDirectedGraph_SubGraphsSingleAnother(t *testing.T) {
	// Create a new DirectedGraph
	graph := NewDirectedGraph[string]()

	// Add vertices and edges
	graph.AddEdge("0xFFFF", "0xFFFF1")
	graph.AddEdge("0xFFFF1", "0xFFFF2")
	graph.AddVertex("0xFFFFFF")

	// Get the subgraphs
	subGraphs := graph.SubGraphs()

	if len(subGraphs) != 1 {
		t.Errorf("Expected 1 subgraph")
	}

	if !reflect.DeepEqual(subGraphs[0], []string{"0xFFFF", "0xFFFF1", "0xFFFF2"}) {
		t.Errorf("Got %v", subGraphs[0])
	}
}

func TestDirectedGraph_SubGraphs(t *testing.T) {

}
