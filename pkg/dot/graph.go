package dot

import (
	"fmt"
	"github.com/awalterschulze/gographviz"
)

type Ref string

func (ref Ref) String() string {
	return fmt.Sprintf(`"%s"`, string(ref))
}

type MemoryGraph struct {
	graph *gographviz.Graph
	nodes map[string]*Node
}

const defaultGraphName = "G"

func newDirectedGraph() *gographviz.Graph {
	graph := gographviz.NewGraph()
	graph.SetName(defaultGraphName)
	graph.SetDir(true)
	graph.AddAttr(defaultGraphName, "nodesep", "0.5")
	graph.AddNode(defaultGraphName, "node", map[string]string{
		"shape": "plaintext",
	})

	return graph
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		graph: newDirectedGraph(),
		nodes: make(map[string]*Node),
	}
}

type Node struct {
	Address string
	Data    string
	Next    string
}

const addressPort = "ref1"
const dataPort = "data"
const nextPort = "ref2"

func newTableFromNode(node *Node) string {
	const Table = `<<table border="0" cellspacing="0" cellborder="1">
		<tr>
			<td port="ref1" width="28" height="36">%s</td>
			<td port="data" width="28" height="36">%s</td>
			<td port="ref2" width="28" height="36">%s</td>
		</tr>
		<tr>
			<td BORDER="0">addr</td>
			<td BORDER="0">val</td>
			<td BORDER="0">Next</td>
		</tr>
	</table>>`

	return fmt.Sprintf(Table, node.Address, node.Data, node.Next)
}

func (mem *MemoryGraph) AddNode(ref Ref, node *Node) {
	table := newTableFromNode(node)
	mem.nodes[ref.String()] = node
	mem.graph.AddNode(defaultGraphName, ref.String(), map[string]string{
		"label": table,
	})
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
	mem.graph.AddPortEdge(from.String(), nextPort, to.String(), addressPort, true, attrs)
}

func (mem *MemoryGraph) String() string {
	return mem.graph.String()
}
