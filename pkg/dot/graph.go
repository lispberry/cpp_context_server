package dot

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
)

type Ref string

func NewRef() string {
	// TODO(ivo)
	return ""
}

func (ref Ref) Value() string {
	return string(ref)
}

func (ref Ref) String() string {
	return fmt.Sprintf(`"%s"`, string(ref))
}

type MemoryGraph struct {
	graph *gographviz.Graph

	pointers map[string]*Pointer
	nodes    map[string]*Node
}

const defaultGraphName = "G"

func newDirectedGraph() *gographviz.Graph {
	graph := gographviz.NewGraph()
	graph.SetName(defaultGraphName)
	graph.SetDir(true)
	graph.AddAttr(defaultGraphName, "nodesep", "0.5")

	return graph
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		graph:    newDirectedGraph(),
		pointers: make(map[string]*Pointer),
		nodes:    make(map[string]*Node),
	}
}

var defaultAttrs = map[string]string{
	"shape": "plaintext",
}

func (mem *MemoryGraph) addDefaultAttrs(attrs map[string]string) map[string]string {
	for k, v := range defaultAttrs {
		attrs[k] = v
	}
	return attrs
}

func (mem *MemoryGraph) AddPointer(pointer *Pointer) Ref {
	ref := Ref(pointer.Name)

	mem.pointers[ref.String()] = pointer
	mem.graph.AddNode(defaultGraphName, ref.String(), map[string]string{
		"label": pointer.Table(),
	})

	return ref
}

func (mem *MemoryGraph) AddNode(node *Node) Ref {

	ref := Ref(node.Address)

	mem.nodes[ref.String()] = node
	mem.graph.AddNode(defaultGraphName, ref.String(), mem.addDefaultAttrs(map[string]string{
		"label": node.Table(),
	}))

	return ref
}

func (mem *MemoryGraph) DeleteNode(ref Ref) {
	delete(mem.nodes, ref.String())
	mem.graph.Nodes.Remove(ref.String())
}

func (mem *MemoryGraph) AddEdge(from Ref, to Ref) {
	attrs := map[string]string{
		"group":    "mid_straight",
		"tailclip": "false",
	}
	mem.graph.AddPortEdge(from.String(), nextPort, to.String(), addressPort, true, mem.addDefaultAttrs(attrs))
}

func (mem *MemoryGraph) Add() {
}

func (mem *MemoryGraph) String() string {
	dgraph := NewDirectedGraph()
	for _, node := range mem.nodes {
		dgraph.AddEdge(Vertex(node.Address), Vertex(node.Next))
	}

	const pointersGraph = "pointersGraph"
	mem.graph.AddSubGraph(defaultGraphName, pointersGraph, map[string]string{
		"rank": "same",
	})

	subGraphs := dgraph.SubGraphs()
	for _, subgraph := range subGraphs {
		subgraph[0] = ""
	}

	return mem.graph.String()
}
