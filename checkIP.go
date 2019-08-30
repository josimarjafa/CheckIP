/*
 * Copyright (c) 2019. Josimar Andrade, No Rights Reserved
 */

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
)


// port configuration
func port() string{
	port := os.Getenv("PORT")
	if len(port) == 0{
		port = "8080"
	}
	return ":" + port
}


// host info retrived from ipapi.co (json format)
type Location struct {
	Ip string `json:"ip"`
	Timezone string `json:"timezone"`
	City string `json:"city"`
	Region string `json:"region"`
	Country_name string `json:"country_name"`
}

// html template
const html  =
	`<html><head><title>Current IP Check</title></head><body>Current IP Address: %v <br>Location: %v<br>Timezone: %v</body></html>`

// http webHome
func webHome(w http.ResponseWriter, r *http.Request) {

	w.WriteHeader(http.StatusOK)
	var ip string
	if r.Header.Get("X-Forwarded-For") != ""{
		ip = r.Header.Get("X-Forwarded-For")
	} else {
		ip, _, _ = net.SplitHostPort(r.RemoteAddr)
	}

	location, timezone, _:= ipInfo(ip)

	log.Printf("<< New request, [ip: %v, location: %s, timezone: %s]", ip, location, timezone)

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	fmt.Fprintf(w, html, ip, location, timezone)
}

// webHealth
func webHealth(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "up")
}

func main() {

	// comand line execution, expect an IP address
	if len(os.Args) == 2 {
		cmd(os.Args[1])
		return
	}

	// print server info
	fmt.Print("\n\n")
	log.Print("Web Service started")
	log.Print("Server address: ",myaddres(), port(), "\n\n")

	// register the function which will handle the requests
	http.HandleFunc("/", webHome)
	http.HandleFunc("/health", webHealth)
	log.Fatal(http.ListenAndServe(port(), nil)) // start listening server requests

}

// handle the command line execution
func cmd(ip string) (string) {
	location, timezone, ipResult := ipInfo(ip)

	fmt.Printf("\nCurrent IP Address: %v \nLocation: %v\nTimezone: %v\n", ip, location, timezone)

	return ipResult
}


// Retrieve IP information from ipapi.co
// receive IP address as argument and return the location, timezone and the IP from ipapi.co
func ipInfo(ip string) (string, string, string){

	url := "https://ipapi.co/" + ip + "/json"
	log.Print("request sent to : ", url)
	// prepare the request
	req, err := http.NewRequest("GET", url, bytes.NewBuffer([]byte(``)))
	if err != nil {
		log.Print("Error: ", err.Error())
		return "", "", ""
	}
	req.Header.Set("Content-Type", "application/json")

	// do the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Print("Error: ", err.Error())
		return "", "", ""
	}
	defer resp.Body.Close() // close connection when finishing

	body, _ := ioutil.ReadAll(resp.Body) // read the body result

	// parse JSON data to location object
	hostLocation := Location{}
	jsonErr := json.Unmarshal(body, &hostLocation)
	if jsonErr != nil {
		log.Print(jsonErr)
	}

	// build the location string with the available information
	var location string

	switch  {
		case hostLocation.Country_name != "":
			location = hostLocation.Country_name
			fallthrough
		case hostLocation.Region != "":
			location = location + "/" + hostLocation.Region
			fallthrough
		case hostLocation.City != "":
			location = location + "/" + hostLocation.City
	}

	return location, hostLocation.Timezone, hostLocation.Ip
}

// Retrieve server IP, for usage help
func myaddres() string {
	addrs, err := net.InterfaceAddrs() // get all interfaces

	if err != nil {
		fmt.Println(err)
	}

	var currentIP string

	for _, address := range addrs { // get the address info

		// check the address type and if it is not a loopback the display it
		// = GET LOCAL IP ADDRESS
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				currentIP = ipnet.IP.String()
			}
		}
	}

	return currentIP
}
