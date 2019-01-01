package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"os"
)

func main() {
	inPlace := flag.Bool("w", false, "write result to source instead of stdout")
	flag.Parse()

	if err := FormatCSV(*inPlace, flag.Args()); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
}

func FormatCSV(inPlace bool, filenames []string) error {
	if len(filenames) == 0 {
		lines, err := ExtractLines(os.Stdin)
		if err != nil {
			return fmt.Errorf("extracting csv lines: %v", err)
		}
		if err := WriteLines(os.Stdout, Align(lines)); err != nil {
			return fmt.Errorf("writing csv lines: %v", err)
		}
	} else {
		return fmt.Errorf("READING FROM FILES NOT SUPPORTED")
	}
	return nil
}

func ExtractLines(r io.Reader) ([][]string, error) {
	return csv.NewReader(r).ReadAll()
}

func WriteLines(w io.Writer, lines [][]string) error {
	for _, line := range lines {
		for j, field := range line {
			var delim = ","
			if j == len(line)-1 {
				delim = ""
			}
			if _, err := fmt.Fprintf(w, " %s%s", field, delim); err != nil {
				return err
			}
		}
		if _, err := w.Write([]byte{'\n'}); err != nil {
			return err
		}
	}
	return nil
}

func Align(lines [][]string) [][]string {
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
	aligned := make([][]string, len(lines))
	for i, line := range lines {
		aligned[i] = make([]string, len(line))
		for j, field := range line {
			// TODO: could make this more efficient
			aligned[i][j] = fmt.Sprintf("%*s", max[j], field)
		}
	}
	return aligned
}
