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
	Trim(lines)
	Escape(lines)
	Align(lines)
	if err := Write(lines); err != nil {
		return fmt.Errorf("writing csv lines: %v", err)
	}
	return nil
}

func Trim(lines [][]string) {
	for _, line := range lines {
		for j, field := range line {
			line[j] = strings.TrimSpace(field)
		}
	}
}

func Escape(lines [][]string) {
	for _, line := range lines {
		for j, field := range line {
			if strings.Contains(field, ",") {
				line[j] = strconv.Quote(field)
			}
		}
	}
}

func Write(lines [][]string) error {
	// Can't use the standard csv writer, because it will escape any fields
	// that have leading spaces.
	for _, line := range lines {
		for j, field := range line {
			var delim = ","
			if j == len(line)-1 {
				delim = ""
			}
			//if strings.Contains(field, ",") {
			//	field = strconv.Quote(field)
			//}
			if _, err := fmt.Printf("%s%s", field, delim); err != nil {
				return err
			}
		}
		if _, err := fmt.Println(); err != nil {
			return err
		}
	}
	return nil
}

func Align(lines [][]string) {
	var max []int
	for _, line := range lines {
		for j, field := range line {
			if j == len(max) {
				max = append(max, len(field))
			} else if max[j] < len(field) {
				max[j] = len(field)
			}
		}
	}
	for j, m := range max {
		max[j] = m + 1
	}
	for _, line := range lines {
		for j, field := range line {
			if isQuoted(field) {
				line[j] = fmt.Sprintf("\"%*s", max[j]-1, field[1:])
			} else {
				line[j] = fmt.Sprintf("%*s", max[j], field)
			}
		}
	}
}

func isQuoted(field string) bool {
	return len(field) > 0 && field[0] == '"'
}
