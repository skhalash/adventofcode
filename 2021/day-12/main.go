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

type path []node

func (p *path) push(n node) {
	*p = append(*p, n)
}

func (p *path) pop() {
	*p = (*p)[:len(*p)-1]
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

	return allPathCount(g), nil
}

var seen map[node]bool
var current path
var count int

func allPathCount(g *graph) int {
	seen = make(map[node]bool)
	dfs(g.start(), g)
	return count
}

func dfs(n node, g *graph) {
	if n.small() {
		if seen, exists := seen[n]; exists && seen {
			return
		}
		seen[n] = true
	}

	current.push(n)

	if n.end() {
		count++
		seen[n] = false
		current.pop()
		return
	}

	for _, next := range g.neighbours(n) {
		dfs(next, g)
	}

	current.pop()
	seen[n] = false
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
