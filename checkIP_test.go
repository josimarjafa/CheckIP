/*
 * Copyright (c) 2019. Josimar Andrade, No Rights Reserved
 */

package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// test command line execution
func Test_ipapiAPIconectivity(t *testing.T) {

	url := "https://ipapi.co/8.8.8.8/json"

	resp, err := http.Get(url)
	if err != nil {
		t.Errorf("Error to connet to www.ipapi.cp (%v)", err)
	} else {
		if resp.StatusCode <= 200 && resp.StatusCode >= 299 {
			t.Errorf("handler returned wrong status code: got %v want %v",
				resp.StatusCode, http.StatusOK)
		}
	}

}

// test command line execution
func Test_cmd(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{"Google DNS", args{"8.8.8.8"}, "8.8.8.8"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := cmd(tt.args.ip); got != tt.want {
				t.Errorf("cmd() = %v, location %v", got, tt.want)
			}
		})
	}
}

// test ipapi.co to retrieve information's
func Test_ipInfo(t *testing.T) {
	type args struct {
		ip string
	}
	tests := []struct {
		name     string
		args     args
		location string
		timezone string
		returnIP string
	}{
		{"Google DNS", args{"8.8.8.8"}, "United States/California/Mountain View", "America/Los_Angeles", "8.8.8.8"},
		{"Dyn", args{"216.146.35.35"}, "United States/New Jersey/Township of Piscataway", "", "216.146.35.35"},
		{"home.cern", args{"188.184.9.235"}, "Switzerland/Geneva/Geneva", "Europe/Zurich", "188.184.9.235"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, got2 := ipInfo(tt.args.ip)
			if got != tt.location {
				t.Errorf("ipInfo() got = %v, location %v", got, tt.location)
			}
			if got1 != tt.timezone {
				t.Errorf("ipInfo() got1 = %v, timezone %v", got1, tt.timezone)
			}
			if got2 != tt.returnIP {
				t.Errorf("ipInfo() got2 = %v, IP %v", got2, tt.returnIP)
			}
		})
	}
}

func Test_webHome(t *testing.T) {

	req, err := http.NewRequest("GET", "", nil)
	req.Header.Set("X-Forwarded-For", "8.8.8.8")
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webHome)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("webHome returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	expected :=
		`<html><head><title>Current IP Check</title></head><body>Current IP Address: 8.8.8.8 <br>Location: United States/California/Mountain View<br>Timezone: America/Los_Angeles</body></html>`

	t.Run("IP from User Agent", func(t *testing.T) {
		if rr.Body.String() != expected {
			t.Errorf("webHome returned unexpected body: got \n|%v| want \n|%v|",
				rr.Body.String(), expected)
		}
	})
}

func Test_webHealth(t *testing.T) {
	req, err := http.NewRequest("GET", "/status/", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(webHealth)
	handler.ServeHTTP(rr, req)
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("/health returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.

	expected := "up"

	t.Run("health", func(t *testing.T) {
		if rr.Body.String() != expected {
			t.Errorf("/health returned unexpected body: got \n|%v| want \n|%v|",
				rr.Body.String(), expected)
		}
	})
}