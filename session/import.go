package session

import (
	"encoding/csv"
	"io"
	"os"
	"strconv"
	"strings"
)

// Import data from CSV file
func (s *Session) ImportFromCsv(fileName string, delimiter string, primary int, verbose bool, linking []int) {

	// Checking if file exist
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		s.Stream.Error(err.Error())
		return
	}

	csvfile, err := os.Open(fileName)
	if err != nil {
		s.Stream.Error(err.Error())
		return
	}

	// Parse the file
	r := csv.NewReader(csvfile)

	// Iterate through the records
	lineCount := 0
	importedLine := 0
	line := []string{}

	for {
		var main string
		keys := []string{}

		lineCount = lineCount + 1
		// Read each record from csv
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			s.Stream.Error(err.Error())
			return
		}

		if len(record) == 1 {
			if strings.Contains(record[0], delimiter) {
				line = strings.Split(record[0], delimiter)
			}
		}

		if lineCount == 1 {
			continue
		}

		if len(line) >= (primary + 1) {

			main = line[primary]
			for key, t := range line {
				if key != primary {
					keys = append(keys, t)
				}
			}

			targetId, err := s.AddTarget("import", main)
			if err != nil {
				continue
			}

			target, err := s.GetTarget(targetId)
			if err != nil {
				continue
			}

			for key, tag := range line {
				if IntSliceKeyExist(linking, key) {
					_, _ = s.AddTag(target, tag)
				}
			}

			if verbose {
				s.Stream.Standard("Imported line: [" + strings.Join(line, ",") + "]")
			}
			importedLine = importedLine + 1
		}
	}

	s.Stream.Backgound("'" + strconv.Itoa(importedLine) + "' target(s) imported")
	return
}

func IntSliceKeyExist(l []int, i int) bool {
	for _, x := range l {
		if x == i {
			return true
		}
	}
	return false
}
