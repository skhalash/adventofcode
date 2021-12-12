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

func (s path) push(n node) path {
	return append(s, n)
}

func (p path) pop() (path, node) {
	if p.empty() {
		return p, node{}
	}

	lastIndex := len(p) - 1
	return p[:lastIndex], p[lastIndex]
}

func (p *path) empty() bool {
	return len(*p) == 0
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

	return len(allPaths(g)), nil
}

func allPaths(g *graph) []path {
	seen := make(map[node]bool)
	var current path
	var all []path

	dfs(g.start(), g, seen, current, all)
	return all
}

func dfs(n node, g *graph, seen map[node]bool, current path, all []path) {

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
	}

	return graph, nil
}
