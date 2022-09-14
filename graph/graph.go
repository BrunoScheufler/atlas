package graph

import (
	"errors"
	"fmt"
	"strings"
)

type Graph[T comparable] struct {
	nodes *OrderedSet[T]
	edges [][]T
}

func New[T comparable]() *Graph[T] {
	return &Graph[T]{
		edges: make([][]T, 0),
		nodes: NewOrderedSet[T](),
	}
}

func (g *Graph[T]) AddNode(node T) {
	g.nodes.Add(node)
}

func (g *Graph[T]) SetNodes(nodes []T) {
	g.nodes = OrderedSetFromSlice(nodes)
}

func (g *Graph[T]) HasNode(node T) bool {
	return g.nodes.Has(node)
}

func (g *Graph[T]) AddEdge(from, to T) {
	g.edges = append(g.edges, []T{from, to})
}

func (g *Graph[T]) HasEdge(from, to T) bool {
	for _, edge := range g.edges {
		if edge[0] == from && edge[1] == to {
			return true
		}
	}
	return false
}

func (g *Graph[T]) RemoveEdge(from, to T) {
	for i, edge := range g.edges {
		if edge[0] == from && edge[1] == to {
			g.edges = append(g.edges[:i], g.edges[i+1:]...)
			return
		}
	}
}
func (g *Graph[T]) NodesWithoutIncomingEdges() []T {
	var nodes []T
	for _, node := range g.nodes.Values() {
		if g.hasNoIncomingEdges(node) {
			nodes = append(nodes, node)
		}
	}
	return nodes
}

func (g *Graph[T]) hasNoIncomingEdges(node T) bool {
	for _, edge := range g.edges {
		if edge[1] == node {
			return false
		}
	}
	return true
}

func (g *Graph[T]) NodesWithEdgeFromN(n T) []T {
	var nodes []T
	for _, edge := range g.edges {
		if edge[0] == n {
			nodes = append(nodes, edge[1])
		}
	}
	return nodes
}

func (g *Graph[T]) NodesWithEdgeToN(n T) []T {
	var nodes []T
	for _, edge := range g.edges {
		if edge[1] == n {
			nodes = append(nodes, edge[0])
		}
	}
	return nodes
}

func (g *Graph[T]) CountIncomingEdges(n T) int {
	var count int
	for _, edge := range g.edges {
		if edge[1] == n {
			count++
		}
	}
	return count
}

func (g *Graph[T]) Indegree() *OrderedMap[T, int] {
	nodeIncomingEdges := NewOrderedMap[T, int]()
	for _, node := range g.nodes.Values() {
		nodeIncomingEdges.Set(node, g.CountIncomingEdges(node))
	}

	return nodeIncomingEdges
}

func (g *Graph[T]) Clone() *Graph[T] {
	newGraph := New[T]()
	for _, node := range g.nodes.Values() {
		newGraph.AddNode(node)
	}

	for _, edge := range g.edges {
		newGraph.AddEdge(edge[0], edge[1])
	}

	return newGraph
}

type errorCycles struct{}

func (e errorCycles) Error() string {
	return "cycle detected"
}

func (g *Graph[T]) TopologicalSort() ([]T, error) {
	cloned := g.Clone()

	var sorted []T

	tmp := cloned.NodesWithoutIncomingEdges()

	for len(tmp) > 0 {
		n := tmp[0]
		tmp = tmp[1:]

		sorted = append(sorted, n)

		for _, m := range cloned.NodesWithEdgeFromN(n) {
			cloned.RemoveEdge(n, m)
			if cloned.hasNoIncomingEdges(m) {
				tmp = append(tmp, m)
			}
		}
	}

	if len(cloned.edges) > 0 {
		return nil, errorCycles{}
	}

	return sorted, nil
}

