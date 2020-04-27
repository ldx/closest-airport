# A simple tool to check what the closest airport/weather station is

## Install

Use `go get github.com/ldx/closest-airport`, or download a [release](https://github.com/ldx/closest-airport/releases).

## Usage

Note: you need [Geoclue2](https://developer.gnome.org/platform-overview/unstable/tech-geoclue2.html.en) installed and running for geolocation.

    $ closest-airport
    {"Latitude":37.741501439884814,"Longitude":-122.4367,"Timestamp":1587946519,"Distance":14.66392645140384,"ClosestAirport":{"iso_country":"US","iata_code":"SFO","local_code":"SFO","ident":"KSFO","continent":"NA","iso_region":"US-CA","coordinates":"-122.375, 37.61899948120117","name":"San Francisco International Airport","municipality":"San Francisco","elevation_ft":"13","type":"large_airport","gps_code":"KSFO"}}

If you only need an airport code for looking up weather in your area:

    $ closest-airport | jq -r '.ClosestAirport.ident'
    KSFO
