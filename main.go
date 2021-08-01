package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func main() {
	if err := FormatCSV(); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func FormatCSV() error {
	lines, err := csv.NewReader(os.Stdin).ReadAll()
	if err != nil {
		return fmt.Errorf("extracting csv lines: %v", err)
	}
	if err := WriteLines(Align(lines)); err != nil {
		return fmt.Errorf("writing csv lines: %v", err)
	}
	return nil
}

func WriteLines(lines [][]string) error {
	// Can't use the standard csv writer, because it will escape any fields
	// that have leading spaces.
	for _, line := range lines {
		for j, field := range line {
			var delim = ","
			if j == len(line)-1 {
				delim = ""
			}
			if strings.Contains(field, ",") {
				field = strconv.Quote(field)
			}
			if _, err := fmt.Printf(" %s%s", field, delim); err != nil {
				return err
			}
		}
		if _, err := fmt.Println(); err != nil {
			return err
		}
	}
	return nil
}

func Align(lines [][]string) [][]string {
	var max []int
	for _, line := range lines {
		for j, field := range line {
			field = strings.TrimSpace(field)
			if j == len(max) {
				max = append(max, len(field))
			} else if max[j] < len(field) {
				max[j] = len(field)
			}
		}
	}
	aligned := make([][]string, len(lines))
	for i, line := range lines {
		aligned[i] = make([]string, len(line))
		for j, field := range line {
			aligned[i][j] = fmt.Sprintf("%*s", max[j], strings.TrimSpace(field))
		}
	}
	return aligned
}
