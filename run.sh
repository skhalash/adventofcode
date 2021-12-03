#!/bin/sh

for dir in ./2021/*/ ; do
    echo "$dir"
    go run "${dir}main.go" "${dir}input.txt"
    echo ""
done
