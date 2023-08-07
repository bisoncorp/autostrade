package dijkstra

import (
	graph3 "github.com.bisoncorp.autostrade/graph"
	"reflect"
	"testing"
)

type link struct {
	weight    int
	nodeIndex int
}

func (l *link) NodeIndex() int {
	return l.nodeIndex
}

func (l *link) Weight() int {
	return l.weight
}

type node struct {
	links []*link
}

func (n *node) Links() []graph3.Link {
	links := make([]graph3.Link, len(n.links))
	for i, l := range n.links {
		links[i] = l
	}
	return links
}

type graph struct {
	nodes []*node
}

func (g *graph) Nodes() []graph3.Node {
	nodes := make([]graph3.Node, len(g.nodes))
	for i, n := range g.nodes {
		nodes[i] = n
	}
	return nodes
}

func (g *graph) addLink(a, b, cost int) {
	l := new(link)
	l.weight = cost
	l.nodeIndex = b
	g.nodes[a].links = append(g.nodes[a].links, l)
}

func (g *graph) addNode() {
	n := new(node)
	n.links = make([]*link, 0)
	g.nodes = append(g.nodes, n)
}

func newGraph(mat [][]int) graph3.Graph {
	g := &graph{nodes: make([]*node, 0)}
	for a, ws := range mat {
		g.addNode()
		for b, w := range ws {
			if w > 0 {
				g.addLink(a, b, w)
			}
		}
	}
	return g
}

var graph1 = newGraph([][]int{
	{0, 3, 0, 7, 8},
	{3, 0, 1, 4, 0},
	{0, 1, 0, 2, 0},
	{7, 4, 2, 0, 3},
	{8, 0, 0, 3, 0},
})

var graph2 = newGraph([][]int{
	{0, 3, 0, 2, 0, 0, 0, 0, 4}, // 0
	{3, 0, 0, 0, 0, 0, 0, 4, 0}, // 1
	{0, 0, 0, 6, 0, 1, 0, 2, 0}, // 2
	{2, 0, 6, 0, 1, 0, 0, 0, 0}, // 3
	{0, 0, 0, 1, 0, 0, 0, 0, 8}, // 4
	{0, 0, 1, 0, 0, 0, 8, 0, 0}, // 5
	{0, 0, 0, 0, 0, 8, 0, 0, 0}, // 6
	{0, 4, 2, 0, 0, 0, 0, 0, 0}, // 7
	{4, 0, 0, 0, 8, 0, 0, 0, 0}, // 8
})

func TestShortestPath(t *testing.T) {
	type args struct {
		graph graph3.Graph
		start int
		end   int
	}
	tests := []struct {
		name string
		args args
		want []int
	}{
		{
			name: "Graph1.02",
			args: args{graph: graph1, start: 0, end: 2},
			want: []int{0, 1, 2},
		},
		{
			name: "Graph1.03",
			args: args{graph: graph1, start: 0, end: 3},
			want: []int{0, 1, 2, 3},
		},
		{
			name: "Graph1.31",
			args: args{graph: graph1, start: 3, end: 1},
			want: []int{3, 2, 1},
		},
		{
			name: "Graph1.24",
			args: args{graph: graph1, start: 2, end: 4},
			want: []int{2, 3, 4},
		},
		{
			name: "Graph2.24",
			args: args{graph: graph2, start: 2, end: 4},
			want: []int{2, 3, 4},
		},
		{
			name: "Graph2.06",
			args: args{graph: graph2, start: 0, end: 6},
			want: []int{0, 3, 2, 5, 6},
		},
		{
			name: "Graph2.02",
			args: args{graph: graph2, start: 0, end: 2},
			want: []int{0, 3, 2},
		},
		{
			name: "Graph2.12",
			args: args{graph: graph2, start: 1, end: 2},
			want: []int{1, 7, 2},
		},
		{
			name: "Graph2.78",
			args: args{graph: graph2, start: 7, end: 8},
			want: []int{7, 1, 0, 8},
		},
		{
			name: "Graph2.28",
			args: args{graph: graph2, start: 2, end: 8},
			want: []int{2, 3, 0, 8},
		},
		{
			name: "Graph2.00",
			args: args{graph: graph2, start: 0, end: 0},
			want: []int{0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ShortestPath(tt.args.graph, tt.args.start, tt.args.end); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("ShortestPath() = %v, want %v", got, tt.want)
			}
		})
	}
}

func BenchmarkShortestPath(b *testing.B) {
	for i := 0; i < b.N; i++ {
		ShortestPath(graph2, 8, 6)
	}
}
