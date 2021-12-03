package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
)

const recordBitCount = 12

type ratingType string

var oxygenGeneratorRating ratingType = "oxygenGeneratorRating"
var co2ScrubberRating ratingType = "co2ScrubberRating"

func allRatingTypes() []ratingType {
	return []ratingType{oxygenGeneratorRating, co2ScrubberRating}
}

type bitCriteria func(val uint64, pos, zeros, ones int) bool

func bitCriteriaByRatingType(rt ratingType) bitCriteria {
	switch rt {
	case oxygenGeneratorRating:
		return oxygenGeneratorRatingCriteria

	case co2ScrubberRating:
		return co2ScrubberRatingCriteria
	}

	panic(fmt.Sprintf("unknown rating type %s", rt))
}

func main() {
	result, err := run(os.Args[1])
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}

	fmt.Println(result)
}

func run(filepath string) (uint64, error) {
	entries, err := loadReport(filepath)
	if err != nil {
		return 0, err
	}

	return ratings(entries), nil
}

func ratings(all []uint64) uint64 {
	var result uint64 = 1

	for _, ratingType := range allRatingTypes() {
		bitCriteria := bitCriteriaByRatingType(ratingType)

		excluded := make(map[uint64]bool)
		for pos := recordBitCount - 1; pos >= 0; pos-- {
			zeros, ones := zerosOnesCount(all, pos, excluded)

			var lastNonExcluded uint64
			for _, value := range all {
				if excluded[value] {
					continue
				}

				if !bitCriteria(value, pos, zeros, ones) {
					excluded[value] = true
				} else {
					lastNonExcluded = value
				}
			}

			if len(all)-len(excluded) == 1 {
				result *= lastNonExcluded
				break
			}
		}
	}

	return result
}

func zerosOnesCount(vals []uint64, pos int, excluded map[uint64]bool) (zeros, ones int) {
	var mask uint64 = 1 << pos
	for _, v := range vals {
		if _, exists := excluded[v]; exists {
			continue
		}

		if v&mask == 0 {
			zeros++
		} else {
			ones++
		}
	}

	return
}

func oxygenGeneratorRatingCriteria(val uint64, pos, zeros, ones int) bool {
	var mostCommonBit uint64 = 0
	if ones >= zeros {
		mostCommonBit = 1
	}

	return mostCommonBit<<pos == val&(1<<pos)
}

func co2ScrubberRatingCriteria(val uint64, pos, zeros, ones int) bool {
	var leastCommonBit uint64 = 1
	if ones >= zeros {
		leastCommonBit = 0
	}

	return leastCommonBit<<pos == val&(1<<pos)
}

func loadReport(filepath string) ([]uint64, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to open data file: %v", err)
	}

	var result []uint64

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		i, err := strconv.ParseUint(line, 2, 64)
		if err != nil {
			return nil, fmt.Errorf("failed to parse line %s: %v", line, err)
		}

		result = append(result, i)
	}

	return result, nil
}
