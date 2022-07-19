package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

const ipApiEndpoint = "https://api.ipify.org?format=json"

type apiClient struct {
	IpAddress string
	Latitude  float32
	Longitude float32
}

type ipAddressResponse struct {
	Ip string `json:"ip"`
}

func newClient() *apiClient {
	c := apiClient{}
	//c.getIp()
	return &c
}

func (api apiClient) getIp() {
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
	//fmt.Printf("Ip address is %s.", api.IpAddress)

}
