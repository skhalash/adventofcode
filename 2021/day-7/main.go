package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (int, error) {
	crabPositions, err := loadCrabPositions(filepath)
	if err != nil {
		return 0, err
	}

	return alignCrabs(crabPositions), nil
}

func alignCrabs(crabPositions []int) int {
	return 0
}

func loadCrabPositions(filepath string) ([]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("failed to read file content: %v", err)
	}

	strs := strings.Split(strings.TrimSpace(string(content)), ",")

	var result []int
	for _, s := range strs {
		n, err := strconv.Atoi(s)
		if err != nil {
			return nil, fmt.Errorf("failed to parse number: %v", err)
		}
		result = append(result, n)
	}
	return result, nil
}
