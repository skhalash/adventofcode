package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const dayCount = 80

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (int, error) {
	fishes, err := loadFishes(filepath)
	if err != nil {
		return 0, err
	}

	counts := make([]int, 9, 9)
	for _, f := range fishes {
		counts[f]++
	}

	for i := 0; i < dayCount; i++ {
		counts = nextDay(counts)
	}

	sum := 0
	for _, c := range counts {
		sum += c
	}

	return sum, nil
}

func nextDay(counts []int) []int {
	newCounts := make([]int, len(counts), len(counts))
	newCounts[6] = counts[0]
	newCounts[8] = counts[0]
	newCounts[0] = 0
	for i := 1; i < 9; i++ {
		newCounts[i-1] += counts[i]
	}
	return newCounts
}

func loadFishes(filepath string) ([]int, error) {
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
