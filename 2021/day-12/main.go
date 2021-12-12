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

type graph struct {
	adjacent map[node][]node
}

func newGraph() *graph {
	return &graph{
		adjacent: make(map[node][]node),
	}
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
	_, err := loadGraph(filepath)
	if err != nil {
		return 0, err
	}

	return 0, nil
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
