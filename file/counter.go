package file

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"strings"
)

const FORMTYPE_CSV_POSITION = 25

// readFileForCountTotals reads the file for counting
// the total expected CIs and total sample count.
// Will return a buffer thats written to by the tee reader
// to ensure that the file can be processed further.
func readFileForCountTotals(r io.Reader) (int, int, *bytes.Buffer, error) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	scanner := bufio.NewScanner(tee)
	sampleCount := 0
	formTypes := make(map[string]string)
	for scanner.Scan() {
		sampleCount++
		line := scanner.Text()
		s := strings.Split(line, ":")
		if len(s) <= FORMTYPE_CSV_POSITION {
			return 0, 0, nil, errors.New("Too few columns in CSV file")
		}
		formTypes[s[FORMTYPE_CSV_POSITION]] = s[FORMTYPE_CSV_POSITION]
	}
	return len(formTypes), sampleCount, &buf, nil
}
