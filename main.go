package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
	"strings"

	dbus "github.com/godbus/dbus/v5"
	"github.com/ldx/go-geoclue2"
)

const (
	earthRadiusKm = 6371
)

type Airport struct {
	IsoCountry   string `json:"iso_country"`
	IATACode     string `json:"iata_code"`
	LocalCode    string `json:"local_code"`
	Ident        string `json:"ident"`
	Continent    string `json:"continent"`
	ISORegion    string `json:"iso_region"`
	Coordinates  string `json:"coordinates"`
	Name         string `json:"name"`
	Municipality string `json:"municipality"`
	ElevationFt  string `json:"elevation_ft"`
	Type         string `json:"type"`
	GPSCode      string `json:"gps_code"`
}

func degreesToRadians(degrees float64) float64 {
	return degrees * math.Pi / 180
}

func distanceBetween(lat1, lon1, lat2, lon2 float64) float64 {
	dLat := degreesToRadians(lat2 - lat1)
	dLon := degreesToRadians(lon2 - lon1)
	rad1 := degreesToRadians(lat1)
	rad2 := degreesToRadians(lat2)
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Sin(dLon/2)*math.Sin(dLon/2)*math.Cos(rad1)*math.Cos(rad2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return earthRadiusKm * c
}

func getLocation() (*geoclue2.Location, error) {
	conn, err := dbus.SystemBus()
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	g := geoclue2.NewGeoClue2(conn, "")
	g.Start()
	defer g.Stop()
	return g.WaitForLocation(context.Background())
}

func getAirports(file string) ([]Airport, error) {
	contents, err := ioutil.ReadFile(file)
	if err != nil {
		return nil, err
	}
	var airports []Airport
	err = json.Unmarshal(contents, &airports)
	if err != nil {
		return nil, err
	}
	return airports, nil
}

func getClosestAirport(lat, lon float64, airports []Airport) (*Airport, float64) {
	minDistance := math.MaxFloat64
	var closest Airport
	for _, airport := range airports {
		if airport.IATACode == "" ||
			airport.Ident == "" ||
			!strings.Contains(airport.Type, "airport") {
			continue
		}
		coordinates := airport.Coordinates
		parts := strings.SplitN(coordinates, ",", 2)
		lonA, err := strconv.ParseFloat(strings.Trim(parts[0], " "), 64)
		if err != nil {
			fmt.Fprintf(
				os.Stderr, "invalid longitude for %+v %v\n", airport, err)
			continue
		}
		latA, err := strconv.ParseFloat(strings.Trim(parts[1], " "), 64)
		if err != nil {
			fmt.Fprintf(
				os.Stderr, "invalid latitude for %+v %v\n", airport, err)
			continue
		}
		d := distanceBetween(lat, lon, latA, lonA)
		if d < minDistance {
			closest = airport
			minDistance = d
		}
	}
	if minDistance < math.MaxFloat64 {
		return &closest, minDistance
	}
	return nil, minDistance
}

type Result struct {
	Latitude       float64
	Longitude      float64
	Timestamp      uint64
	Distance       float64
	ClosestAirport Airport
}

func main() {
	location, err := getLocation()
	if err != nil {
		panic(err)
	}
	airports, err := getAirports("airport-codes-pp.json")
	if err != nil {
		panic(err)
	}
	closest, distance := getClosestAirport(
		location.Latitude, location.Longitude, airports)
	if closest == nil {
		panic("{\"error\": \"no airport found\"}")
	}
	result := Result{
		Latitude:       location.Latitude,
		Longitude:      location.Longitude,
		Timestamp:      location.Timestamp.Seconds,
		Distance:       distance,
		ClosestAirport: *closest,
	}
	buf, err := json.Marshal(&result)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s\n", string(buf))
}
