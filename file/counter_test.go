package file

import (
	"bufio"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetCount(t *testing.T) {
	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	ciCount, totalCount, _, _ := readFileForCountTotals(file)

	assert.Equal(t, 9, totalCount)
	assert.Equal(t, 2, ciCount)
}

func TestCanReadFileAfterGettingCount(t *testing.T) {
	file, err := os.Open("sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	_, _, buf, _ := readFileForCountTotals(file)

	additionalCount := 0
	scanner := bufio.NewScanner(buf)
	for scanner.Scan() {
		additionalCount++
	}
	assert.Equal(t, 9, additionalCount)
}

func TestHandlesInvalidNumberOfSampleFileColumns(t *testing.T) {
	file, err := os.Open("bad_sample_test_file.csv")
	assert.Nil(t, err)
	defer file.Close()

	_, _, _, err = readFileForCountTotals(file)
	assert.Equal(t, err.Error(), "Too few columns in CSV file")
}

func TestDuplicateSampleIdCausesError(t *testing.T) {
    file, err := os.Open("duplicate_sample_test_file.csv")
    assert.Nil(t, err)
    defer file.Close()

    _, _, _, err = readFileForCountTotals(file)
    assert.Equal(t, err.Error(), "Duplicate sample unit in file")
}