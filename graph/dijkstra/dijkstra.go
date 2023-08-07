package dijkstra

import (
	"container/heap"
	"github.com.bisoncorp.autostrade/graph"
	"math"
)

type nodeInfo struct {
	cost           int
	previous       int
	links          []graph.Link
	visited        bool
	qIndex, gIndex int
}

type priorityQueue []*nodeInfo

func (p priorityQueue) Len() int {
	return len(p)
}

func (p priorityQueue) Less(i, j int) bool {
	return p[i].cost < p[j].cost
}

func (p priorityQueue) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
	p[i].qIndex, p[j].qIndex = i, j
}

func (p *priorityQueue) Push(x any) {
	n := x.(*nodeInfo)
	n.qIndex = len(*p)
	*p = append(*p, n)
}

func (p *priorityQueue) Pop() any {
	r := (*p)[len(*p)-1]
	*p = (*p)[:len(*p)-1]
	r.qIndex = -1
	return r
}

func (p *priorityQueue) update(n *nodeInfo, cost int) {
	n.cost = cost
	heap.Fix(p, n.qIndex)
}

func ShortestPath(graph graph.Graph, start, end int) []int {
	graphNodes := graph.Nodes()
	var unvisitedQueue priorityQueue = make([]*nodeInfo, 0, len(graphNodes))
	var nodesInfo = make([]*nodeInfo, len(graphNodes))
	for i, n := range graphNodes {
		n := &nodeInfo{
			cost:     math.MaxInt,
			previous: -1,
			links:    n.Links(),
			visited:  false,
			qIndex:   0,
			gIndex:   i,
		}
		if i == start {
			n.cost = 0
		}
		nodesInfo[i] = n
		heap.Push(&unvisitedQueue, n)
	}
	var currentNode *nodeInfo
	for {
		if unvisitedQueue.Len() == 0 {
			break
		}
		currentNode = heap.Pop(&unvisitedQueue).(*nodeInfo)
		currentNode.visited = true
		for _, link := range currentNode.links {
			if nodesInfo[link.NodeIndex()].visited {
				continue
			}
			cost := link.Weight() + currentNode.cost
			if destination := nodesInfo[link.NodeIndex()]; cost < destination.cost {
				destination.previous = currentNode.gIndex
				unvisitedQueue.update(destination, cost)
			}
		}
		if currentNode.gIndex == end {
			break
		}
	}
	result := make([]int, 0)
	for i := end; i != -1; {
		result = append(result, i)
		i = nodesInfo[i].previous
	}
	reverse(result)
	return result
}

func reverse[S ~[]E, E any](s S) {
	for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
		s[i], s[j] = s[j], s[i]
	}
}
