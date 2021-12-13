package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type node struct {
	name string
}

func (n node) start() bool {
	return n.name == "start"
}

func (n node) end() bool {
	return n.name == "end"
}

func (n node) small() bool {
	return strings.ToLower(n.name) == n.name
}

type graph struct {
	adjacent map[node][]node
}

func newGraph() *graph {
	return &graph{
		adjacent: make(map[node][]node),
	}
}

func (g *graph) start() node {
	for _, n := range g.nodes() {
		if n.start() {
			return n
		}
	}

	return node{}
}

func (g *graph) nodes() []node {
	var nodes []node
	for node := range g.adjacent {
		nodes = append(nodes, node)
	}
	return nodes
}

func (g *graph) neighbours(n node) []node {
	return g.adjacent[n]
}

func (g *graph) addEdge(from, to node) {
	g.adjacent[from] = append(g.adjacent[from], to)
}

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (int, error) {
	g, err := loadGraph(filepath)
	if err != nil {
		return 0, err
	}

	return dfs(g.start(), g, make(map[node]bool), true), nil
}

func dfs(n node, g *graph, seen map[node]bool, canRevisit bool) int {
	if n.end() {
		return 1
	}

	if n.small() {
		if seen, exists := seen[n]; exists && seen {
			if !canRevisit || n.start() {
				return 0
			}
			if !n.start() {
				canRevisit = false
			}
		}
		seen[n] = true
	}

	total := 0
	for _, next := range g.neighbours(n) {
		total += dfs(next, g, copy(seen), canRevisit)
	}
	return total
}

func copy(source map[node]bool) map[node]bool {
	result := make(map[node]bool)
	for k, v := range source {
		result[k] = v
	}
	return result
}

func loadGraph(filepath string) (*graph, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	graph := newGraph()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Split(line, "-")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line format %s", line)
		}

		graph.addEdge(node{parts[0]}, node{parts[1]})
		graph.addEdge(node{parts[1]}, node{parts[0]})
	}

	return graph, nil
}
