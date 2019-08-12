package main

import(
	"os"
	"strings"
	"strconv"
	"log"
	"util"
	"github.com/hbakhtiyor/strsim"
	"bytes"
	"encoding/gob"
	"github.com/umahmood/haversine"
	"sort"
    "html/template"
    "net/http"
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
func computeFuel(strPort1, strPort2 string, airportMap map[string]latlon) (float64, float64) {
	var ret float64
	var coords1, coords2 latlon

	coords1 = airportMap[strPort1]
	coords2 = airportMap[strPort2]


	r1 := haversine.Coord{Lat: coords1.Lat, Lon: coords1.Lon}  // Oxford, UK
	r2  := haversine.Coord{Lat: coords2.Lat, Lon: coords2.Lon}  // Turin, Italy
	_, km := haversine.Distance(r1, r2)
	
	ret= float64(litersKilometer)*km
	return ret, km
}
func main() {
	lines,err := util.ReadLines("airports.dat")
	if err != nil{
		log.Fatal(err)
	}
	airportMap := getMap(lines)


	tmpl := template.Must(template.ParseFiles("index.html"))

    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method != http.MethodPost {
            tmpl.Execute(w, nil)
            return
        }

		port1 := r.FormValue("from")
		port2 := r.FormValue("to")

		strPort1 := findAirport(port1, airportMap)
		strPort2 := findAirport(port2, airportMap)
		
		fromCountry := airportMap[strPort1].Country
		toCountry := airportMap[strPort2].Country

		consumption, km := computeFuel(strPort1, strPort2, airportMap)

        tmpl.Execute(w, struct{ 
			Success bool; 
			Km float64; 
			Consumption float64;
			FromAirport string;
			ToAirport string;
			FromCountry string;
			ToCountry string
			}{
				true, 
				km, 
				consumption,
				strPort1,
				strPort2,
				fromCountry,
				toCountry,
			})
	})
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("/static/"))))
    http.ListenAndServe(":8080", nil)
}