package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	time2 "time"
)

const ipApiEndpoint = "https://api.ipify.org?format=json"
const latLongApiEndpoint = "http://ip-api.com/json"
const weatherApiEndpoint = "https://api.open-meteo.com/v1/forecast"

//const temperatureUnit = "celsius"
const temperatureUnit = "fahrenheit"

type apiClient struct {
	IpAddress       string
	Latitude        float32
	Longitude       float32
	Timezone        string
	WeatherCodesMap map[float32]string
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
	Times       []time        `json:"time"`
	Temp        []temp        `json:"temperature_2m"`
	WeatherCode []weatherCode `json:"weathercode"`
}
type time string
type temp float32
type weatherCode float32

func newClient() *apiClient {
	c := apiClient{}
	c.WeatherCodesMap = make(map[float32]string)
	c.WeatherCodesMap[0.] = string('\U0001F31E')  // 🌞
	c.WeatherCodesMap[1.] = string('\U0001F324')  // 🌤
	c.WeatherCodesMap[2.] = string('\U0001F325')  // 🌥
	c.WeatherCodesMap[3.] = string('\U0001F325')  // 🌥
	c.WeatherCodesMap[51.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[53.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[55.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[45.] = string('\U0001F32B') // 🌫 fog
	c.WeatherCodesMap[48.] = string('\U0001F32B') // 🌫 fog
	c.WeatherCodesMap[61.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[63.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[65.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[71.] = string('\U0001F328') // 🌨 snow
	c.WeatherCodesMap[73.] = string('\U0001F328') // 🌨 snow
	c.WeatherCodesMap[75.] = string('\U0001F328') // 🌨 snow
	c.WeatherCodesMap[77.] = string('\U0001F328') // 🌨 snow
	c.WeatherCodesMap[80.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[81.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[82.] = string('\U0001F327') // 🌧 rain
	c.WeatherCodesMap[85.] = string('\U0001F328') // 🌨 snow
	c.WeatherCodesMap[86.] = string('\U0001F328') // 🌨 snow
	c.WeatherCodesMap[95.] = string('\U0001F329') // 🌩 thunder
	c.WeatherCodesMap[96.] = string('\U0001F329') // 🌩 thunder
	c.WeatherCodesMap[99.] = string('\U0001F329') // 🌩 thunder

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

}

func (api *apiClient) getWeather() {
	var weather weatherResponse
	endpoint := fmt.Sprintf("%s?temperature_unit=%s&latitude=%f&longitude=%f&timezone=%s&hourly=temperature_2m,weathercode", weatherApiEndpoint, temperatureUnit, api.Latitude, api.Longitude, api.Timezone)
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
	var days []string
	var times []string
	var dataPoints []string

	for i, time := range weather.Hourly.Times {
		temp := weather.Hourly.Temp[i]
		weatherCode := weather.Hourly.WeatherCode[i]
		t := parseTime(string(time), api.Timezone)

		if t.After(now) && len(dataPoints) <= 18 {
			days = append(days, t.Format("Mon    "))
			times = append(times, t.Format("03:04PM"))
			dataPoints = append(dataPoints, fmt.Sprintf("%.1f %s", temp, api.WeatherCodesMap[float32(weatherCode)]))
		}

	}

	fmt.Println(strings.Join(days, " | "))
	fmt.Println(strings.Join(times, " | "))
	fmt.Println(strings.Join(dataPoints, " | "))

}

func parseTime(timeString string, location string) time2.Time {
	loc, _ := time2.LoadLocation(location)
	const layout = "2006-01-02T15:04"
	t, _ := time2.ParseInLocation(layout, timeString, loc)
	return t
}
