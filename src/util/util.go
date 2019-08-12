package util

import (
    "os"
    "bufio"
)

//https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
// readLines reads a whole file into memory
// and returns a slice of its lines.
func ReadLines(path string) ([]string, error) {
    file, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    var lines []string
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        lines = append(lines, scanner.Text())
    }
    return lines, scanner.Err()
}

//https://stackoverflow.com/questions/34259800/is-there-a-built-in-min-function-for-a-slice-of-int-arguments-or-a-variable-numb
func MinMax(array []float64) (float64, float64, int, int) {
    var max float64 = array[0]
	var min float64 = array[0]
	var indmin int = 0
	var indmax int = 0
    for i, value := range array {
        if max < value {
			max = value
			indmax = i
        }
        if min > value {
			min = value
			indmin = i
        }
    }
    return min, max, indmin, indmax
}