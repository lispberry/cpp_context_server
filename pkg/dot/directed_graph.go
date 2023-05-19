package dot

type DirectedGraph[Vertex comparable] struct {
	graph map[Vertex][]Vertex
}

func NewDirectedGraph[Vertex comparable]() *DirectedGraph[Vertex] {
	return &DirectedGraph[Vertex]{
		graph: make(map[Vertex][]Vertex),
	}
}

func (g *DirectedGraph[Vertex]) AddVertex(vertex Vertex) {
	if _, exists := g.graph[vertex]; !exists {
		g.graph[vertex] = []Vertex{}
	}
}

func (g *DirectedGraph[Vertex]) AddEdge(to Vertex, from Vertex) {
	g.AddVertex(to)
	g.AddVertex(from)
	g.graph[to] = append(g.graph[to], from)
}

type data struct {
	order int
	color int
}

func (g *DirectedGraph[Vertex]) SubGraphs() [][]Vertex {
	usedColors := map[int]int{}
	unusedColor := 0
	colored := make(map[Vertex]data)

	var dfs func(vertex Vertex, order int)
	dfs = func(vertex Vertex, order int) {
		if _, exists := usedColors[unusedColor]; !exists {
			usedColors[unusedColor] = order
		}
		if usedColors[unusedColor] < order {
			usedColors[unusedColor] = order
		}
		colored[vertex] = data{
			order: order,
			color: unusedColor,
		}
		for _, adj := range g.graph[vertex] {
			delete(usedColors, colored[adj].color)
			order++
			dfs(adj, order)
		}
	}

	for vertex := range g.graph {
		if _, ok := colored[vertex]; !ok {
			unusedColor++
			dfs(vertex, 0)
		}
	}

	var result [][]Vertex
	for color, amount := range usedColors {
		nodes := make([]Vertex, amount+1)
		for node, data := range colored {
			if data.color == color {
				nodes[data.order] = node
			}
		}
		result = append(result, nodes)
	}

	return result
}
