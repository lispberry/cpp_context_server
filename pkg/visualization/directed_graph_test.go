package visualization

import (
	"testing"
)

func TestDirectedGraph_SubGraphs(t *testing.T) {
	// Create a new DirectedGraph
	graph := NewDirectedGraph()

	// Add vertices and edges
	graph.AddEdge("A", "B")
	graph.AddEdge("B", "C")
	graph.AddEdge("C", "D")
	graph.AddEdge("E", "F")

	// Get the subgraphs
	subGraphs := graph.SubGraphs()

	//// Define the expected subgraphs
	//expected := [][]Vertex{
	//	{"E", "F"},
	//	{"A", "B", "C", "D"},
	//}

	if len(subGraphs) != 2 {
		t.Errorf("Expected 2 subgraphs")
	}
}
