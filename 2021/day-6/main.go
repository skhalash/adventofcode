package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

const dayCount = 256

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (int, error) {
	lanternfishDaysLeft, err := loadLanternfishDaysLeft(filepath)
	if err != nil {
		return 0, err
	}

	countByDaysLeft := make([]int, 9)
	for _, daysLeft := range lanternfishDaysLeft {
		countByDaysLeft[daysLeft]++
	}

	for i := 0; i < dayCount; i++ {
		nextDay(countByDaysLeft)
	}

	sum := 0
	for _, c := range countByDaysLeft {
		sum += c
	}

	return sum, nil
}

func nextDay(countByDaysLeft []int) {
	newbornCount := countByDaysLeft[0]

	for daysLeft := 1; daysLeft < 9; daysLeft++ {
		countByDaysLeft[daysLeft-1] = countByDaysLeft[daysLeft]
	}

	countByDaysLeft[6] += newbornCount
	countByDaysLeft[8] = newbornCount
}

func loadLanternfishDaysLeft(filepath string) ([]int, error) {
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
