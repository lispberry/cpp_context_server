package visualization

type DirectedGraph struct {
	graph map[Vertex][]Vertex
}

type Vertex string

func NewDirectedGraph() *DirectedGraph {
	return &DirectedGraph{
		graph: make(map[Vertex][]Vertex),
	}
}

func (g *DirectedGraph) AddVertex(vertex Vertex) {
	if _, exists := g.graph[vertex]; !exists {
		g.graph[vertex] = []Vertex{}
	}
}

func (g *DirectedGraph) AddEdge(to Vertex, from Vertex) {
	g.AddVertex(to)
	g.AddVertex(from)
	g.graph[to] = append(g.graph[to], from)
	g.graph[from] = append(g.graph[from], to)
}

func (g *DirectedGraph) SubGraphs() [][]Vertex {
	unusedColor := 0
	colored := make(map[Vertex]int)
	var result [][]Vertex

	var dfs func(vertex Vertex)
	dfs = func(vertex Vertex) {
		if _, exists := colored[vertex]; exists {
			return
		}
		colored[vertex] = unusedColor
		for _, adj := range g.graph[vertex] {
			dfs(adj)
		}
	}

	for vertex := range g.graph {
		dfs(vertex)
		unusedColor++
	}

	colorMap := make(map[int][]Vertex)
	for vertex, color := range colored {
		colorMap[color] = append(colorMap[color], vertex)
	}

	for _, vertices := range colorMap {
		result = append(result, vertices)
	}

	return result
}
