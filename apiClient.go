package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	time2 "time"
)

const ipApiEndpoint = "https://api.ipify.org?format=json"
const latLongApiEndpoint = "http://ip-api.com/json"
const weatherApiEndpoint = "https://api.open-meteo.com/v1/forecast"

//const temperatureUnit = "celsius"
const temperatureUnit = "fahrenheit"

type apiClient struct {
	IpAddress string
	Latitude  float32
	Longitude float32
	Timezone  string
}

type ipAddressResponse struct {
	Ip string `json:"ip"`
}

type locationResponse struct {
	Timezone string  `json:"timezone"`
	Lat      float32 `json:"lat"`
	Lon      float32 `json:"lon"`
}

type weatherResponse struct {
	Hourly hourly `json:"hourly"`
}

type hourly struct {
	Times []time `json:"time"`
	Temp  []temp `json:"temperature_2m"`
}
type time string
type temp float32

func newClient() *apiClient {
	c := apiClient{}
	//c.getIp()
	return &c
}

func (api *apiClient) getIp() {
	var ip ipAddressResponse

	resp, requestErr := http.Get(ipApiEndpoint)
	if requestErr != nil {
		log.Fatalln(requestErr)
	}

	body, bodyReadErr := ioutil.ReadAll(resp.Body)
	if bodyReadErr != nil {
		log.Fatalln(bodyReadErr)
	}

	unmarshallErr := json.Unmarshal(body, &ip)
	if unmarshallErr != nil {
		log.Fatalln(unmarshallErr)
	}

	api.IpAddress = ip.Ip
	//fmt.Printf("Ip address is %s.\n", api.IpAddress)

}

func (api *apiClient) getLocation() {
	var location locationResponse

	endpoint := fmt.Sprintf("%s/%s", latLongApiEndpoint, api.IpAddress)
	resp, requestErr := http.Get(endpoint)
	if requestErr != nil {
		log.Fatalln(requestErr)
	}

	body, bodyReadErr := ioutil.ReadAll(resp.Body)
	if bodyReadErr != nil {
		log.Fatalln(bodyReadErr)
	}

	unmarshallErr := json.Unmarshal(body, &location)
	if unmarshallErr != nil {
		log.Fatalln(unmarshallErr)
	}

	api.Latitude = location.Lat
	api.Longitude = location.Lon
	api.Timezone = location.Timezone
	//fmt.Printf("API %#v.\n", api)

}

func (api *apiClient) getWeather() {
	var weather weatherResponse

	endpoint := fmt.Sprintf("%s?temperature_unit=%s&latitude=%f&longitude=%f&timezone=%s&hourly=temperature_2m", weatherApiEndpoint, temperatureUnit, api.Latitude, api.Longitude, api.Timezone)
	fmt.Println(endpoint)
	resp, requestErr := http.Get(endpoint)
	if requestErr != nil {
		log.Fatalln(requestErr)
	}

	body, bodyReadErr := ioutil.ReadAll(resp.Body)
	if bodyReadErr != nil {
		log.Fatalln(bodyReadErr)
	}

	unmarshallErr := json.Unmarshal(body, &weather)
	if unmarshallErr != nil {
		log.Fatalln(unmarshallErr)
	}

	now := time2.Now()
	var dataPoints []string

	for i, time := range weather.Hourly.Times {
		temp := weather.Hourly.Temp[i]
		t := parseTime(string(time), api.Timezone)

		if t.After(now) && len(dataPoints) <= 24 {
			dataPoints = append(dataPoints, fmt.Sprintf("Weather at %v is %v.\n", time, temp))
		}

	}
	fmt.Println(dataPoints)

}

func parseTime(timeString string, location string) time2.Time {
	loc, _ := time2.LoadLocation(location)

	const layout = "2006-01-02T15:04"
	t, _ := time2.ParseInLocation(layout, timeString, loc)
	return t
}
