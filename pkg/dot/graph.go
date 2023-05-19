package dot

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
	"github.com/google/uuid"
)

type Ref string

func NewRef() Ref {
	return Ref(uuid.New().String())
}

func (ref Ref) Value() string {
	return string(ref)
}

func (ref Ref) String() string {
	return fmt.Sprintf(`"%s"`, string(ref))
}

type MemoryGraph struct {
	pointers  []*Pointer
	listNodes []*ListNode
	changes   []string
}

const defaultGraphName = "G"

const NullptrRef = Ref("0x0")

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		changes:   []string{},
		pointers:  []*Pointer{},
		listNodes: []*ListNode{},
	}
}

func newDirectedGraph() *gographviz.Graph {
	graph := gographviz.NewGraph()
	graph.SetName(defaultGraphName)
	graph.SetDir(true)
	graph.AddAttr(defaultGraphName, "nodesep", "0.5")

	return graph
}

func (mem *MemoryGraph) Changes() []string {
	res := mem.changes
	mem.changes = []string{}
	return res
}

func (mem *MemoryGraph) changed() {
	mem.changes = append(mem.changes, mem.String())
}

func (mem *MemoryGraph) String() string {
	graph := newDirectedGraph()
	dummy := mem.addPointersGraph(graph)
	mem.addListNodesGraph(graph, dummy)

	return graph.String()
}

func addSubGraph(graph *gographviz.Graph) (Ref, Ref) {
	name := NewRef()
	graph.AddSubGraph(defaultGraphName, name.String(), map[string]string{
		"rank": "same",
	})

	dummyName := name + "_dummy"
	graph.AddNode(name.String(), dummyName.String(), map[string]string{
		"style": "invis",
		"shape": "none",
		"label": "1",
	})

	return name, dummyName
}

func (mem *MemoryGraph) addPointersGraph(graph *gographviz.Graph) Ref {
	name, dummy := addSubGraph(graph)
	if len(mem.pointers) == 0 {
		return dummy
	}

	for _, ptr := range mem.pointers {
		graph.AddNode(name.String(), ptr.node.Ref.String(), map[string]string{
			"shape": "plaintext",
			"label": ptr.node.String(),
		})

		if ptr.Address != NullptrRef {
			graph.AddPortEdge(ptr.node.Ref.String(), pointerValuePort+":e", ptr.Address.String(), addressPort+":n", true, nil)
		}
	}

	graph.AddEdge(dummy.String(), mem.pointers[0].node.Ref.String(), true, map[string]string{
		"style": "invis",
	})

	// Add edges to preserve pointers order
	for i := 1; i < len(mem.pointers); i++ {
		prev := mem.pointers[i-1].node.Ref.String()
		curr := mem.pointers[i].node.Ref.String()
		graph.AddEdge(prev, curr, true, map[string]string{
			"style": "invis",
		})
	}

	return dummy
}

func (mem *MemoryGraph) addListNodesGraph(graph *gographviz.Graph, prevDummy Ref) {
	const nullptr = Ref("0x0")

	nodes := map[string]*ListNode{}
	dgraph := NewDirectedGraph[string]()
	for _, node := range mem.listNodes {
		dgraph.AddVertex(node.Address.String())
		if node.Next != nullptr {
			dgraph.AddEdge(node.Address.String(), node.Next.String())
		}
		nodes[node.Address.String()] = node
	}

	connected := dgraph.SubGraphs()
	for _, list := range connected {
		graphName, dummy := addSubGraph(graph)
		graph.AddEdge(prevDummy.String(), dummy.String(), true, map[string]string{
			"style": "invis",
		})
		graph.AddEdge(dummy.String(), list[0], true, map[string]string{
			"style": "invis",
		})
		prevDummy = dummy

		for _, el := range list {
			graph.AddNode(graphName.String(), el, map[string]string{
				"shape": "plaintext",
				"label": nodes[el].node.String(),
			})
		}
		for i := 1; i < len(list); i++ {
			graph.AddPortEdge(list[i-1], nextPort, list[i], addressPort, true, map[string]string{
				"group":    "mid_straight",
				"tailclip": "false",
			})
		}
	}
}

func (mem *MemoryGraph) GetPointer(name string) *Pointer {
	for _, ptr := range mem.pointers {
		if ptr.Name == name {
			return ptr
		}
	}
	return nil
}

func (mem *MemoryGraph) GetListNode(address Ref) *ListNode {
	for _, node := range mem.listNodes {
		if node.Address == address {
			return node
		}
	}
	return nil
}

func (mem *MemoryGraph) Pointer(name string, address Ref) *Pointer {
	ptr := newPointer(mem, name, address)
	mem.pointers = append(mem.pointers, ptr)
	mem.changed()
	return ptr
}

func (mem *MemoryGraph) ListNode(address Ref, data string, next Ref) *ListNode {
	node := NewListNode(mem, address, data, next)
	mem.listNodes = append(mem.listNodes, node)
	mem.changed()
	return node
}
