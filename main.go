package main

import(
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
	"fmt"
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

func main() {
	lines,err := readLines("airports.dat")
	if err != nil{
		log.Fatal(err)
	}
	airportMap := makeMap(lines)
	fmt.Println(airportMap)

}