func (g *Graph[T]) TopologicalSortWithLayers() ([][]T, error) {
	input := g
	cloned := input.Clone()

	//  Start with a set S0 containing all nodes with no incoming edges
	S0 := cloned.NodesWithoutIncomingEdges()

	layers := make([][]T, 0)
	layers = append(layers, S0, []T{})

	Sn := func(layer int) []T {
		return layers[layer]
	}

	i := 0
	for {
		for _, n := range Sn(i) {
			for _, m := range cloned.NodesWithEdgeFromN(n) {
				cloned.RemoveEdge(n, m)

				if cloned.hasNoIncomingEdges(m) {
					layers[i+1] = append(layers[i+1], m)
				}
			}
		}

		if Sn(i+1) == nil || len(Sn(i+1)) == 0 {
			break
		}

		layers = append(layers, []T{})
		i++
	}

	if len(cloned.edges) > 0 {
		return nil, errorCycles{}
	}

	// only return layers with entries
	var layersWithEntries [][]T
	for _, layer := range layers {
		if len(layer) > 0 {
			layersWithEntries = append(layersWithEntries, layer)
		}
	}

	return layersWithEntries, nil
}

func (g *Graph[T]) HasCycles() bool {
	_, err := g.TopologicalSort()
	if errors.Is(err, errorCycles{}) {
		return true
	}
	return false
}

func (g *Graph[T]) TransitiveReduction() (*Graph[T], error) {
	input := g

	if input.HasCycles() {
		return nil, errorCycles{}
	}

	transitiveReduction := New[T]()
	transitiveReduction.SetNodes(input.Nodes())

	descendants := NewOrderedMap[T, *OrderedSet[T]]()
	var checkCount = input.Indegree()

	for _, u := range input.Nodes() {
		uNeighbours := OrderedSetFromSlice(input.NodesWithEdgeFromN(u))

		for _, v := range input.NodesWithEdgeFromN(u) {
			if uNeighbours.Has(v) {
				if !descendants.Has(v) {
					walkedEdges := OrderedSetFromSlice(input.DFS(v, -1))
					descendants.Set(v, walkedEdges)
				}

				for _, d := range descendants.Get(v).Values() {
					uNeighbours.Remove(d)
				}
			}

			checkCount.Set(v, checkCount.Get(v)-1)
			if checkCount.Get(v) == 0 {
				descendants.Delete(v)
			}

		}

		for _, v := range uNeighbours.Values() {
			transitiveReduction.AddEdge(u, v)
		}
	}

	return transitiveReduction, nil
}

func Remove[T comparable](s []T, e T) []T {
	for i, v := range s {
		if v == e {
			s = append(s[:i], s[i+1:]...)
		}
	}
	return s
}

func Has[T comparable](s []T, e T) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func (g *Graph[T]) DFS(start T, maxDepth int) []T {
	cloned := g.Clone()

	visited := make([]T, 0)
	stack := []T{start}

	for len(stack) > 0 {
		n := stack[len(stack)-1]
		stack = stack[:len(stack)-1]

		if Has(visited, n) {
			continue
		}

		visited = append(visited, n)

		for _, m := range cloned.NodesWithEdgeFromN(n) {
			if !Has(visited, m) {
				stack = append(stack, m)
			}
		}

		if maxDepth != -1 && len(visited) >= maxDepth {
			break
		}
	}

	// Return all visited except for start
	return Filter(visited, func(n T) bool {
		return n != start
	})
}

func Filter[T comparable](s []T, fn func(T) bool) []T {
	var r []T
	for _, v := range s {
		if fn(v) {
			r = append(r, v)
		}
	}
	return r
}

func (g *Graph[T]) String() string {
	var b strings.Builder

	for _, node := range g.Nodes() {
		b.WriteString(fmt.Sprintf("%v -> %v\n", node, g.NodesWithEdgeFromN(node)))
	}

	return b.String()
}

func (g *Graph[T]) Nodes() []T {
	return g.nodes.Values()
}

func (g *Graph[T]) Edges() [][]T {
	return g.edges
}
