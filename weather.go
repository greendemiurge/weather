package main

func main() {
	client := newClient()
	client.getIp()
	client.getLocation()
	client.getWeather()

}
