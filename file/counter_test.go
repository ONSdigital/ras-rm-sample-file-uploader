package file

import (
	"bufio"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestGetCount(t *testing.T) {
	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	ciCount, totalCount, _ := readFileForCountTotals(file)

	assert.Equal(t, 9, totalCount)
	assert.Equal(t, 2, ciCount)
}

func TestCanReadFileAfterGettingCount(t *testing.T) {
	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	_, _, buf := readFileForCountTotals(file)

	additionalCount := 0
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		additionalCount++
	}
	assert.Equal(t, 9, additionalCount)
}