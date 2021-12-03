package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
	measurements, err := loadMeasurements(filepath)
	if err != nil {
		return 0, err
	}

	return increases(measurements), nil
}

func increases(measurements []int) int {
	result := 0

	for i := 2; i < len(measurements)-1; i++ {
		if sumOfThree(measurements, i) > sumOfThree(measurements, i-1) {
			result++
		}
	}

	return result
}

func sumOfThree(values []int, middle int) int {
	return values[middle-1] + values[middle] + values[middle+1]
}

func loadMeasurements(filepath string) ([]int, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result []int

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.Atoi(line)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %s: %v", line, err)
		}

		result = append(result, i)
	}

	return result, nil
}
