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
	if apiKey == "" {
		log.Println( "OPENWEATHER_API_KEY is required, please set the environment variable, e.g. by running the following command in the terminal: export OPENWEATHER_API_KEY=<Your API key>")
        return 0, "OPENWEATHER_API_KEY is missing"
	}

	const openWeatherApiUrl = "http://api.openweathermap.org/data/2.5/weather"
	reqWeather, err := http.NewRequest("GET", openWeatherApiUrl + "?units=metric&lat=" + lat + "&lon=" + long + "&appid=" + apiKey, nil)

	client := &http.Client{}
	respWeather, err := client.Do(reqWeather)
    if err != nil {
		log.Println( openWeatherApiUrl + " returned an error: " + err.Error())
        return 0, openWeatherApiUrl + " returned an error"
    }
	defer respWeather.Body.Close()
	
	bodyWeather, _ := ioutil.ReadAll(respWeather.Body)

	var result map[string]interface{}

	json.Unmarshal(bodyWeather, &result)
	if result == nil {
		log.Println(openWeatherApiUrl + " returned an empty response")
		return 0, openWeatherApiUrl + " returned an empty response"
	}

	main := result["main"].(map[string]interface{})

	if main == nil {
		log.Println(openWeatherApiUrl + " returned an empty response")
		return 0, openWeatherApiUrl + " returned an empty response"
	}
	temp := main["temp"].(float64)

	log.Printf("Temperature is %f\n", temp)

	return temp, ""
}

// Returns latitude and longitude of the client's IP address
// Requires X-FORWARDED-FOR header to be set to the IP address in question
func GetCoord(r *http.Request) (string, string, string) {

	clientIp := r.Header.Get("X-FORWARDED-FOR")

	if clientIp == "" {
		log.Println("Forwarded header not found")
		return "0", "0", "Could not identify the client's IP"
	}

	log.Println("Client IP: " + clientIp)

	var coords Coord

	if os.Getenv("IPIFY_API_KEY") != "" {
		coords = new(Ipify)
	} else {
		coords = new(IpApi)
	}

	return coords.GetCoord(clientIp)
}

type Coord interface {
	GetCoord(clientIp string) (string, string, string)
}

type IpApi struct { // throws RateLimited error all too often on a free plan
}

type Ipify struct { // 1000 requests per month on free subscription. Requires environment variable IPIFY_API_KEY to be set
}

// Returns latitude and longitude of the client's IP address
// Requires X-FORWARDED-FOR header to be set to the IP address in question
func (coord IpApi) GetCoord(clientIp string) (string, string, string) {
	
	ipApiUrl := "https://ipapi.co/" + string(clientIp) + "/latlong/" // throws RateLimited error all too often on a free plan

	reqLatLong, err := http.NewRequest("GET", ipApiUrl, nil)
    client := &http.Client{}
	respLatLong, err := client.Do(reqLatLong)
	
    if err != nil {
        log.Println(ipApiUrl + " returned an error " + err.Error())
		return "0", "0", ipApiUrl + " returned an error"
    }
	defer respLatLong.Body.Close()
	
	if respLatLong.StatusCode != http.StatusOK {
		bodyLatLong, _ := ioutil.ReadAll(respLatLong.Body)
		log.Printf("%s returned a status code %d %s", ipApiUrl, respLatLong.StatusCode, string(bodyLatLong))
		return "0", "0", fmt.Sprintf("%s returned a status code %d %s", ipApiUrl, respLatLong.StatusCode, string(bodyLatLong))
	}
	
	bodyLatLong, _ := ioutil.ReadAll(respLatLong.Body)

	latlng :=strings.Split(string(bodyLatLong), ",")
	log.Println("Lat, Long: " + latlng[0] + ", " + latlng[1])
	return latlng[0], latlng[1], ""
}

// Returns latitude and longitude of the client's IP address
// Requires X-FORWARDED-FOR header to be set to the IP address in question
// Requires environment variable IPIFY_API_KEY to be set
func (coord Ipify) GetCoord(clientIp string) (string, string, string) {
	
	ipifyApiKey := os.Getenv("IPIFY_API_KEY") // 1000 requests per month on free subscription
	ipApiUrl := "https://geo.ipify.org/api/v1?apiKey=" + ipifyApiKey + "&ipAddress=" + string(clientIp)

	reqLatLong, err := http.NewRequest("GET", ipApiUrl, nil)
    client := &http.Client{}
	respLatLong, err := client.Do(reqLatLong)
	
    if err != nil {
        log.Println(ipApiUrl + " returned an error " + err.Error())
		return "0", "0", ipApiUrl + " returned an error"
    }
	defer respLatLong.Body.Close()
	
	if respLatLong.StatusCode != http.StatusOK {
		bodyLatLong, _ := ioutil.ReadAll(respLatLong.Body)
		log.Printf("%s returned a status code %d %s", ipApiUrl, respLatLong.StatusCode, string(bodyLatLong))
		return "0", "0", fmt.Sprintf("%s returned a status code %d %s", ipApiUrl, respLatLong.StatusCode, string(bodyLatLong))
	}
	
	bodyLatLong, _ := ioutil.ReadAll(respLatLong.Body)

	var result map[string]interface{}
	json.Unmarshal(bodyLatLong, &result)
	if result == nil {
		log.Println(ipApiUrl + " returned an empty response")
		return "0", "0", ipApiUrl + " returned an empty response"
	}
	location := result["location"].(map[string]interface{})
	lat := fmt.Sprintf("%f", location["lat"].(float64))
	lng := fmt.Sprintf("%f", location["lng"].(float64))
	log.Println("Lat, Long: " + lat + ", " + lng)
	return lat, lng, ""
}