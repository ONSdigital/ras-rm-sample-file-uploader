package file

import (
	"bufio"
	"bytes"
	"io"
	"strings"
)

// readFileForCountTotals reads the file for counting
// the total expected CIs and total sample count.
// Will return a buffer thats written to by the tee reader
// to ensure that the file can be processed further.
func readFileForCountTotals(r io.Reader) (int, int, *bytes.Buffer) {
	var buf bytes.Buffer
	tee := io.TeeReader(r, &buf)
	scanner := bufio.NewScanner(tee)
	sampleCount := 0
	formTypes := make(map[string]string)
	for scanner.Scan() {
		sampleCount++
		line := scanner.Text()
		s := strings.Split(line, ":")
		formTypes[s[25]] = s[25]
	}
	return len(formTypes), sampleCount, &buf
}
