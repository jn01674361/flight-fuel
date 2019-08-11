package main

import(
	"bufio"
	"os"
	"strings"
	"strconv"
	"log"
	"fmt"
	"util"
	"github.com/hbakhtiyor/strsim"
	"bytes"
	"encoding/gob"
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
	Lat float64
	Lon float64
	Country string
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
			Lat: lat,
			Lon: lon,
			Country: country,
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
			distances[i] = strsim.Compare(key, k)
			i = i + 1
		}
		_,_, _, iClosest := util.MinMax(distances)

		closestKey := keys[iClosest]
		return closestKey
	}
}
func getMap(lines []string) map[string]latlon {
	var network bytes.Buffer        // Stand-in for a network connection
	var m map[string]latlon
	if _, err := os.Stat("airports.gob"); err == nil {
		dec := gob.NewDecoder(&network)
		
		err = dec.Decode(&m)
		if err != nil {
			log.Fatal("decode error:", err)
		}

	} else {
		m = makeMap(lines)
		enc := gob.NewEncoder(&network) // Will write to network.
		err := enc.Encode(m)
		if err != nil {
			log.Fatal("encode error:", err)
		}
	}
	return m
}
func main() {
	lines,err := util.ReadLines("airports.dat")
	if err != nil{
		log.Fatal(err)
	}
	airportMap := getMap(lines)

	reader := bufio.NewReader(os.Stdin)
	fmt.Print("Airport 1: ")
	port1, _ := reader.ReadString('\n')
	fmt.Print("Airport 2: ")
	port2, _ := reader.ReadString('\n')
	fmt.Println(findAirport(port1, airportMap))
	fmt.Println(findAirport(port2, airportMap))

}