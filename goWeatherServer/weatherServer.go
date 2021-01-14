// Modified from source: https://golangcode.com/get-the-request-ip-addr/
// Modified from source: https://github.com/ipapi-co/weather-from-ip-address

package main

import (
	"fmt"
	"encoding/json"
	"net/http"
	"io/ioutil"
	"strings"
	"os"
	"log"
)

// Returns temperature based on the client's IP address
// Requires X-FORWARDED-FOR header to be set to the IP address in question
func main() {
	http.HandleFunc("/", RequestHandler)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func RequestHandler(w http.ResponseWriter, r *http.Request) {
	
	w.Header().Add("Content-Type", "application/json")

	temp, message := GetTemperature(r)
	var resp []byte

	if message != "" {
		w.WriteHeader(http.StatusInternalServerError)
		resp, _ = json.Marshal(map[string]string{
			"message": message,
		})
	} else {
		resp, _ = json.Marshal(map[string]float64{
			"temperature": temp,
		})
	}

	w.Write(resp)
}

// Returns temperature based on the client's IP address
// Requires X-FORWARDED-FOR header to be set to the IP address in question
func GetTemperature(r *http.Request) (float64, string) {

	lat, long, message := GetCoord(r)

	if message != "" {
		return 0, message
	}

	apiKey := os.Getenv("OPENWEATHER_API_KEY")
	openWeatherApiUrl := "http://api.openweathermap.org/data/2.5/weather"
	reqWeather, err := http.NewRequest("GET", openWeatherApiUrl + "?units=metric&lat=" + lat + "&lon=" + long + "&appid=" + apiKey, nil)

	client := &http.Client{}
	respWeather, err := client.Do(reqWeather)
    if err != nil {
		// log.Fatal( openWeatherApiUrl + " returned an error: " + err.Error())
        return 0, openWeatherApiUrl + " returned an error"
    }
	defer respWeather.Body.Close()
	
	bodyWeather, _ := ioutil.ReadAll(respWeather.Body)

	var result map[string]interface{}

	json.Unmarshal(bodyWeather, &result)
	if result == nil {
		// log.Fatal(openWeatherApiUrl + " returned an empty response")
		return 0, openWeatherApiUrl + " returned an empty response"
	}

	main := result["main"].(map[string]interface{})
	temp := main["temp"].(float64)

	log.Printf("Temperature is %f\n", temp)

	return temp, ""
}

// Returns latitude and longitude of the client's IP address
// Requires X-FORWARDED-FOR header to be set to the IP address in question
func GetCoord(r *http.Request) (string, string, string) {

	clientIp := r.Header.Get("X-FORWARDED-FOR")

	if clientIp == "" {
		// log.Fatal("Forwarded header not found")
		return "0", "0", "Could not identify the client's IP"
	}

	log.Println("Client IP: " + clientIp)

	ipApiUrl := "https://ipapi.co/"

	reqLatLong, err := http.NewRequest("GET", ipApiUrl + string(clientIp) + "/latlong/", nil)
    client := &http.Client{}
	respLatLong, err := client.Do(reqLatLong)
	
    if err != nil {
        // log.Fatal(ipApiUrl + " returned an error " + err.Error())
		return "0", "0", ipApiUrl + " returned an error"
    }
	defer respLatLong.Body.Close()
	
	if respLatLong.StatusCode != 200 {
		bodyLatLong, _ := ioutil.ReadAll(respLatLong.Body)
		// log.Fatalf("%s returned a status code %d %s", ipApiUrl, respLatLong.StatusCode, string(bodyLatLong))
		return "0", "0", fmt.Sprintf("%s returned a status code %d %s", ipApiUrl, respLatLong.StatusCode, string(bodyLatLong))
	}
	
	bodyLatLong, _ := ioutil.ReadAll(respLatLong.Body)

	latlong :=strings.Split(string(bodyLatLong), ",")

	log.Println("Lat, Long: " + latlong[0] + ", " + latlong[1])
	return latlong[0], latlong[1], ""
}