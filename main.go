package main

import(
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
	"fmt"
	"github.com/hbakhtiyor/strsim"
	//"github.com/umahmood/haversine"
	"sort"
)
const(
	gallonsMile int = 5
	litersKilometer int = 12
	iAirport int = 1
	iLat int = 6
	iLon int = 7
	iCountry int = 3
)
type latlon struct{
	lat float64
	lon float64
	country string
}

//https://stackoverflow.com/questions/5884154/read-text-file-into-string-array-and-write
// readLines reads a whole file into memory
// and returns a slice of its lines.
func readLines(path string) ([]string, error) {
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
func makeMap(lines []string) map[string]latlon{
	ret := make(map[string]latlon)
	for _, line := range lines{
		splitted := strings.Split(line, ",")
		lat,err := strconv.ParseFloat(splitted[iLat], 64)
		lon,err := strconv.ParseFloat(splitted[iLon], 64)
		country := splitted[iCountry]
		airport := splitted[iAirport]
		if err != nil{
			log.Fatal(err)
		}
		info := latlon{
			lat: lat,
			lon: lon,
			country: country,
		}
		ret[airport] = info

	}
	return ret
}
func findAirport(key string, m map[string]latlon) string{
	_, ok := m[key]
	if ok {
		return key
	} else {
		var keys []string
		for akey := range m {
			keys = append(keys, akey)
		}
		sort.Strings(keys)
		distances := make([]float64, len(m))
		i:=0
		for _,k := range(keys){
			// fmt.Println(k)
			distances[i] = strsim.Compare(key, k)
			i = i + 1
		}
		_,_, _, iClosest := MinMax(distances)

		closestKey := keys[iClosest]
		return closestKey
	}
}
func main() {
	lines,err := readLines("airports.dat")
	if err != nil{
		log.Fatal(err)
	}
	airportMap := makeMap(lines)

	fmt.Println(findAirport("Arlanda", airportMap))

}