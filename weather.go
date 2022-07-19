package main

func main() {
	client := newClient()
	client.getIp()
	client.getLocation()
	client.getWeather()
	//t := parseTime("2022-07-18T00:00", "America/Chicago")
	//fmt.Printf("time is %v.\n", t)

}
